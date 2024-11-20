package chat

import (
	"LaoQGChat/api/models/chat"
	"github.com/gin-gonic/gin"
)

type openaiAPI struct{}

func newOpenAIAPI() (*openaiAPI, error) {
	return &openaiAPI{}, nil
}

func (api *openaiAPI) chat(ctx *gin.Context, model string, contents []chat.Content) (*chat.Response, error) {
	response := &chat.Response{}
	return response, nil
}