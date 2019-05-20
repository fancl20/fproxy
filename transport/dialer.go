package transport

import (
	"net"
	"net/url"

	"github.com/fancl20/fproxy/transport/kcp"
	"github.com/fancl20/fproxy/transport/tcp"
)

func NewDialFunc(scheme string, u *url.Userinfo) DialFunc {
	switch scheme {
	case "tcp":
		return tcp.Dial
	case "kcp":
		conf := kcp.DefaultConfig
		conf.Crypt = u.Username()
		conf.Key, _ = u.Password()
		l := func(addr string) (net.Conn, error) {
			return kcp.Dial(addr, &conf)
		}
		return WithMuxDial(l)
	}
	return nil
}
