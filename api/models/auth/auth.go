package auth

import "github.com/google/uuid"

type Entity struct {
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	LoginToken uuid.UUID `json:"loginToken"`
	Permission string    `json:"permission"`
}
