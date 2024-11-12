package chat

import (
	"encoding/json"
	"errors"
)

type PartWrapper struct {
	Part Part
}

func (cartWrapper *PartWrapper) MarshalJSON() ([]byte, error) {
	return json.Marshal(cartWrapper.Part)
}

func (cartWrapper *PartWrapper) UnmarshalJSON(data []byte) error {
	// 获取类型
	typeContent := TypePart{}
	if err := json.Unmarshal(data, &typeContent); err != nil {
		return err
	}

	// 根据类型选择反序列化实例
	switch typeContent.Type {
	case PartTypeText:
		textContent := TextPart{}
		if err := json.Unmarshal(data, &textContent); err != nil {
			return err
		}
		cartWrapper.Part = &textContent
	case PartTypeImage:
		imageContent := ImagePart{}
		if err := json.Unmarshal(data, &imageContent); err != nil {
			return err
		}
		cartWrapper.Part = &imageContent
	default:
		return errors.New("unknown content type: " + typeContent.Type)
	}
	return nil
}

type Part interface {
	GetContentType() string
}

const (
	PartTypeText  = "Text"
	PartTypeImage = "Image"
)

type TypePart struct {
	Type string `json:"type"`
}

func (requestPart *TypePart) GetContentType() string {
	return requestPart.Type
}

type TextPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (requestPart *TextPart) GetContentType() string {
	return requestPart.Type
}

type ImagePart struct {
	Type     string `json:"type"`
	ImageUrl string `json:"imageUrl"`
}

func (requestPart *ImagePart) GetContentType() string {
	return requestPart.Type
}
