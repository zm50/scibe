package global

import "github.com/tmc/langchaingo/llms/ollama"

var llm *ollama.LLM

func InitLLM() {
	var err error
	llm, err = ollama.New(ollama.WithServerURL(Config().OllamaEndpoint), 
	ollama.WithModel(Config().Model))
	if err != nil {
		panic(err)
	}
}

func LLM() *ollama.LLM {
	return llm
}
