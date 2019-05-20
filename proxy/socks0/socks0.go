package socks0

import (
	"net"

	"github.com/shadowsocks/go-shadowsocks2/socks"
)

type Proxy struct{}

func (*Proxy) Accept(conn net.Conn) (net.Conn, error) {
	addr, err := socks.ReadAddr(conn)
	if err != nil {
		return nil, err
	}
	return net.Dial("tcp", addr.String())
}
