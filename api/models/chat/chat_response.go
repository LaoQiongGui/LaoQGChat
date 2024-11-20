package chat

import (
	"github.com/google/uuid"
)

type Response struct {
	SessionId uuid.UUID       `json:"sessionId"`
	Answer    []PartWrapper   `json:"answer"`
	Options   [][]PartWrapper `json:"options"`
}
