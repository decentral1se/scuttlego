package rpc

import (
	"github.com/planetary-social/go-ssb/network/rpc/transport"
	"io"
)

type ResponseWriter struct {
	req  *Request
	conn *Connection
}

func NewResponseWriter(req *Request, conn *Connection) ResponseWriter {
	return ResponseWriter{
		req:  req,
		conn: conn,
	}
}

func (rw ResponseWriter) OpenResponseStream(bodyType transport.MessageBodyType) io.WriteCloser {
	panic("not implemented")
}