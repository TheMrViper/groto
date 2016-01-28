package groto

import (
	"encoding/gob"
	"net"
)

type gConn struct {
	net.Conn

	encoder *gob.Encoder
	decoder *gob.Decoder
}

func newGConn(conn net.Conn) *gConn {
	GConn := &gConn{
		Conn:    conn,
		encoder: gob.NewEncoder(conn),
		decoder: gob.NewDecoder(conn),
	}

	return GConn
}

func (conn *gConn) send(inPtr interface{}) error {
	return conn.encoder.Encode(inPtr)
}

func (conn *gConn) recv(outPtr interface{}) error {
	return conn.decoder.Decode(outPtr)
}
