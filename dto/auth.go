package dto

import "github.com/google/uuid"

type AuthInDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthOutDto struct {
	LoginToken uuid.UUID `json:"login_token"`
}
