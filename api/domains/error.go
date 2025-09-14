package domains

import "fmt"

type HttpError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e HttpError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func NewHttpError(code int, message string) HttpError {
	return HttpError{
		Code:    code,
		Message: message,
	}
}
