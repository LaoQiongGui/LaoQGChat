package dao

import (
	"LaoQGChat/api/models/auth"
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type AuthDao interface {
	GetUserInfoByUserName(userName string) (*auth.UserInfo, error)
	GetLoginStatusByToken(loginToken string) (*auth.LoginStatus, error)
	UpdateLoginStatus(entity auth.LoginStatus) error
}

type authDao struct {
	getUserInfoByUserNameStmt *sql.Stmt
	getLoginStatusStmtByToken *sql.Stmt
	updateLoginStatusStmt     *sql.Stmt
}

func NewAuthDao(db *sql.DB) (AuthDao, error) {
	var (
		getUserInfoByUserNameStmt *sql.Stmt
		getLoginStatusStmtByToken *sql.Stmt
		updateLoginStatusStmt     *sql.Stmt
		err                       error
	)
	getUserInfoByUserNameStmt, err = db.Prepare(`
		SELECT password, permission FROM account
		WHERE user_name = $1`)
	if err != nil {
		return nil, err
	}

	getLoginStatusStmtByToken, err = db.Prepare(`
		SELECT user_name, last_login_time FROM login_record
		WHERE login_token = $1`)
	if err != nil {
		return nil, err
	}

	updateLoginStatusStmt, err = db.Prepare(`
		INSERT INTO login_record (user_name, last_login_time, login_token)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_name)
		DO UPDATE SET last_login_time = $2, login_token = $3`)
	if err != nil {
		return nil, err
	}

	return authDao{
		getUserInfoByUserNameStmt: getUserInfoByUserNameStmt,
		getLoginStatusStmtByToken: getLoginStatusStmtByToken,
		updateLoginStatusStmt:     updateLoginStatusStmt,
	}, nil
}

func (dao authDao) GetUserInfoByUserName(userName string) (*auth.UserInfo, error) {
	var (
		password   string
		permission string
		err        error
	)

	row := dao.getUserInfoByUserNameStmt.QueryRow(userName)
	if err = row.Scan(&password, &permission); err != nil {
		return nil, err
	} else {
		return &auth.UserInfo{
			Username:   userName,
			Password:   password,
			Permission: permission,
		}, nil
	}
}

func (dao authDao) GetLoginStatusByToken(loginToken string) (*auth.LoginStatus, error) {
	var (
		userName      string
		lastLoginTime time.Time
		err           error
	)

	row := dao.getLoginStatusStmtByToken.QueryRow(loginToken)
	if err = row.Scan(&userName, &lastLoginTime); err != nil {
		return nil, err
	} else {
		return &auth.LoginStatus{
			UserName:      userName,
			LastLoginTime: lastLoginTime,
		}, nil
	}
}

func (dao authDao) UpdateLoginStatus(loginStatus auth.LoginStatus) error {
	_, err := dao.updateLoginStatusStmt.Exec(
		loginStatus.UserName, loginStatus.LastLoginTime, loginStatus.LoginToken.String())
	return err
}
