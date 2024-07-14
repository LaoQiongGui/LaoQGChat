package dto

import (
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/google/uuid"
)

type ChatInDto struct {
	Permission string    `json:"permission"`
	SessionId  uuid.UUID `json:"sessionId"`
	Question   string    `json:"question"`
}

type ChatOutDto struct {
	SessionId uuid.UUID `json:"sessionId"`
	Answer    string    `json:"answer"`
	Choices   []string  `json:"choices"`
}

type ChatContext struct {
	ChatMessages []azopenai.ChatRequestMessageClassification
}
