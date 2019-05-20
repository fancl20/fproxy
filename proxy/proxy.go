package proxy

import (
	"errors"
	"net"
	"net/url"

	"github.com/fancl20/fproxy/proxy/socks0"
	"github.com/fancl20/fproxy/proxy/socks5"
)

var (
	ErrProxySchemeNotSupported = errors.New("proxy: scheme not supported")
)

type Proxy interface {
	Accept(net.Conn) (net.Conn, error)
}

func NewProxy(scheme string, u *url.Userinfo) (Proxy, error) {
	switch scheme {
	case "socks0":
		return &socks0.Proxy{}, nil
	case "socks5":
		return &socks5.Proxy{}, nil
	}
	return nil, ErrProxySchemeNotSupported
}
