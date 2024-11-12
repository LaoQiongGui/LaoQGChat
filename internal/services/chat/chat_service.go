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
		sessionId   = request.SessionId
		chatContext = make([]chat.Content, 0)
		response    *chat.Response
		err         error
	)

	if sessionId != uuid.Nil {
		// 有sessionId则获取历史对话
		chatContext, err = service.getSessionContext(ctx, sessionId)
		if err != nil {
			_ = ctx.Error(err)
			return nil
		}
	} else {
		// 无sessionId则创建新的对话
		sessionId = uuid.New()
	}

	// 加入本次提问
	chatContext = append(chatContext, request.Question)

	// 调用外部API
	response, err = service.callExternalChatAPI(ctx, request.Model, chatContext)
	if err != nil {
		_ = ctx.Error(err)
		return nil
	}
	response.SessionId = sessionId

	// 保存上下文
	err = service.saveSessionContext(ctx, chatContext)
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

func (service *chatService) getSessionContext(ctx *gin.Context, sessionId uuid.UUID) ([]chat.Content, error) {
	// TODO
	chatContext := make([]chat.Content, 0)
	return chatContext, nil
}

func (service *chatService) callExternalChatAPI(ctx *gin.Context, model string, chatContexts []chat.Content) (*chat.Response, error) {
	externalAPI, err := getExternalAPI(model)
	if err != nil {
		return nil, err
	}
	return externalAPI.chat(ctx, model, chatContexts)
}

func (service *chatService) saveSessionContext(ctx *gin.Context, chatContents []chat.Content) error {
	// TODO
	return nil
}
