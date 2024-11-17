package chat

import (
	"LaoQGChat/api/models/chat"
	"LaoQGChat/internal/myerrors"
	"LaoQGChat/internal/shared"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"os"
	"reflect"
)

type geminiAPI struct {
	apiKey       string
	defaultModel string
}

func newGeminiAPI() (*geminiAPI, error) {
	// 设置API KEY
	apiKey := os.Getenv("GEMINI_API_KEY")

	// 设置默认模型
	defaultModel := os.Getenv("GEMINI_MODEL")
	return &geminiAPI{apiKey: apiKey, defaultModel: defaultModel}, nil
}

func (api *geminiAPI) chat(ctx *gin.Context, modelStr string, contents []chat.Content) (*chat.Response, error) {
	var (
		client          *genai.Client
		model           *genai.GenerativeModel
		chatSession     *genai.ChatSession
		geminiRequests  []genai.Part
		geminiResponses *genai.GenerateContentResponse
		err             error
	)

	// 创建gemini客户端
	client, err = genai.NewClient(ctx.Request.Context(), option.WithAPIKey(api.apiKey))
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "ECH0200",
			MessageText: "Gemini客户端创建失败，请联系管理员。",
		}
		return nil, err
	}
	defer func(client *genai.Client) {
		_ = client.Close()
	}(client)

	// 设置模型
	modelName := getModelId(modelStr)
	if modelName == "" {
		modelName = api.defaultModel
	}
	model = client.GenerativeModel(modelName)

	// 将历史记录转为gemini历史记录
	chatSession = model.StartChat()
	chatSession.History, err = toGeminiContents(contents[:len(contents)-1])
	if err != nil {
		return nil, err
	}

	// 将请求转为gemini请求
	parts := make([]chat.Part, len(contents))
	for _, part := range contents[len(contents)-1].Parts {
		parts = append(parts, part.Part)
	}
	geminiRequests, err = toGeminiParts(parts)
	if err != nil {
		return nil, err
	}

	// 发送gemini请求
	geminiResponses, err = chatSession.SendMessage(ctx.Request.Context(), geminiRequests...)
	if err != nil {
		return nil, err
	}

	var (
		answer  *chat.Content
		options = make([]chat.Content, len(geminiResponses.Candidates)-1)
	)

	// 设置回答
	answer, err = toContent(geminiResponses.Candidates[0].Content)
	if err != nil {
		return nil, err
	}

	// 设置选项
	for _, candidate := range geminiResponses.Candidates[1:] {
		var responseOption *chat.Content
		responseOption, err = toContent(candidate.Content)
		if err != nil {
			return nil, err
		}
		responseOption.Type = chat.ContentTypeOption
		options = append(options, *responseOption)
	}
	response := &chat.Response{
		Answer:  *answer,
		Options: options,
	}

	return response, nil
}

func toGeminiContents(contents []chat.Content) ([]*genai.Content, error) {
	geminiContents := make([]*genai.Content, len(contents))
	for _, content := range contents {
		geminiContent, err := toGeminiContent(content)
		if err != nil {
			return nil, err
		}
		geminiContents = append(geminiContents, geminiContent)
	}
	return geminiContents, nil
}

func toGeminiContent(content chat.Content) (*genai.Content, error) {
	parts := make([]chat.Part, len(content.Parts))
	for _, part := range content.Parts {
		parts = append(parts, part.Part)
	}
	geminiParts, err := toGeminiParts(parts)
	if err != nil {
		return nil, err
	}

	switch content.Type {
	case chat.ContentTypeQuestion:
		return &genai.Content{
			Parts: geminiParts,
			Role:  "role",
		}, nil
	case chat.ContentTypeAnswer:
		return &genai.Content{
			Parts: geminiParts,
			Role:  "models",
		}, nil
	}

	err = &myerrors.CustomError{
		StatusCode:  200,
		MessageCode: "ECH0250",
		MessageText: fmt.Sprintf("不支持的Content类型：%s。", content.Type),
	}
	return nil, err
}

func toGeminiParts(parts []chat.Part) ([]genai.Part, error) {
	geminiParts := make([]genai.Part, len(parts))
	for _, part := range parts {
		geminiPart, err := toGeminiPart(part)
		if err != nil {
			return nil, err
		}
		geminiParts = append(geminiParts, geminiPart)
	}
	return geminiParts, nil
}

func toGeminiPart(part chat.Part) (genai.Part, error) {
	switch part.(type) {
	case *chat.TextPart:
		textPart := part.(*chat.TextPart)
		geminiPart := genai.Text(textPart.Text)
		return geminiPart, nil
	case *chat.ImagePart:
		imagePart := part.(*chat.ImagePart)
		mimeType, data, err := shared.ExtractBase64ImageData(imagePart.ImageUrl)
		if err != nil {
			err = &myerrors.CustomError{
				StatusCode:  200,
				MessageCode: "ECH0254",
				MessageText: "ImageUrl解析出错",
			}
			return nil, err
		}
		geminiPart := genai.Blob{MIMEType: mimeType, Data: data}
		return geminiPart, nil
	}
	err := &myerrors.CustomError{
		StatusCode:  200,
		MessageCode: "ECH0251",
		MessageText: fmt.Sprintf("不支持的Part类型：%s。", part.GetContentType()),
	}
	return nil, err
}

func toContent[T *genai.Content](geminiContent T) (*chat.Content, error) {
	return func(geminiContent *genai.Content) (*chat.Content, error) {
		contentType := ""
		switch geminiContent.Role {
		case "user":
			contentType = chat.ContentTypeQuestion
		case "models":
			contentType = chat.ContentTypeAnswer
		}
		if contentType == "" {
			err := &myerrors.CustomError{
				StatusCode:  200,
				MessageCode: "ECH0253",
				MessageText: fmt.Sprintf("不支持的GeminiContent类型：%s。", geminiContent.Role),
			}
			return nil, err
		}

		parts := make([]chat.PartWrapper, len(geminiContent.Parts))
		for _, geminiPart := range geminiContent.Parts {
			part, err := toPart(geminiPart)
			if err != nil {
				return nil, err
			}
			parts = append(parts, chat.PartWrapper{Part: part})
		}
		content := &chat.Content{
			Type:  contentType,
			Parts: parts,
		}
		return content, nil
	}(geminiContent)
}

func toPart[T genai.Part](geminiPart T) (chat.Part, error) {
	return func(geminiPart genai.Part) (chat.Part, error) {
		switch geminiPart.(type) {
		case *genai.Text:
			textPart := geminiPart.(*genai.Text)
			part := chat.TextPart{
				Type: chat.PartTypeText,
				Text: string(*textPart),
			}
			return &part, nil
		}
		err := &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "ECH0252",
			MessageText: fmt.Sprintf("不支持的GeminiPart类型：%s。", reflect.TypeOf(geminiPart).String()),
		}
		return nil, err
	}(geminiPart)
}
