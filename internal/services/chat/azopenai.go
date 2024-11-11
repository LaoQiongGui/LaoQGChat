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
			MessageCode: "ECH01",
			MessageText: "Azure OpenAI认证失败，请联系管理员。",
		}
		return nil, err
	}

	// 设置默认模型
	defaultModel := os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL")

	return &azopenaiAPI{client: *azopenaiClient, defaultModel: defaultModel}, nil
}

func (api *azopenaiAPI) chat(ctx *gin.Context, model string, contexts []chat.Context) (*chat.Response, error) {
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
	messages, err = toAzopenaiContexts(contexts)
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
			MessageCode: "ECH02",
			MessageText: "Azure OpenAI获取答案失败，请联系管理员。",
		}
		return nil, err
	}

	// 设置回答
	answer := *azopenaiResponse.Choices[0].Message.Content
	choices := make([]string, len(azopenaiResponse.Choices)-1)
	for _, azopenaiChoice := range azopenaiResponse.Choices[1:] {
		choices = append(choices, *azopenaiChoice.Message.Content)
	}
	response := &chat.Response{
		Answer:  answer,
		Choices: choices,
	}

	return response, nil
}

func toAzopenaiContexts(contexts []chat.Context) ([]azopenai.ChatRequestMessageClassification, error) {
	azopenaiContexts := make([]azopenai.ChatRequestMessageClassification, len(contexts))
	for _, context := range contexts {
		azopenaiContext, err := toAzopenaiContext(&context)
		if err != nil {
			return nil, err
		}
		azopenaiContexts = append(azopenaiContexts, azopenaiContext)
	}
	return azopenaiContexts, nil
}

func toAzopenaiContext(context *chat.Context) (azopenai.ChatRequestMessageClassification, error) {
	contents := make([]chat.Content, len(context.Contents))
	for _, content := range context.Contents {
		contents = append(contents, content.Content)
	}
	azopenaiContents, err := toAzopenaiContents(contents)
	if err != nil {
		return nil, err
	}

	switch context.Type {
	case chat.ContextTypeQuestion:
		return &azopenai.ChatRequestUserMessage{
			Content: azopenai.NewChatRequestUserMessageContent(azopenaiContents),
		}, nil
	case chat.ContextTypeAnswer:
		if content, ok := context.Contents[0].Content.(*chat.TextContent); ok {
			return &azopenai.ChatRequestAssistantMessage{
				Content: &content.Text,
			}, nil
		} else {
			err = &myerrors.CustomError{
				StatusCode:  200,
				MessageCode: "ECH52",
				MessageText: fmt.Sprintf("不支持的Content类型：%s。", (context.Contents[0].Content).GetContentType()),
			}
			return nil, err
		}
	}
	err = &myerrors.CustomError{
		StatusCode:  200,
		MessageCode: "ECH51",
		MessageText: fmt.Sprintf("不支持的Context类型：%s。", context.Type),
	}
	return nil, nil
}

func toAzopenaiContents(contents []chat.Content) ([]azopenai.ChatCompletionRequestMessageContentPartClassification, error) {
	azopenaiContents := make([]azopenai.ChatCompletionRequestMessageContentPartClassification, len(contents))
	for _, content := range contents {
		azopenaiContent, err := toAzopenaiContent(&content)
		if err != nil {
			return nil, err
		}
		azopenaiContents = append(azopenaiContents, azopenaiContent)
	}
	return azopenaiContents, nil
}

func toAzopenaiContent(content *chat.Content) (azopenai.ChatCompletionRequestMessageContentPartClassification, error) {
	switch (*content).(type) {
	case *chat.TextContent:
		textContent := (*content).(*chat.TextContent)
		azopenaiContent := &azopenai.ChatCompletionRequestMessageContentPartText{
			Text: &textContent.Text,
		}
		return azopenaiContent, nil
	case *chat.ImageContent:
		imageContent := (*content).(*chat.ImageContent)
		azopenaiContent := &azopenai.ChatCompletionRequestMessageContentPartImage{
			ImageURL: &azopenai.ChatCompletionRequestMessageContentPartImageURL{
				URL: &imageContent.ImageUrl,
			},
		}
		return azopenaiContent, nil
	}
	err := myerrors.CustomError{
		StatusCode:  200,
		MessageCode: "ECH52",
		MessageText: fmt.Sprintf("不支持的Content类型：%s。", (*content).GetContentType()),
	}
	return nil, &err
}
