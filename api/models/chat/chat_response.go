package chat

import (
	"github.com/google/uuid"
)

type Response struct {
	SessionId uuid.UUID `json:"sessionId"`
	Answer    Content   `json:"answer"`
	Options   []Content `json:"options"`
}
