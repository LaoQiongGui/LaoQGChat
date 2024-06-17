package dto

import "github.com/google/uuid"

type AuthInDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthOutDto struct {
	AuthToken uuid.UUID `json:"auth_token"`
}
