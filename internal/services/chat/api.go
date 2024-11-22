package chat

import (
	"LaoQGChat/api/models/chat"
	"LaoQGChat/internal/myerrors"
	"strings"

	"github.com/gin-gonic/gin"
)

type externalAPI interface {
	chat(ctx *gin.Context, model string, contexts []chat.Content) (*chat.Response, error)
}

func getExternalAPI(model string) (externalAPI, error) {
	var (
		api externalAPI
		err error
	)

	switch getModelName(model) {
	case "azopenai":
		if api, err = newAzopenaiAPI(); err != nil {
			return nil, err
		}
		return api, nil
	case "gemini":
		if api, err = newGeminiAPI(); err != nil {
			return nil, err
		}
		return api, nil
	}
	err = &myerrors.CustomError{
		StatusCode:  200,
		MessageCode: "ECH0050",
		MessageText: "不受支持的模型名称。",
	}
	return nil, err
}

func getModelName(model string) string {
	if strings.Contains(model, "$") {
		return strings.Split(model, "$")[0]
	} else {
		return model
	}
}

func getModelId(model string) string {
	if strings.Contains(model, "$") {
		return strings.Split(model, "$")[1]
	} else {
		return ""
	}
}
