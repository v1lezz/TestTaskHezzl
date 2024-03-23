package http

import (
	"fmt"
	"net/http"
)

func NewHTTPServer(port int) *http.Server {
	return &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}
}
