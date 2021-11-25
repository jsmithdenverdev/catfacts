package http

import (
	"net/http"
)

type Server struct {
	Handler http.Handler
}

func (server Server) Listen(addr string) error {
	return http.ListenAndServe(addr, server.Handler)
}
