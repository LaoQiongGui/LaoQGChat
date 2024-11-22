package auth

import (
	model "LaoQGChat/api/models/auth"
	"LaoQGChat/internal/dao"
	"LaoQGChat/internal/myerrors"
	"database/sql"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Service interface {
	Login(ctx *gin.Context, request model.UserInfo) *model.UserInfo
	Check(loginToken uuid.UUID) (*model.UserInfo, error)
}

type authService struct {
	authDao dao.AuthDao
}

func NewService(db *sql.DB) Service {
	authDao, err := dao.NewAuthDao(db)
	if err != nil {
		return nil
	}
	service := &authService{authDao: authDao}
	return service
}

func (service *authService) Login(ctx *gin.Context, request model.UserInfo) *model.UserInfo {
	var (
		permission  string
		currentTime = time.Now()
		loginToken  = uuid.New()
		userInfo    *model.UserInfo
		loginStatus *model.LoginStatus
		err         error
	)
	// 验证账号密码
	userInfo, err = service.authDao.GetUserInfoByUserName(request.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = &myerrors.CustomError{
				StatusCode:  200,
				MessageCode: "EAU0000",
				MessageText: "账号或密码错误。",
			}
		}
		_ = ctx.Error(err)
		return nil
	}
	if userInfo.Password != request.Password {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "EAU0000",
			MessageText: "账号或密码错误。",
		}
		_ = ctx.Error(err)
		return nil
	}

	// 更新登录凭证
	loginStatus = &model.LoginStatus{
		UserName:      request.Username,
		LastLoginTime: currentTime,
		LoginToken:    loginToken,
	}
	err = service.authDao.UpdateLoginStatus(*loginStatus)
	if err != nil {
		_ = ctx.Error(err)
		return nil
	}

	// 返回登录凭证与账号权限
	outDto := &model.UserInfo{
		LoginToken: loginToken,
		Permission: permission,
	}
	return outDto
}

func (service *authService) Check(loginToken uuid.UUID) (*model.UserInfo, error) {
	var (
		currentTime = time.Now()
		userInfo    *model.UserInfo
		loginStatus *model.LoginStatus
		err         error
	)

	// 登录状态检测
	loginStatus, err = service.authDao.GetLoginStatusByToken(loginToken.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = &myerrors.CustomError{
				StatusCode:  200,
				MessageCode: "EAU0001",
				MessageText: "用户未登录。",
			}
		}
		return nil, err
	}

	// 登陆时间检测
	if currentTime.Sub(loginStatus.LastLoginTime).Hours() > 24 {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "EAU0002",
			MessageText: "登录已超时，请重新登录。",
		}
		return nil, err
	}

	userInfo, err = service.authDao.GetUserInfoByUserName(loginStatus.UserName)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "EAU0003",
			MessageText: "用户已注销。",
		}
		return nil, err
	}

	outDto := &model.UserInfo{
		LoginToken: loginToken,
		Permission: userInfo.Permission,
	}
	return outDto, nil
}
