package tcp

import (
	"net"
)

func Dial(addr string) (net.Conn, error) {
	return net.Dial("tcp", addr)
}

func Listen(addr string) (net.Listener, error) {
	return net.Listen("tcp", addr)
}
