package chat

import (
	"LaoQGChat/api/models/chat"
	"LaoQGChat/internal/dao"
	"LaoQGChat/internal/myerrors"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Service interface {
	Chat(ctx *gin.Context, request chat.Request) *chat.Response
	EndChat(ctx *gin.Context, request chat.Request) *chat.Response
}

type chatService struct {
	chatDao dao.ChatDao
}

func NewService(db *sql.DB) Service {
	chatDao, err := dao.NewChatDao(db)
	if err != nil {
		return nil
	}
	service := &chatService{chatDao: chatDao}
	return service
}

func (service *chatService) Chat(ctx *gin.Context, request chat.Request) *chat.Response {
	var (
		sessionId uuid.UUID
		contents  = make([]chat.Content, 0)
		response  *chat.Response
		err       error
	)

	if request.SessionId != nil {
		sessionId = *request.SessionId

		// 有sessionId则获取历史对话
		contents, err = service.getSessionContents(ctx, sessionId)
		if err != nil {
			_ = ctx.Error(err)
			return nil
		}
	} else {
		// 无sessionId则创建新的对话
		sessionId = uuid.New()
	}

	// 加入本次提问
	questionContent := chat.Content{
		Type:  chat.ContentTypeQuestion,
		Parts: request.Question,
	}
	contents = append(contents, questionContent)

	// 调用外部API
	response, err = service.callExternalChatAPI(ctx, request.Model, contents)
	if err != nil {
		_ = ctx.Error(err)
		return nil
	}
	response.SessionId = sessionId
	answerContent := chat.Content{
		Type:  chat.ContentTypeAnswer,
		Parts: response.Answer,
	}

	// 保存本次对话
	err = service.saveSessionContents(ctx, sessionId, []chat.Content{
		questionContent,
		answerContent,
	})
	if err != nil {
		_ = ctx.Error(err)
		return nil
	}

	return response
}

func (service *chatService) EndChat(ctx *gin.Context, request chat.Request) *chat.Response {
	return &chat.Response{}
}

func (service *chatService) getSessionContents(ctx *gin.Context, sessionId uuid.UUID) ([]chat.Content, error) {
	var (
		userName   string
		permission string
		contents   []chat.Content
		ok         bool
		err        error
	)
	value, exists := ctx.Get("UserName")
	if !exists {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "ECH0000",
			MessageText: "用户权限不足。",
		}
		return nil, err
	}
	userName = value.(string)

	value, exists = ctx.Get("Permission")
	if !exists {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "ECH0000",
			MessageText: "用户权限不足。",
		}
		return nil, err
	}
	permission = value.(string)

	// 根据用户权限检测会话
	if permission == "super" {
		// 超级用户检测当前会话是否存在
		ok, err = service.chatDao.CheckSessionById(sessionId.String())
		if err != nil {
			return nil, err
		} else if !ok {
			err = &myerrors.CustomError{
				StatusCode:  200,
				MessageCode: "ECH0001",
				MessageText: "当前会话不存在。",
			}
			return nil, err
		}
	} else {
		// 普通用户检测当前会话是否处于自己的会话列表中
		ok, err = service.chatDao.CheckUserSessionById(userName, sessionId.String())
		if err != nil {
			return nil, err
		} else if !ok {
			err = &myerrors.CustomError{
				StatusCode:  200,
				MessageCode: "ECH0000",
				MessageText: "用户权限不足。",
			}
			return nil, err
		}
	}

	// 获取会话内容
	contents, err = service.chatDao.GetSessionContentsById(sessionId.String())
	if err != nil {
		return nil, err
	}
	return contents, nil
}

func (service *chatService) callExternalChatAPI(ctx *gin.Context, model string, contents []chat.Content) (*chat.Response, error) {
	api, err := getExternalAPI(model)
	if err != nil {
		return nil, err
	}
	return api.chat(ctx, model, contents)
}

func (service *chatService) saveSessionContents(ctx *gin.Context, sessionId uuid.UUID, contents []chat.Content) error {
	err := service.chatDao.InsertSessionContents(sessionId.String(), contents)
	return err
}
