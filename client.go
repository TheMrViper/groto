package groto

import (
	"errors"
	"net"
	"sync"
)

type (
	Client struct {
		// config from file
		config ClientConfig

		// idle conns mutex
		mutex sync.Mutex
		// idle conns pool
		idleConns []*gConn
	}

	ClientConfig struct {
		// connection addres
		Addres string `json:"addres"`
		// max idle conns
		MaxIdleConns int `json:"max_idle_conns"`
	}
)

var (
	UnknownAction = errors.New("groto-server: received unknown action from client")
)

func (conf ClientConfig) NewClient() *Client {
	cli := &Client{
		config: conf,

		idleConns: make([]*gConn, 0, conf.MaxIdleConns),
	}
	return cli
}

func (cli *Client) newConn() (*gConn, error) {
	conn, err := net.Dial("tcp", cli.config.Addres)
	if err != nil {
		return nil, err
	}
	return newGConn(conn), nil
}

func (cli *Client) getConn() (*gConn, error) {
	cli.mutex.Lock()
	defer cli.mutex.Unlock()

	if len(cli.idleConns) > 0 {
		cli.idleConns = cli.idleConns[1:]
		return cli.idleConns[0], nil
	}

	return cli.newConn()
}

func (cli *Client) putConn(conn *gConn) {
	cli.mutex.Lock()
	defer cli.mutex.Unlock()

	if len(cli.idleConns) < cli.config.MaxIdleConns || cli.config.MaxIdleConns == 0 {
		cli.idleConns = append(cli.idleConns, conn)
		return
	}
	conn.Close()
	return
}

func (cli *Client) Send(request ClientRequest, response ClientResponse) (err error) {

	// get conn from pool or new
	cConn, err := cli.getConn()
	if err != nil {
		return err
	}

	// send request headers
	requestHeader := &RequestHeader{}
	requestHeader.Action = request.Name()
	if err := cConn.send(requestHeader); err != nil {
		cConn.Close()
		return err
	}

	// send request data
	if err := cConn.send(request); err != nil {
		cConn.Close()
		return err
	}

	// recv response headers
	responseHeader := &ResponseHeader{}
	if err := cConn.recv(responseHeader); err != nil {
		cConn.Close()
		return err
	}

	// check response error
	if responseHeader.ErrorCode != 0 {
		cli.putConn(cConn)
		return UnknownAction
	}

	// recv response data
	if err := cConn.recv(response); err != nil {
		cConn.Close()
		return err
	}

	// put conn in pool
	cli.putConn(cConn)
	return nil
}

func NewClient(addr string, maxIdleConns int) *Client {
	cli := &Client{
		config: ClientConfig{
			Addres:       addr,
			MaxIdleConns: maxIdleConns,
		},

		idleConns: make([]*gConn, 0, maxIdleConns),
	}
	return cli
}
