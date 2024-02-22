package conn

import "net/http"

type ConnectionHandler func(conn Conn)

type Server interface {
	SetConnHandler(handler ConnectionHandler)
	Run(host string, port int, endpoints []Router) error
}

type Router struct {
	Path    string
	Method  string
	Handler func(w http.ResponseWriter, r *http.Request)
}
