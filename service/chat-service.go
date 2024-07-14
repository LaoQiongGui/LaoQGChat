package service

import (
	"LaoQGChat/dto"
	"context"
	"database/sql"
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
	chatContextMap      map[uuid.UUID]dto.ChatContext
	azureOpenAIKey      string
	modelDeploymentID   string
	azureOpenAIEndpoint string
}

func NewChatService(db *sql.DB) ChatService {
	service := &chatService{
		chatContextMap:      make(map[uuid.UUID]dto.ChatContext),
		azureOpenAIKey:      os.Getenv("AOAI_API_KEY"),
		modelDeploymentID:   os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL"),
		azureOpenAIEndpoint: os.Getenv("AOAI_ENDPOINT"),
	}
	return service
}

func (service *chatService) StartChat(ctx *gin.Context, inDto dto.ChatInDto) *dto.ChatOutDto {
	outDto := new(dto.ChatOutDto)

	keyCredential := azcore.NewKeyCredential(service.azureOpenAIKey)
	client, err := azopenai.NewClientWithKeyCredential(service.azureOpenAIEndpoint, keyCredential, nil)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "ECH01"
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
		ctx.Keys["MessageCode"] = "ECH02"
		ctx.Keys["MessageText"] = "Azure OpenAI获取答案失败，请联系管理员。"
		panic(err)
	}
	if resp.Choices == nil || len(resp.Choices) == 0 {
		ctx.Keys["StatusCode"] = 100
		ctx.Keys["MessageCode"] = "WCH01"
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
	service.chatContextMap[sessionId] = chatContext

	outDto.SessionId = sessionId
	outDto.Answer = answer
	outDto.Choices = choices
	return outDto
}

func (service *chatService) Chat(ctx *gin.Context, inDto dto.ChatInDto) *dto.ChatOutDto {
	outDto := new(dto.ChatOutDto)
	chatContext, ok := service.chatContextMap[inDto.SessionId]
	if !ok {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "ECH03"
		ctx.Keys["MessageText"] = "不存在该会话或该会话已被删除。"
		panic(nil)
	}

	keyCredential := azcore.NewKeyCredential(service.azureOpenAIKey)
	client, err := azopenai.NewClientWithKeyCredential(service.azureOpenAIEndpoint, keyCredential, nil)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "ECH01"
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
		ctx.Keys["MessageCode"] = "ECH02"
		ctx.Keys["MessageText"] = "Azure OpenAI获取答案失败，请联系管理员。"
		panic(err)
	}
	if resp.Choices == nil || len(resp.Choices) == 0 {
		ctx.Keys["StatusCode"] = 100
		ctx.Keys["MessageCode"] = "WCH01"
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
	service.chatContextMap[inDto.SessionId] = chatContext

	outDto.SessionId = inDto.SessionId
	outDto.Answer = answer
	outDto.Choices = choices
	return outDto
}

func (service *chatService) EndChat(ctx *gin.Context, inDto dto.ChatInDto) *dto.ChatOutDto {
	_, ok := service.chatContextMap[inDto.SessionId]
	if !ok {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "ECH03"
		ctx.Keys["MessageText"] = "不存在该会话或该会话已被删除。"
		panic(nil)
	}
	delete(service.chatContextMap, inDto.SessionId)
	return nil
}
