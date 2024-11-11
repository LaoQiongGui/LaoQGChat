package chat

import (
	"LaoQGChat/api/models/chat"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Service interface {
	Chat(ctx *gin.Context, request chat.Request) *chat.Response
	EndChat(ctx *gin.Context, request chat.Request) *chat.Response
}

type chatService struct{}

func NewService() Service {
	service := &chatService{}
	return service
}

func (service *chatService) Chat(ctx *gin.Context, request chat.Request) *chat.Response {
	var (
		sessionId    = request.SessionId
		chatContexts []chat.Context
		response     *chat.Response
		err          error
	)

	if sessionId != uuid.Nil {
		// 有sessionId则获取历史对话
		chatContexts, err = service.getSessionContexts(ctx, sessionId)
		if err != nil {
			_ = ctx.Error(err)
			return nil
		}
	} else {
		// 无sessionId则创建新的对话
		sessionId = uuid.New()
	}

	// 调用外部API
	response, err = service.callExternalChatAPI(ctx, request.Model, chatContexts)
	if err != nil {
		_ = ctx.Error(err)
		return nil
	}
	response.SessionId = sessionId

	// 保存上下文
	err = service.saveContext(ctx, chatContexts)
	if err != nil {
		_ = ctx.Error(err)
		return nil
	}

	return response
}

func (service *chatService) EndChat(ctx *gin.Context, request chat.Request) *chat.Response {
	// TODO
	return &chat.Response{}
}

func (service *chatService) getSessionContexts(ctx *gin.Context, sessionId uuid.UUID) ([]chat.Context, error) {
	// TODO
	chatContexts := make([]chat.Context, 0)
	return chatContexts, nil
}

func (service *chatService) callExternalChatAPI(ctx *gin.Context, model string, chatContexts []chat.Context) (*chat.Response, error) {
	externalAPI, err := getExternalAPI(model)
	if err != nil {
		return nil, err
	}
	return externalAPI.chat(ctx, model, chatContexts)
}

func (service *chatService) saveContext(ctx *gin.Context, chatContexts []chat.Context) error {
	// TODO
	return nil
}
