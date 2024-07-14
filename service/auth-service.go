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
	Check(ctx *gin.Context, tokenId uuid.UUID)
}

type authService struct {
	selectUserInfo        *sql.Stmt
	updateLoginStatus     *sql.Stmt
	getLoginStatusByToken *sql.Stmt
}

func NewAuthService(db *sql.DB) AuthService {
	var (
		err                   error
		selectUserInfo        *sql.Stmt
		updateLoginStatus     *sql.Stmt
		getLoginStatusByToken *sql.Stmt
	)
	selectUserInfo, err = db.Prepare(
		"SELECT permission FROM account WHERE user_name = $1 AND password = $2")
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
		selectUserInfo:        selectUserInfo,
		updateLoginStatus:     updateLoginStatus,
		getLoginStatusByToken: getLoginStatusByToken,
	}
	return service
}

func (service *authService) Login(ctx *gin.Context, inDto dto.AuthInDto) *dto.AuthOutDto {
	var (
		err         error
		permission  string
		currentTime time.Time = time.Now()
		loginToken  uuid.UUID = uuid.New()
	)
	err = service.selectUserInfo.QueryRow(inDto.Username, inDto.Password).Scan(&permission)
	if err != nil {
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

func (service *authService) Check(ctx *gin.Context, loginToken uuid.UUID) {
	var (
		err           error
		userName      string
		currentTime   time.Time = time.Now()
		lastLoginTime time.Time
	)
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
}
