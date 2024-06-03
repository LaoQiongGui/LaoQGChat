package service

import (
	"LaoQGChat/dto"
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"os"
)

type ChatService interface {
	StartChat(ctx *gin.Context, inDto dto.ChatInDto) *dto.ChatOutDto
	Chat(ctx *gin.Context, inDto dto.ChatInDto) *dto.ChatOutDto
	EndChat(ctx *gin.Context, inDto dto.ChatInDto) *dto.ChatOutDto
}

type chatService struct {
	ChatContextMap      map[uuid.UUID]dto.ChatContext
	azureOpenAIKey      string
	modelDeploymentID   string
	azureOpenAIEndpoint string
}

func NewChatService() ChatService {
	service := &chatService{}
	service.ChatContextMap = make(map[uuid.UUID]dto.ChatContext)
	service.azureOpenAIKey = os.Getenv("AOAI_API_KEY")
	service.modelDeploymentID = os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL")
	service.azureOpenAIEndpoint = os.Getenv("AOAI_ENDPOINT")
	return service
}

func (service *chatService) StartChat(ctx *gin.Context, inDto dto.ChatInDto) *dto.ChatOutDto {
	outDto := dto.ChatOutDto{}

	keyCredential := azcore.NewKeyCredential(service.azureOpenAIKey)
	client, err := azopenai.NewClientWithKeyCredential(service.azureOpenAIEndpoint, keyCredential, nil)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "E0001"
		ctx.Keys["MessageText"] = "Azure OpenAI认证失败，请联系管理员。"
		panic(err)
	}

	messages := []azopenai.ChatRequestMessageClassification{
		&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(inDto.Question)},
	}
	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &service.modelDeploymentID,
	}, nil)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "E0002"
		ctx.Keys["MessageText"] = "Azure OpenAI获取答案失败，请联系管理员。"
		panic(err)
	}
	if resp.Choices == nil || len(resp.Choices) == 0 {
		ctx.Keys["StatusCode"] = 100
		ctx.Keys["MessageCode"] = "W0001"
		ctx.Keys["MessageText"] = "无法回答该问题。"
		panic(nil)
	}

	answer := *resp.Choices[0].Message.Content
	var choices = make([]string, 0)
	if len(resp.Choices) > 1 {
		for _, respChoice := range resp.Choices[1:] {
			choices = append(choices, *respChoice.Message.Content)
		}
	}

	sessionId := uuid.New()
	chatContext := dto.ChatContext{
		ChatMessages: append(messages, &azopenai.ChatRequestAssistantMessage{
			Content: to.Ptr(answer),
		}),
	}
	service.ChatContextMap[sessionId] = chatContext

	outDto.SessionId = sessionId
	outDto.Answer = answer
	outDto.Choices = choices
	return &outDto
}

func (service *chatService) Chat(ctx *gin.Context, inDto dto.ChatInDto) *dto.ChatOutDto {
	outDto := dto.ChatOutDto{}
	chatContext, ok := service.ChatContextMap[inDto.SessionId]
	if !ok {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "E0003"
		ctx.Keys["MessageText"] = "不存在该会话或该会话已被删除。"
		panic(nil)
	}

	keyCredential := azcore.NewKeyCredential(service.azureOpenAIKey)
	client, err := azopenai.NewClientWithKeyCredential(service.azureOpenAIEndpoint, keyCredential, nil)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "E0001"
		ctx.Keys["MessageText"] = "Azure OpenAI认证失败，请联系管理员。"
		panic(err)
	}

	messages := append(chatContext.ChatMessages, &azopenai.ChatRequestUserMessage{
		Content: azopenai.NewChatRequestUserMessageContent(inDto.Question),
	})
	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &service.modelDeploymentID,
	}, nil)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "E0002"
		ctx.Keys["MessageText"] = "Azure OpenAI获取答案失败，请联系管理员。"
		panic(err)
	}
	if resp.Choices == nil || len(resp.Choices) == 0 {
		ctx.Keys["StatusCode"] = 100
		ctx.Keys["MessageCode"] = "W0001"
		ctx.Keys["MessageText"] = "无法回答该问题。"
		panic(nil)
	}

	answer := *resp.Choices[0].Message.Content
	var choices = make([]string, 0)
	if len(resp.Choices) > 1 {
		for _, respChoice := range resp.Choices[1:] {
			choices = append(choices, *respChoice.Message.Content)
		}
	}
	chatContext.ChatMessages = append(
		messages, &azopenai.ChatRequestAssistantMessage{
			Content: to.Ptr(answer),
		})
	service.ChatContextMap[inDto.SessionId] = chatContext

	outDto.SessionId = inDto.SessionId
	outDto.Answer = answer
	outDto.Choices = choices
	return &outDto
}

func (service *chatService) EndChat(ctx *gin.Context, inDto dto.ChatInDto) *dto.ChatOutDto {
	_, ok := service.ChatContextMap[inDto.SessionId]
	if !ok {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "E0003"
		ctx.Keys["MessageText"] = "不存在该会话或该会话已被删除。"
		panic(nil)
	}
	delete(service.ChatContextMap, inDto.SessionId)
	return nil
}
