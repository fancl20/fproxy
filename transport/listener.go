package transport

import (
	"errors"
	"net"
	"net/url"

	"github.com/fancl20/fproxy/transport/kcp"
	"github.com/fancl20/fproxy/transport/tcp"
)

var (
	ErrListenSchemeNotSupported = errors.New("transport: listen scheme not supported")
)

func NewListener(scheme, host string, u *url.Userinfo) (net.Listener, error) {
	switch scheme {
	case "tcp":
		return tcp.Listen(host)
	case "kcp":
		conf := kcp.DefaultConfig
		conf.Crypt = u.Username()
		conf.Key, _ = u.Password()
		l := func(addr string) (net.Listener, error) {
			return kcp.Listen(addr, &conf)
		}
		return WithMuxListen(l)(host)
	}
	return nil, ErrListenSchemeNotSupported
}
