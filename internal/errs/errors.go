package errs

import "net/http"

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
	status  int
}

func (e *Error) Error() string { return e.Message }
func (e *Error) Status() int   { return e.status }

func NotFound(msg string) *Error {
	return &Error{
		Code:    "NOT_FOUND",
		Message: msg,
		status:  http.StatusNotFound,
	}
}

func BadRequest(msg string, details ...any) *Error {
	err := &Error{
		Code:    "BAD_REQUEST",
		Message: msg,
		status:  http.StatusBadRequest,
	}
	if len(details) > 0 {
		err.Details = details[0]
	}
	return err
}

func Internal(msg string) *Error {
	return &Error{
		Code:    "INTERNAL_ERROR",
		Message: msg,
		status:  http.StatusInternalServerError,
	}
}
