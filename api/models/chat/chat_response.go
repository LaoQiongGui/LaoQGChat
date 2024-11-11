package chat

import (
	"github.com/google/uuid"
)

type Response struct {
	SessionId uuid.UUID `json:"sessionId"`
	Answer    string    `json:"answer"`
	Choices   []string  `json:"choices"`
}
