package chat

import (
	"LaoQGChat/api/models/chat"
	"github.com/gin-gonic/gin"
)

type geminiAPI struct{}

func newGeminiAPI() (*geminiAPI, error) {
	return &geminiAPI{}, nil
}

func (api *geminiAPI) chat(ctx *gin.Context, model string, contexts []chat.Context) (*chat.Response, error) {
	response := &chat.Response{}
	return response, nil
}
