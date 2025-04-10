package util

import (
	"fmt"
	"net/http"
)

type WrappedError struct {
	Err     error
	Message string
}

func (w WrappedError) Error() string {
	return fmt.Sprintf("%s: %s", w.Err.Error(), w.Message)
}

func WriteError(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(http.StatusText(statusCode)))
}
