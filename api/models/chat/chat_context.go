package chat

type Context struct {
	Type     string           `json:"type"`
	Contents []ContentWrapper `json:"contents"`
}

const (
	ContextTypeQuestion = "Question"
	ContextTypeAnswer   = "Answer"
	ContextTypeOption   = "Option"
)
