package chat

import (
	"github.com/google/uuid"
)

type Request struct {
	SessionId uuid.UUID        `json:"sessionId"`
	Model     string           `json:"model"`
	Contents  []ContentWrapper `json:"contents"`
}
