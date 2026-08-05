package llm

type LLMHandler struct {
	wordchunks chan string
}

/*
 *
 func (a *STT) InvokeLLM(ctx context.Context, prompt string) error {
	tokens := make(chan string, 32)
	go func() {
		systemPrompt := "You are an F1 race engineer. While testing, named Dumber Bono. Make shit up for now. Give conversational very very short answers Its a high panic situation you get a few seconds to speak. Like one liners"
		if err := llm.StreamLLM(ctx, systemPrompt, prompt, "openai/gpt-oss-20b", tokens); err != nil {
			log.Println("stream error:", err)
		}
	}()
	for tok := range tokens {
		a.wordchunks <- tok
	}
	return nil
 }
*/
