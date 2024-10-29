package services

import (
	"LaoQGChat/api/models"
	"LaoQGChat/internal/myerrors"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type ChatService interface {
	StartChat(ctx *gin.Context, inDto models.ChatInDto) *models.ChatOutDto
	Chat(ctx *gin.Context, inDto models.ChatInDto) *models.ChatOutDto
	EndChat(ctx *gin.Context, inDto models.ChatInDto) *models.ChatOutDto
}

type chatService struct {
	azureOpenAIKey      string
	modelDeploymentID   string
	azureOpenAIEndpoint string

	getAllChatContexts  *sql.Stmt
	getUserChatContexts *sql.Stmt
	getChatContextById  *sql.Stmt
	insertChatContext   *sql.Stmt
	updateChatContext   *sql.Stmt
	deleteChatContext   *sql.Stmt
}

func NewChatService(db *sql.DB) ChatService {
	var (
		err                 error
		getAllChatContexts  *sql.Stmt
		getUserChatContexts *sql.Stmt
		getChatContextById  *sql.Stmt
		insertChatContext   *sql.Stmt
		updateChatContext   *sql.Stmt
		deleteChatContext   *sql.Stmt
	)

	getAllChatContexts, err = db.Prepare(`
		SELECT user_name, session_id
		FROM chat_record
		ORDER BY user_name, create_timestamp`)
	if err != nil {
		return nil
	}

	getUserChatContexts, err = db.Prepare(`
		SELECT session_id
		FROM chat_record
		WHERE user_name = $1
		ORDER BY create_timestamp`)
	if err != nil {
		return nil
	}

	getChatContextById, err = db.Prepare(`
		SELECT context
		FROM chat_record
		WHERE session_id = $1`)
	if err != nil {
		return nil
	}

	insertChatContext, err = db.Prepare(`
		INSERT INTO chat_record
		(user_name, session_id, context, create_timestamp, update_timestamp)
		VALUES ($1, $2, $3, $4, $4)`)
	if err != nil {
		return nil
	}

	updateChatContext, err = db.Prepare(`
		UPDATE chat_record
		SET context = $2, update_timestamp = $3
		WHERE session_id = $1`)
	if err != nil {
		return nil
	}

	deleteChatContext, err = db.Prepare(`
		DELETE FROM chat_record
		WHERE session_id = $1`)
	if err != nil {
		return nil
	}

	service := &chatService{
		azureOpenAIKey:      os.Getenv("AOAI_API_KEY"),
		modelDeploymentID:   os.Getenv("AOAI_CHAT_COMPLETIONS_MODEL"),
		azureOpenAIEndpoint: os.Getenv("AOAI_ENDPOINT"),
		getAllChatContexts:  getAllChatContexts,
		getUserChatContexts: getUserChatContexts,
		getChatContextById:  getChatContextById,
		insertChatContext:   insertChatContext,
		updateChatContext:   updateChatContext,
		deleteChatContext:   deleteChatContext,
	}
	return service
}

func (service *chatService) StartChat(ctx *gin.Context, inDto models.ChatInDto) *models.ChatOutDto {
	fmt.Println("Calling service StartChat")
	defer fmt.Println("Returning service StartChat")
	var (
		userName       = ctx.GetString("UserName")
		err            error
		keyCredential  = azcore.NewKeyCredential(service.azureOpenAIKey)
		currentTime    = time.Now()
		sessionId      = uuid.New()
		chatContextStr []byte
		chatContext    models.ChatContext
		outDto         = new(models.ChatOutDto)
	)

	// azopenai认证
	client, err := azopenai.NewClientWithKeyCredential(service.azureOpenAIEndpoint, keyCredential, nil)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "ECH01",
			MessageText: "Azure OpenAI认证失败，请联系管理员。",
		}
		_ = ctx.Error(err)
		return nil
	}

	// 将inDto转为azopenai的输入
	chatRequestUserMessage := inDto.ToAzopenai()
	messages := []azopenai.ChatRequestMessageClassification{
		&chatRequestUserMessage,
	}

	// 发送azopenai请求
	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &service.modelDeploymentID,
	}, nil)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "ECH02",
			MessageText: "Azure OpenAI获取答案失败，请联系管理员。",
		}
		_ = ctx.Error(err)
		return nil
	}
	if resp.Choices == nil || len(resp.Choices) == 0 {
		err = &myerrors.CustomError{
			StatusCode:  100,
			MessageCode: "WCH01",
			MessageText: "无法回答该问题。",
		}
		_ = ctx.Error(err)
		return nil
	}

	// 设置回答
	answer := *resp.Choices[0].Message.Content
	var choices = make([]string, 0)
	if len(resp.Choices) > 1 {
		for _, respChoice := range resp.Choices[1:] {
			choices = append(choices, *respChoice.Message.Content)
		}
	}

	// TODO：设置选项

	//
	chatContext = models.ChatContext{
		ChatMessages: append(messages, &azopenai.ChatRequestAssistantMessage{
			Content: to.Ptr(answer),
		}),
	}

	// 插入对话上下文
	chatContextStr, err = json.Marshal(chatContext)
	if err != nil {
		_ = ctx.Error(err)
		return nil
	}
	_, err = service.insertChatContext.Exec(userName, sessionId, chatContextStr, currentTime)
	if err != nil {
		_ = ctx.Error(err)
		return nil
	}

	outDto.SessionId = sessionId
	outDto.Answer = answer
	outDto.Choices = choices
	return outDto
}

