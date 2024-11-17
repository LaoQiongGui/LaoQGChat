package auth

import (
	"github.com/google/uuid"
	"time"
)

type UserInfo struct {
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	LoginToken uuid.UUID `json:"loginToken"`
	Permission string    `json:"permission"`
}

type LoginStatus struct {
	UserName      string    `json:"username"`
	LastLoginTime time.Time `json:"last_login_time"`
	LoginToken    uuid.UUID `json:"login_token"`
}
