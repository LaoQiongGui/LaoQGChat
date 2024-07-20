package dto

import (
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/google/uuid"
)

type ChatInDto struct {
	SessionId uuid.UUID `json:"sessionId"`
	Question  string    `json:"question"`
}

type ChatOutDto struct {
	SessionId uuid.UUID `json:"sessionId"`
	Answer    string    `json:"answer"`
	Choices   []string  `json:"choices"`
}

type ChatContext struct {
	ChatMessages []azopenai.ChatRequestMessageClassification `json:"chatMessages"`
}

func (chatContext *ChatContext) UnmarshalJSON(data []byte) error {
	var (
		err          error
		chatMessage  azopenai.ChatRequestMessageClassification
		chatMessages []azopenai.ChatRequestMessageClassification
	)
	chatStrContext := struct {
		ChatMessageRaw []json.RawMessage `json:"chatMessages"`
	}{}
	if err = json.Unmarshal(data, &chatStrContext); err != nil {
		return err
	}
	for _, chatMessageRaw := range chatStrContext.ChatMessageRaw {
		chatRequestMessage := struct {
			Role azopenai.ChatRole `json:"role"`
		}{}
		if err = json.Unmarshal(chatMessageRaw, &chatRequestMessage); err != nil {
			return err
		}
		switch chatRequestMessage.Role {
		case azopenai.ChatRoleAssistant:
			chatMessage = &azopenai.ChatRequestAssistantMessage{}
		case azopenai.ChatRoleFunction:
			chatMessage = &azopenai.ChatRequestFunctionMessage{}
		case azopenai.ChatRoleSystem:
			chatMessage = &azopenai.ChatRequestSystemMessage{}
		case azopenai.ChatRoleTool:
			chatMessage = &azopenai.ChatRequestToolMessage{}
		case azopenai.ChatRoleUser:
			chatMessage = &azopenai.ChatRequestUserMessage{}
		default:
			chatMessage = &azopenai.ChatRequestMessage{}
		}
		if err = json.Unmarshal(chatMessageRaw, chatMessage); err != nil {
			return err
		}
		chatMessages = append(chatMessages, chatMessage)
	}
	chatContext.ChatMessages = chatMessages
	return nil
}