func (service *chatService) Chat(ctx *gin.Context, inDto models.ChatInDto) *models.ChatOutDto {
	fmt.Println("Calling service ChatInDto")
	defer fmt.Println("Returning service ChatInDto")
	var (
		userName       = ctx.GetString("UserName")
		permission     = ctx.GetString("Permission")
		err            error
		keyCredential  = azcore.NewKeyCredential(service.azureOpenAIKey)
		currentTime    = time.Now()
		sessionId      uuid.UUID
		chatContextStr []byte
		chatContext    models.ChatContext
		outDto         = new(models.ChatOutDto)
	)

	// 非管理员用户检测SessionId是否在自己的对话记录中
	if permission != "super" {
		rows, err := service.getUserChatContexts.Query(userName)
		if err != nil {
			err = &myerrors.CustomError{
				StatusCode:  200,
				MessageCode: "ECH03",
				MessageText: "不存在该会话或该会话已被删除。",
			}
			_ = ctx.Error(err)
			return nil
		}
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "ECH03",
			MessageText: "不存在该会话或该会话已被删除。",
		}
		for rows.Next() {
			if rows.Scan(&sessionId) != nil {
				_ = ctx.Error(err)
				return nil
			}
			if sessionId == inDto.SessionId {
				err = nil
				break
			}
		}
		if err != nil {
			_ = ctx.Error(err)
			return nil
		}
	}

	// 获取对话上下文
	err = service.getChatContextById.QueryRow(inDto.SessionId).Scan(&chatContextStr)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "ECH03",
			MessageText: "不存在该会话或该会话已被删除。",
		}
		_ = ctx.Error(err)
		return nil
	}

	err = json.Unmarshal(chatContextStr, &chatContext)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  990,
			MessageCode: "ECH91",
			MessageText: "JSON反序列化失败。",
		}
		_ = ctx.Error(err)
		return nil
	}

	// azopenai认证
	client, err := azopenai.NewClientWithKeyCredential(service.azureOpenAIEndpoint, keyCredential, nil)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "ECH01",
			MessageText: "Azure OpenAI认证失败，请联系管理员。",
		}
		_ = ctx.Error(err)
		return nil
	}

	// 将inDto转为azopenai的输入
	chatRequestUserMessage := inDto.ToAzopenai()
	messages := []azopenai.ChatRequestMessageClassification{
		&chatRequestUserMessage,
	}

	// 将转换后的inDto拼接在原回答之后
	messages = append(chatContext.ChatMessages, messages...)
	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &service.modelDeploymentID,
	}, nil)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "ECH02",
			MessageText: "Azure OpenAI获取答案失败，请联系管理员。",
		}
		_ = ctx.Error(err)
		return nil
	}
	if resp.Choices == nil || len(resp.Choices) == 0 {
		err = &myerrors.CustomError{
			StatusCode:  100,
			MessageCode: "WCH01",
			MessageText: "无法回答该问题。",
		}
		_ = ctx.Error(err)
		return nil
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

	// 更新对话上下文
	chatContextStr, err = json.Marshal(chatContext)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  990,
			MessageCode: "ECH90",
			MessageText: "JSON序列化失败。",
		}
		_ = ctx.Error(err)
		return nil
	}
	_, err = service.updateChatContext.Exec(inDto.SessionId, chatContextStr, currentTime)
	if err != nil {
		_ = ctx.Error(err)
		return nil
	}

	outDto.SessionId = inDto.SessionId
	outDto.Answer = answer
	outDto.Choices = choices
	return outDto
}

func (service *chatService) EndChat(ctx *gin.Context, inDto models.ChatInDto) *models.ChatOutDto {
	fmt.Println("Calling service EndChat")
	defer fmt.Println("Returning service EndChat")
	var (
		userName   = ctx.GetString("UserName")
		permission = ctx.GetString("Permission")
		err        error
		sessionId  uuid.UUID
	)

	// 非管理员用户检测SessionId是否在自己的对话记录中
	if permission != "super" {
		rows, err := service.getUserChatContexts.Query(userName)
		if err != nil {
			err = &myerrors.CustomError{
				StatusCode:  200,
				MessageCode: "ECH03",
				MessageText: "不存在该会话或该会话已被删除。",
			}
			_ = ctx.Error(err)
		}
		err = errors.New("")
		for rows.Next() {
			if rows.Scan(&sessionId) != nil {
				_ = ctx.Error(err)
			}
			if sessionId == inDto.SessionId {
				err = nil
				break
			}
		}
		if err != nil {
			err = &myerrors.CustomError{
				StatusCode:  200,
				MessageCode: "ECH03",
				MessageText: "不存在该会话或该会话已被删除。",
			}
			_ = ctx.Error(err)
			return nil
		}
	}

	// 删除对话上下文
	_, err = service.deleteChatContext.Exec(inDto.SessionId)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "ECH03",
			MessageText: "不存在该会话或该会话已被删除。",
		}
		_ = ctx.Error(err)
		return nil
	}
	return nil
}
