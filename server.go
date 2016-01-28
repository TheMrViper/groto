package groto

import "net"

type (
	Server struct {
		// config
		config ServerConfig

		// handlers
		handlers map[string]ServerRequest
	}

	ServerConfig struct {
		// listen addres
		Addres string `json:"addres"`
	}
)

func (conf ServerConfig) NewServer() *Server {
	serv := &Server{
		config: conf,

		handlers: make(map[string]ServerRequest),
	}
	return serv
}

func (serv *Server) Listen() error {

	listener, err := net.Listen("tcp", serv.config.Addres)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go serv.listenConn(newGConn(conn))
	}
}

func (serv *Server) Handle(request ServerRequest) {
	serv.handlers[request.Name()] = request
}

func (serv *Server) listenConn(sConn *gConn) {
	for {

		// receive request headers
		requestHeader := &RequestHeader{}
		if err := sConn.recv(requestHeader); err != nil {
			sConn.Close()
			return
		}

		// check request action
		creater, ok := serv.handlers[requestHeader.Action]
		if !ok {
			// when we receive unknown action
			// we must clear socket
			// and send error

			// clear socket
			if err := sConn.recv(&struct{}{}); err != nil {
				sConn.Close()
				return
			}

			// send error
			responseHeader := &ResponseHeader{}
			responseHeader.ErrorCode = 1
			if err := sConn.send(responseHeader); err != nil {
				sConn.Close()
				return
			}

			continue
		}

		// create new ptr for data
		requestStruct := creater.Create()

		// recv request data
		if err := sConn.recv(requestStruct); err != nil {
			sConn.Close()
			return
		}

		// send response headers
		responseHeader := &ResponseHeader{}
		responseHeader.ErrorCode = 0
		if err := sConn.send(responseHeader); err != nil {
			sConn.Close()
			return
		}

		// handle action and
		// send response data
		responseStruct := requestStruct.Handler()
		if err := sConn.send(responseStruct); err != nil {
			sConn.Close()
			return
		}
	}
}

func NewServer(addr string) *Server {
	serv := &Server{
		config: ServerConfig{
			Addres: addr,
		},

		handlers: make(map[string]ServerRequest),
	}
	return serv
}
