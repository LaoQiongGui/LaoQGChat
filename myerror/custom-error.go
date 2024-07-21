package myerror

import "fmt"

type CustomError struct {
	StatusCode  int
	MessageCode string
	MessageText string
}

func (e *CustomError) Error() string {
	return fmt.Sprintf(
		"StatusCode: %d, MessageCode: %s, MessageText: %s", e.StatusCode, e.MessageCode, e.MessageText)
}
