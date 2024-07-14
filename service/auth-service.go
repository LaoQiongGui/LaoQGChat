package service

import (
	"LaoQGChat/dto"
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"time"
)

type AuthService interface {
	Login(ctx *gin.Context, inDto dto.AuthInDto) *dto.AuthOutDto
	Check(ctx *gin.Context, permissionCheck func(permission string))
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

func (service *authService) Login(ctx *gin.Context, inDto dto.AuthInDto) *dto.AuthOutDto {
	var (
		err         error
		password    string
		permission  string
		currentTime = time.Now()
		loginToken  = uuid.New()
	)
	err = service.getUserInfo.QueryRow(inDto.Username).Scan(&password, &permission)
	if err != nil || password != inDto.Password {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "EAU00"
		ctx.Keys["MessageText"] = "账号或密码错误。"
		panic(err)
	}
	_, err = service.updateLoginStatus.Exec(inDto.Username, currentTime, loginToken)
	if err != nil {
		panic(err)
	}
	outDto := &dto.AuthOutDto{
		LoginToken: loginToken,
	}
	return outDto
}

func (service *authService) Check(ctx *gin.Context, permissionCheck func(permission string)) {
	var (
		err           error
		loginToken    uuid.UUID
		userName      string
		password      string
		permission    string
		currentTime   = time.Now()
		lastLoginTime time.Time
	)
	// 从http头中取得loginToken
	loginToken, err = uuid.Parse(ctx.GetHeader("LoginToken"))
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "EAU01"
		ctx.Keys["MessageText"] = "用户未登录。"
		panic(err)
	}
	// 用户存在check
	err = service.getLoginStatusByToken.QueryRow(loginToken).Scan(&userName, &lastLoginTime)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "EAU01"
		ctx.Keys["MessageText"] = "用户未登录。"
		panic(err)
	}
	if currentTime.Sub(lastLoginTime).Hours() >= 24 {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "EAU02"
		ctx.Keys["MessageText"] = "登录已超时，请重新登录。"
		panic(err)
	}
	err = service.getUserInfo.QueryRow(userName).Scan(&password, &permission)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "EAU03"
		ctx.Keys["MessageText"] = "用户已注销。"
		panic(err)
	}
	// 不需要权限check直接返回
	if permissionCheck == nil {
		return
	}
	// 权限check
	permissionCheck(permission)
}
