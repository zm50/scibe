package handler

import (
	"fmt"

	"scibe/global"
	"scibe/utils/resp"

	"github.com/gin-gonic/gin"
)

func ChatCompletion(c *gin.Context) {
	// type CompletionRequest struct {
	// Prompt string `json:"prompt"`
	// }
	prompt, ok := c.Params.Get("prompt")
	if !ok {
		c.Abort()
		return
	}
	// req := CompletionRequest{}

	// err := c.BindJSON(&req)
	// if err != nil {
	// 	c.Abort()
	// 	return
	// }
	answer, err := global.ChatCompletion(global.ChatRequest{
		Model: "qwen2.5:1.5b",
		Prompt: prompt,
	})
	if err != nil {
		resp.AbortErr(c, err.Error())
		return
	}
	fmt.Println("6")
	fmt.Println(answer)

	resp.Ok(c, gin.H{
		"answer": answer,
	})
}
