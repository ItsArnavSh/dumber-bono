package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// TyreDataInput defines the parameters the model can pass to the tool.
type TyreDataInput struct {
	Tyre string `json:"tyre" jsonschema_description:"which tyre to check: front-left, front-right, rear-left, rear-right"`
}

// getTyreData is the actual Go function the tool calls.
func getTyreData(ctx context.Context, input *TyreDataInput) (string, error) {
	// Fake data for now — later hook this into real sim telemetry
	return fmt.Sprintf("%s tyre temp: 98°C, wear: 62%%", input.Tyre), nil
}

func main() {
	ctx := context.Background()

	cm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey:  os.Getenv("GROQ_API_KEY"),
		BaseURL: "https://api.groq.com/openai/v1",
		Model:   "llama-3.3-70b-versatile",
	})
	if err != nil {
		log.Fatal(err)
	}

	instruction := "You are a F1 Race Engineer. Answer in very short sentences you get a few seconds. Make stuff up for now"

	tyreTool, err := utils.InferTool(
		"get_tyre_data",
		"Get current temperature and wear percentage for a specific tyre",
		getTyreData,
	)
	if err != nil {
		log.Fatal(err)
	}

	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "Ch02ChatModelAgent",
		Description: "A minimal ChatModelAgent with in-memory multi-turn history.",
		Instruction: instruction,
		Model:       cm,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{tyreTool},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	runner := adk.NewTypedRunner(adk.TypedRunnerConfig[*schema.Message]{
		Agent:           agent,
		EnableStreaming: true,
	})

	events := runner.Query(ctx, "What is the tyre temperature?")

	for {
		event, ok := events.Next()
		if !ok {
			break // Iterator closed, all events consumed
		}
		if event.Err != nil {
			fmt.Println("event error:", event.Err)
			continue
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}
		stream := event.Output.MessageOutput.MessageStream
		if stream == nil {
			continue
		}
		for {
			frame, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(frame.Content)
		}
	}
	fmt.Println()
}
