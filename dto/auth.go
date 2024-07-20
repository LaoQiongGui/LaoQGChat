package dto

import "github.com/google/uuid"

type AuthDto struct {
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	LoginToken uuid.UUID `json:"login_token"`
	Permission string    `json:"permission"`
}
