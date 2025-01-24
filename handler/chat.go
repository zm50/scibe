package handler

import (
	"context"
	"sort"
	"strings"

	"scibe/global"
	"scibe/model"
	"scibe/utils/consts"
	"scibe/utils/ret"
	"scibe/utils/types"

	"github.com/gin-gonic/gin"
	"github.com/tmc/langchaingo/llms"
)

type ChatExtractTextRequest struct {
    Fid uint `json:"fid"`
}

type ChatExtractTextResponseItem struct {
    Title string `json:"title"`
    Texts []string `json:"texts"`
}

type TextItem struct {
    Title string
	Start int
	End int
}

func ChatExtractText(c *gin.Context) {
	prop := types.GetProperty(c)
	log := prop.Logger()

	req := &ChatExtractTextRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.Err(err).Msg("failed to bind request")
		ret.AbortErr(c, err.Error())
		return
	}

	docFile := &model.DocFile{}
	err := global.DB().Where("id = ? AND uid = ?", req.Fid, prop.Uid).First(docFile).Error
	if err != nil {
		log.Err(err).Msg("failed to get file")
		ret.AbortErr(c, err.Error())
		return
	}

	pageText, err := fileExtractText(filePath(docFile.ID))
	if err != nil {
		log.Err(err).Msg("failed to extract text")
		ret.AbortErr(c, err.Error())
		return
	}

	prompt := "以下为论文的原文:\n" + pageText + "\n阅读论文,采用中文回答,必须使用以下为输出格式." +
`
研究背景:
...
研究目的:
...
研究方法:
...
研究结果:
...
研究展望:
...`

	resp, err := llms.GenerateFromSinglePrompt(c.Request.Context(), global.LLM(), prompt)
	if err != nil {
		log.Err(err).Msg("failed to generate content")
		ret.AbortErr(c, err.Error())
		return
	}

	log.Printf("chat extract response: %s", resp)

	array := []TextItem{}
	keys := []string{"研究背景:", "研究目的:", "研究方法:", "研究结果:", "研究展望:"}
	for _, key := range keys {
		idx := strings.Index(resp, key)
		if idx != -1 {
			array = append(array, TextItem{
				Title: key[:len(key) - 1],
				Start: idx,
				End: idx + len(key),
			})
		}
	}

	sort.Slice(array, func(i, j int) bool {
		return array[i].Start < array[j].Start
	})

	items := []ChatExtractTextResponseItem{}
	for i := 0; i < len(array); i++ {
		content := ""
		if i == len(array) - 1 {
			content = resp[array[i].End:]
		} else {
			content = resp[array[i].End:array[i+1].Start]
		}

		texts := strings.Split(content, "\n")

		items = append(items, ChatExtractTextResponseItem{
		    Title: array[i].Title,
		    Texts: texts,
		})
	}

	ret.Ok(c, gin.H{
		"items": items,
	})
}

type ChatPaperSummaryRequest struct {
    Fid uint `json:"fid"`
}

func ChatPaperSummary(c *gin.Context) {
	prop := types.GetProperty(c)
	log := prop.Logger()

	req := &ChatPaperSummaryRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		log.Err(err).Msg("failed to bind request")
		ret.AbortErr(c, err.Error())
		return
	}

	docFile := &model.DocFile{}
	err := global.DB().Where("id = ? AND uid = ?", req.Fid, prop.Uid).First(docFile).Error
	if err != nil {
		log.Err(err).Msg("failed to get file")
		ret.AbortErr(c, err.Error())
		return
	}

    // prompt := "用尽可能简洁的语言，总结下面这篇文章的内容。不得使用 markdown 记号。"
	pageText, err := fileExtractText(filePath(docFile.ID))
	if err != nil {
		log.Err(err).Msg("failed to extract text")
		ret.AbortErr(c, err.Error())
		return
	}
	prompt := "以下为论文的原文:\n" + pageText + "\n简洁的总结这篇文章的内容,采用中文回答,不得使用 markdown 记号"

	resp, err := llms.GenerateFromSinglePrompt(c.Request.Context(), global.LLM(), prompt)
	if err != nil {
		log.Err(err).Msg("failed to generate content")
		ret.AbortErr(c, err.Error())
		return
	}

	ret.Ok(c, gin.H{
		"summary": resp,
	})
}

type ChatCompletionRequest struct {
	Fid string `json:"fid"`
	Prompt string `json:"prompt"`
}

func ChatCompletion(c *gin.Context) {
	prop := types.GetProperty(c)
	log := prop.Logger()

	req := ChatCompletionRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		log.Err(err).Msg("failed to bind request")
		ret.AbortErr(c, err.Error())
		return
	}

	docFile := &model.DocFile{}
	err = global.DB().Where("id = ? AND uid = ?", req.Fid, prop.Uid).First(docFile).Error
	if err != nil {
		log.Err(err).Msg("failed to get file")
		ret.AbortErr(c, err.Error())
		return
	}

	pageText, err := fileExtractText(filePath(docFile.ID))
	if err != nil {
		log.Err(err).Msg("failed to extract text")
		ret.AbortErr(c, err.Error())
		return
	}

	systemPrompt := "以下为论文的原文:\n" + pageText + "\n回答时请注意:" +
`1. 优先使用论文中的内容回答
2. 回答要简洁清晰,并标注信息来源
3. 使用中文回答

输出格式要求：
1. 使用 Markdown 格式输出
2. 重要内容使用**加粗**标记
3. 使用适当的标题层级(###, ####)组织内容
4. 数学公式使用 LaTeX 格式：
   - 行内公式使用单个$包裹
   - 独立公式块使用$$包裹
5. 如有必要使用列表或表格增强可读性
6. 引用原文时使用>引用格式
` + "7. 关键概念或术语使用`代码块`标记"

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	_, err = global.LLM().GenerateContent(c.Request.Context(), []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, req.Prompt),
	}, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		c.SSEvent(consts.SSEventData, string(chunk))
		c.Writer.Flush()
		return nil
	}))
	if err != nil {
		log.Err(err).Msg("failed to generate content")
		ret.AbortErr(c, err.Error())
		return
	}
}
