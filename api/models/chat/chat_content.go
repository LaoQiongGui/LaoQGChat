package chat

type Content struct {
	Type  string        `json:"type"`
	Parts []PartWrapper `json:"parts"`
}

const (
	ContentTypeQuestion = "Question"
	ContentTypeAnswer   = "Answer"
	ContentTypeOption   = "Option"
)
