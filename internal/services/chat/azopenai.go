package chat

import (
	"LaoQGChat/api/models/chat"
	"LaoQGChat/internal/myerrors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/gin-gonic/gin"
	"os"
)

type azopenaiAPI struct {
	client       azopenai.Client
	defaultModel string
}

func newAzopenaiAPI() (*azopenaiAPI, error) {
	// azopenai认证
	azopenaiKey := os.Getenv("AOAI_API_KEY")
	azopenaiEndpoint := os.Getenv("AOAI_ENDPOINT")
	azopenaiClient, err := azopenai.NewClientWithKeyCredential(azopenaiEndpoint, azcore.NewKeyCredential(azopenaiKey), nil)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "ECH0100",
			MessageText: "Azure OpenAI认证失败，请联系管理员。",
		}
		return nil, err
	}

	// 设置默认模型
	defaultModel := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL")

	return &azopenaiAPI{client: *azopenaiClient, defaultModel: defaultModel}, nil
}

func (api *azopenaiAPI) chat(ctx *gin.Context, model string, contents []chat.Content) (*chat.Response, error) {
	var (
		azopenaiRequest  azopenai.ChatCompletionsOptions
		azopenaiResponse azopenai.GetChatCompletionsResponse
		messages         []azopenai.ChatRequestMessageClassification
		err              error
	)

	// 获取模型实例
	deploymentName := getModelId(model)
	if deploymentName == "" {
		deploymentName = api.defaultModel
	}

	// 将request转为azopenai的输入
	messages, err = toAzopenaiContents(contents)
	if err != nil {
		return nil, err
	}

	// 构造azopenai请求
	azopenaiRequest = azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &deploymentName,
	}

	// 发送azopenai请求
	azopenaiResponse, err = api.client.GetChatCompletions(ctx.Request.Context(), azopenaiRequest, nil)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "ECH0101",
			MessageText: "Azure OpenAI获取答案失败，请联系管理员。",
		}
		return nil, err
	}

	// 设置回答
	answer := chat.Content{
		Type: chat.ContentTypeAnswer,
		Parts: []chat.PartWrapper{
			{Part: &chat.TextPart{
				Type: chat.PartTypeText,
				Text: *azopenaiResponse.Choices[0].Message.Content,
			}},
		},
	}

	// 设置选项
	options := make([]chat.Content, len(azopenaiResponse.Choices)-1)
	for _, choice := range azopenaiResponse.Choices[1:] {
		option := chat.Content{
			Type: chat.ContentTypeOption,
			Parts: []chat.PartWrapper{
				{Part: &chat.TextPart{
					Type: chat.PartTypeText,
					Text: *choice.Message.Content,
				}},
			},
		}
		options = append(options, option)
	}
	response := &chat.Response{
		Answer:  answer,
		Options: options,
	}

	return response, nil
}

func toAzopenaiContents(contents []chat.Content) ([]azopenai.ChatRequestMessageClassification, error) {
	azopenaiContents := make([]azopenai.ChatRequestMessageClassification, len(contents))
	for _, content := range contents {
		azopenaiContent, err := toAzopenaiContent(content)
		if err != nil {
			return nil, err
		}
		azopenaiContents = append(azopenaiContents, azopenaiContent)
	}
	return azopenaiContents, nil
}

func toAzopenaiContent(content chat.Content) (azopenai.ChatRequestMessageClassification, error) {
	parts := make([]chat.Part, len(content.Parts))
	for _, part := range content.Parts {
		parts = append(parts, part.Part)
	}
	azopenaiParts, err := toAzopenaiParts(parts)
	if err != nil {
		return nil, err
	}

	switch content.Type {
	case chat.ContentTypeQuestion:
		return &azopenai.ChatRequestUserMessage{
			Content: azopenai.NewChatRequestUserMessageContent(azopenaiParts),
		}, nil
	case chat.ContentTypeAnswer:
		if part, ok := content.Parts[0].Part.(*chat.TextPart); ok {
			return &azopenai.ChatRequestAssistantMessage{
				Content: &part.Text,
			}, nil
		} else {
			err = &myerrors.CustomError{
				StatusCode:  200,
				MessageCode: "ECH0151",
				MessageText: fmt.Sprintf("不支持的Part类型：%s。", part.GetContentType()),
			}
			return nil, err
		}
	}
	err = &myerrors.CustomError{
		StatusCode:  200,
		MessageCode: "ECH0150",
		MessageText: fmt.Sprintf("不支持的Content类型：%s。", content.Type),
	}
	return nil, nil
}

func toAzopenaiParts(parts []chat.Part) ([]azopenai.ChatCompletionRequestMessageContentPartClassification, error) {
	azopenaiParts := make([]azopenai.ChatCompletionRequestMessageContentPartClassification, len(parts))
	for _, part := range parts {
		azopenaiPart, err := toAzopenaiPart(part)
		if err != nil {
			return nil, err
		}
		azopenaiParts = append(azopenaiParts, azopenaiPart)
	}
	return azopenaiParts, nil
}

func toAzopenaiPart(part chat.Part) (azopenai.ChatCompletionRequestMessageContentPartClassification, error) {
	switch part.(type) {
	case *chat.TextPart:
		textPart := part.(*chat.TextPart)
		azopenaiPart := &azopenai.ChatCompletionRequestMessageContentPartText{
			Text: &textPart.Text,
		}
		return azopenaiPart, nil
	case *chat.ImagePart:
		imagePart := part.(*chat.ImagePart)
		azopenaiPart := &azopenai.ChatCompletionRequestMessageContentPartImage{
			ImageURL: &azopenai.ChatCompletionRequestMessageContentPartImageURL{
				URL: &imagePart.ImageUrl,
			},
		}
		return azopenaiPart, nil
	}
	err := &myerrors.CustomError{
		StatusCode:  200,
		MessageCode: "ECH0151",
		MessageText: fmt.Sprintf("不支持的Part类型：%s。", part.GetContentType()),
	}
	return nil, err
}
