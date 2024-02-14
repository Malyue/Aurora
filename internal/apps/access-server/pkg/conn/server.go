package conn

type ConnectionHandler func(conn Conn)

type Server interface {
	SetConnHandler(handler ConnectionHandler)
	Run(host string, port int) error
}
