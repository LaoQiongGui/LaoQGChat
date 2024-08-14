package dto

import (
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/google/uuid"
)

type ChatInDto struct {
	SessionId uuid.UUID                     `json:"sessionId"`
	Contents  []ChatQuestionContentPartsDto `json:"contents"`
}

func (chatInDto *ChatInDto) UnmarshalJSON(data []byte) error {
	var (
		err         error
		typeContent ChatQuestionContentPartsDtoType
	)

	chatTypeInDto := struct {
		SessionId uuid.UUID         `json:"sessionId"`
		Contents  []json.RawMessage `json:"contents"`
	}{}
	if err = json.Unmarshal(data, &chatTypeInDto); err != nil {
		return err
	}
	chatInDto.SessionId = chatTypeInDto.SessionId
	for _, content := range chatTypeInDto.Contents {
		if err = json.Unmarshal(content, &typeContent); err != nil {
			return err
		}
		switch typeContent.Type {
		case "Text":
			var textContent ChatQuestionContentPartsDtoText
			_ = json.Unmarshal(content, &textContent)
			chatInDto.Contents = append(chatInDto.Contents, &textContent)
		case "Image":
			var imageContent ChatQuestionContentPartsDtoImage
			_ = json.Unmarshal(content, &imageContent)
			chatInDto.Contents = append(chatInDto.Contents, &imageContent)
		case "Audio":
			var audioContent ChatQuestionContentPartsDtoAudio
			_ = json.Unmarshal(content, &audioContent)
			chatInDto.Contents = append(chatInDto.Contents, &audioContent)
		case "ImageOCR":
			var imageContent ChatQuestionContentPartsDtoImageOCR
			_ = json.Unmarshal(content, &imageContent)
			chatInDto.Contents = append(chatInDto.Contents, &imageContent)
		}
	}
	return nil
}

func (chatInDto *ChatInDto) ToAzopenai() azopenai.ChatRequestUserMessage {
	var azopenaiContents []azopenai.ChatCompletionRequestMessageContentPartClassification
	for _, content := range chatInDto.Contents {
		azopenaiContent, _ := content.ToChatCompletionRequestMessageContentPartClassification()
		azopenaiContents = append(azopenaiContents, azopenaiContent)
	}
	chatRequestUserMessage := azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(azopenaiContents)}
	return chatRequestUserMessage
}

type ChatQuestionContentPartsDto interface {
	GetContentType() string
	ToChatCompletionRequestMessageContentPartClassification() (azopenai.ChatCompletionRequestMessageContentPartClassification, error)
}

type ChatQuestionContentPartsDtoType struct {
	Type string `json:"type"`
}

func (chatQuestionContentPartsDtoType *ChatQuestionContentPartsDtoType) GetContentType() string {
	return chatQuestionContentPartsDtoType.Type
}

func (chatQuestionContentPartsDtoType *ChatQuestionContentPartsDtoType) ToChatCompletionRequestMessageContentPartClassification() (azopenai.ChatCompletionRequestMessageContentPartClassification, error) {
	return nil, nil
}

type ChatQuestionContentPartsDtoText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (chatQuestionContentPartsDtoText *ChatQuestionContentPartsDtoText) GetContentType() string {
	return chatQuestionContentPartsDtoText.Type
}

func (chatQuestionContentPartsDtoText *ChatQuestionContentPartsDtoText) ToChatCompletionRequestMessageContentPartClassification() (azopenai.ChatCompletionRequestMessageContentPartClassification, error) {
	return &azopenai.ChatCompletionRequestMessageContentPartText{
		Text: &chatQuestionContentPartsDtoText.Text,
	}, nil
}

type ChatQuestionContentPartsDtoImage struct {
	Type     string `json:"type"`
	ImageUrl string `json:"imageUrl"`
}

func (chatQuestionContentPartsDtoImage *ChatQuestionContentPartsDtoImage) GetContentType() string {
	return chatQuestionContentPartsDtoImage.Type
}

func (chatQuestionContentPartsDtoImage *ChatQuestionContentPartsDtoImage) ToChatCompletionRequestMessageContentPartClassification() (azopenai.ChatCompletionRequestMessageContentPartClassification, error) {
	return &azopenai.ChatCompletionRequestMessageContentPartImage{
		ImageURL: &azopenai.ChatCompletionRequestMessageContentPartImageURL{
			URL: &chatQuestionContentPartsDtoImage.ImageUrl,
		},
	}, nil
}

type ChatQuestionContentPartsDtoAudio struct {
	Type string `json:"type"`
	Data string `json:"imageUrl"`
}

func (chatQuestionContentPartsDtoAudio *ChatQuestionContentPartsDtoAudio) GetContentType() string {
	return chatQuestionContentPartsDtoAudio.Type
}

func (chatQuestionContentPartsDtoAudio *ChatQuestionContentPartsDtoAudio) ToChatCompletionRequestMessageContentPartClassification() (azopenai.ChatCompletionRequestMessageContentPartClassification, error) {
	return &azopenai.ChatCompletionRequestMessageContentPartText{
		Text: &chatQuestionContentPartsDtoAudio.Data,
	}, nil
}

type ChatQuestionContentPartsDtoImageOCR struct {
	Type     string `json:"type"`
	ImageUrl string `json:"imageUrl"`
}

func (chatQuestionContentPartsDtoImageOCR *ChatQuestionContentPartsDtoImageOCR) GetContentType() string {
	return chatQuestionContentPartsDtoImageOCR.Type
}

func (chatQuestionContentPartsDtoImageOCR *ChatQuestionContentPartsDtoImageOCR) ToChatCompletionRequestMessageContentPartClassification() (azopenai.ChatCompletionRequestMessageContentPartClassification, error) {
	return &azopenai.ChatCompletionRequestMessageContentPartText{
		Text: &chatQuestionContentPartsDtoImageOCR.ImageUrl,
	}, nil
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
