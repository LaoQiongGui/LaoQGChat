package services

import (
	"LaoQGChat/api/models"
	"LaoQGChat/internal/myerrors"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type AuthService interface {
	Login(ctx *gin.Context, inDto models.AuthDto) *models.AuthDto
	Check(loginToken uuid.UUID) (*models.AuthDto, error)
}

type authService struct {
	getUserInfo           *sql.Stmt
	updateLoginStatus     *sql.Stmt
	getLoginStatusByToken *sql.Stmt
}

func NewAuthService(db *sql.DB) AuthService {
	var (
		err                   error
		getUserInfo           *sql.Stmt
		updateLoginStatus     *sql.Stmt
		getLoginStatusByToken *sql.Stmt
	)
	getUserInfo, err = db.Prepare(
		"SELECT password, permission FROM account WHERE user_name = $1")
	if err != nil {
		return nil
	}
	updateLoginStatus, err = db.Prepare(`
		INSERT INTO login_record (user_name, last_login_time, login_token)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_name)
		DO UPDATE SET
		    last_login_time = $2, login_token = $3`)
	if err != nil {
		return nil
	}
	getLoginStatusByToken, err = db.Prepare(
		"SELECT user_name, last_login_time FROM login_record WHERE login_token = $1")
	if err != nil {
		return nil
	}
	service := &authService{
		getUserInfo:           getUserInfo,
		updateLoginStatus:     updateLoginStatus,
		getLoginStatusByToken: getLoginStatusByToken,
	}
	return service
}

func (service *authService) Login(ctx *gin.Context, inDto models.AuthDto) *models.AuthDto {
	var (
		err         error
		password    string
		permission  string
		currentTime = time.Now()
		loginToken  = uuid.New()
	)
	err = service.getUserInfo.QueryRow(inDto.Username).Scan(&password, &permission)
	if err != nil || password != inDto.Password {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "EAU00",
			MessageText: "账号或密码错误。",
		}
		_ = ctx.Error(err)
		return nil
	}
	_, err = service.updateLoginStatus.Exec(inDto.Username, currentTime, loginToken)
	if err != nil {
		_ = ctx.Error(err)
		return nil
	}
	outDto := &models.AuthDto{
		LoginToken: loginToken,
		Permission: permission,
	}
	return outDto
}

func (service *authService) Check(loginToken uuid.UUID) (*models.AuthDto, error) {
	var (
		err           error
		userName      string
		password      string
		permission    string
		currentTime   = time.Now()
		lastLoginTime time.Time
	)
	// 用户存在check
	err = service.getLoginStatusByToken.QueryRow(loginToken).Scan(&userName, &lastLoginTime)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "EAU01",
			MessageText: "用户未登录。",
		}
		return nil, err
	}
	if currentTime.Sub(lastLoginTime).Hours() >= 24 {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "EAU02",
			MessageText: "登录已超时，请重新登录。",
		}
		return nil, err
	}
	err = service.getUserInfo.QueryRow(userName).Scan(&password, &permission)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "EAU03",
			MessageText: "用户已注销。",
		}
		return nil, err
	}

	outDto := &models.AuthDto{
		LoginToken: loginToken,
		Permission: permission,
	}
	return outDto, nil
}
