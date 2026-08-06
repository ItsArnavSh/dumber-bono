package eino

import (
	"context"
	"dubmer-bono/app/service/internal/llm"
	"fmt"
	"io"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

// GroqEinoChatModel adapts the llm package's QueryLLM/StreamLLM functions to
// Eino's model.BaseChatModel interface (Generate + Stream), so it can be
// used as a ChatModel component inside Eino graphs/chains/agents.
//
// Generate uses the non-streaming QueryLLM call.
// Stream uses the real SSE-based StreamLLM call, run in its own goroutine
// writing into an io.Pipe, with each chunk of text forwarded as a
// schema.Message onto the returned StreamReader as it arrives.
type GroqEinoChatModel struct {
	model        string
	systemPrompt string // optional default system prompt if none present in input messages
}

// GroqEinoChatModelConfig configures GroqEinoChatModel.
type GroqEinoChatModelConfig struct {
	// Model is the Groq model name, e.g. "llama-3.3-70b-versatile".
	// If empty, the llm package's default is used.
	Model string

	// DefaultSystemPrompt is used only if the input message slice contains
	// no schema.System message.
	DefaultSystemPrompt string
}

// NewGroqEinoChatModel constructs a new GroqEinoChatModel.
func NewGroqEinoChatModel(_ context.Context, cfg *GroqEinoChatModelConfig) (*GroqEinoChatModel, error) {
	if cfg == nil {
		cfg = &GroqEinoChatModelConfig{}
	}
	return &GroqEinoChatModel{
		model:        cfg.Model,
		systemPrompt: cfg.DefaultSystemPrompt,
	}, nil
}

// compile-time interface check
var _ model.BaseChatModel = (*GroqEinoChatModel)(nil)

// splitMessages extracts a system prompt (last System message wins, falls
// back to configured default) and flattens the remaining messages into a
// single user-role prompt string, since the llm package only accepts a
// single system + single user string.
func (m *GroqEinoChatModel) splitMessages(input []*schema.Message) (systemPrompt, userPrompt string) {
	systemPrompt = m.systemPrompt

	var b strings.Builder
	for _, msg := range input {
		if msg == nil {
			continue
		}
		switch msg.Role {
		case schema.System:
			systemPrompt = msg.Content
		case schema.User:
			if b.Len() > 0 {
				b.WriteString("\n\n")
			}
			b.WriteString(msg.Content)
		case schema.Assistant:
			if b.Len() > 0 {
				b.WriteString("\n\n")
			}
			b.WriteString("Assistant: ")
			b.WriteString(msg.Content)
		case schema.Tool:
			if b.Len() > 0 {
				b.WriteString("\n\n")
			}
			b.WriteString("Tool result: ")
			b.WriteString(msg.Content)
		default:
			if b.Len() > 0 {
				b.WriteString("\n\n")
			}
			b.WriteString(msg.Content)
		}
	}

	return systemPrompt, b.String()
}

// resolveModel applies any per-call model.Option override on top of the
// configured default model name.
func (m *GroqEinoChatModel) resolveModel(opts ...model.Option) string {
	options := model.GetCommonOptions(&model.Options{Model: &m.model}, opts...)
	if options.Model != nil && *options.Model != "" {
		return *options.Model
	}
	return m.model
}

// Generate implements model.BaseChatModel.
func (m *GroqEinoChatModel) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("groq eino chat model: no input messages provided")
	}

	modelName := m.resolveModel(opts...)
	systemPrompt, userPrompt := m.splitMessages(input)

	content, err := llm.QueryLLM(ctx, systemPrompt, userPrompt, modelName)
	if err != nil {
		return nil, fmt.Errorf("groq eino chat model: %w", err)
	}

	return &schema.Message{
		Role:    schema.Assistant,
		Content: content,
	}, nil
}

// Stream implements model.BaseChatModel.
//
// It runs llm.StreamLLM in a goroutine, writing SSE token deltas into an
// io.Pipe. A second goroutine reads from the pipe and forwards each read
// chunk as a schema.Message onto the returned StreamReader, so downstream
// Eino consumers see real incremental output as it arrives from Groq.
func (m *GroqEinoChatModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("groq eino chat model: no input messages provided")
	}

	modelName := m.resolveModel(opts...)
	systemPrompt, userPrompt := m.splitMessages(input)

	pr, pw := io.Pipe()

	// Producer: streams from Groq and writes token text into the pipe.
	go func() {
		err := llm.StreamLLM(ctx, systemPrompt, userPrompt, modelName, pw)
		_ = pw.CloseWithError(err) // nil err just closes cleanly (io.EOF for the reader)
	}()

	sr, sw := schema.Pipe[*schema.Message](1)

	// Consumer: reads pipe chunks and forwards them as schema.Message onto
	// the Eino stream.
	go func() {
		defer sw.Close()

		buf := make([]byte, 4096)
		for {
			n, err := pr.Read(buf)
			if n > 0 {
				chunk := string(buf[:n])
				sw.Send(&schema.Message{
					Role:    schema.Assistant,
					Content: chunk,
				}, nil)
			}
			if err != nil {
				if err == io.EOF {
					return
				}
				sw.Send(nil, fmt.Errorf("groq eino chat model: stream error: %w", err))
				return
			}
		}
	}()

	return sr, nil
}
