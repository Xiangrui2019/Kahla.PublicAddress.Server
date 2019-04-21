package models

import (
	"fmt"
	"net/http"
)

type ResponseStatusCodeNot200 struct {
	Response   *http.Response
	StatusCode int
}

func (r *ResponseStatusCodeNot200) Error() string {
	return fmt.Sprintf("response status code not 200: %d", r.StatusCode)
}