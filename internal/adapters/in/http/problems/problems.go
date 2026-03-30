package problems

import (
	"delivery/internal/pkg/errs"
	"errors"
	"net/http"
)

type Problem struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewNotFound(message string) Problem {
	return Problem{Status: http.StatusNotFound, Message: message}
}

func NewConflict(message string) Problem {
	return Problem{Status: http.StatusConflict, Message: message}
}

func NewBadRequest(message string) Problem {
	return Problem{Status: http.StatusBadRequest, Message: message}
}

func NewInternal(message string) Problem {
	return Problem{Status: http.StatusInternalServerError, Message: message}
}

func NewUnprocessable(message string) Problem {
	return Problem{Status: http.StatusUnprocessableEntity, Message: message}
}

func FromError(err error) Problem {
	var notFound *errs.NotFoundError
	if errors.As(err, &notFound) {
		return NewNotFound(err.Error())
	}

	var valueIsRequired *errs.ValueIsRequiredError
	if errors.As(err, &valueIsRequired) {
		return NewBadRequest(err.Error())
	}

	var valueIsInvalid *errs.ValueIsInvalidError
	if errors.As(err, &valueIsInvalid) {
		return NewUnprocessable(err.Error())
	}

	return NewInternal(err.Error())
}
