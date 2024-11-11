package chat

import (
	"encoding/json"
	"errors"
)

type ContentWrapper struct {
	Content Content
}

func (contentWrapper *ContentWrapper) MarshalJSON() ([]byte, error) {
	return json.Marshal(contentWrapper.Content)
}

func (contentWrapper *ContentWrapper) UnmarshalJSON(data []byte) error {
	// 获取类型
	typeContent := TypeContent{}
	if err := json.Unmarshal(data, &typeContent); err != nil {
		return err
	}

	// 根据类型选择反序列化实例
	switch typeContent.Type {
	case ContentTypeText:
		textContent := TextContent{}
		if err := json.Unmarshal(data, &textContent); err != nil {
			return err
		}
		contentWrapper.Content = &textContent
	case ContentTypeImage:
		imageContent := ImageContent{}
		if err := json.Unmarshal(data, &imageContent); err != nil {
			return err
		}
		contentWrapper.Content = &imageContent
	default:
		return errors.New("unknown content type: " + typeContent.Type)
	}
	return nil
}

type Content interface {
	GetContentType() string
}

const (
	ContentTypeText  = "Text"
	ContentTypeImage = "Image"
)

type TypeContent struct {
	Type string `json:"type"`
}

func (requestContent *TypeContent) GetContentType() string {
	return requestContent.Type
}

type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (requestContent *TextContent) GetContentType() string {
	return requestContent.Type
}

type ImageContent struct {
	Type     string `json:"type"`
	ImageUrl string `json:"imageUrl"`
}

func (requestContent *ImageContent) GetContentType() string {
	return requestContent.Type
}
