package tunnel

import (
	"errors"
	"net"
	"net/url"

	"github.com/fancl20/fproxy/tunnel/shadow"
)

var (
	ErrTunnelSchemeNotSupported = errors.New("tunnel: scheme not supported")
)

type Tunnel interface {
	With(net.Conn) (net.Conn, error)
}

type TunnelComposer struct {
	tuns []Tunnel
}

func (t *TunnelComposer) With(conn net.Conn) (net.Conn, error) {
	var err error
	for _, v := range t.tuns {
		conn, err = v.With(conn)
		if err != nil {
			return nil, err
		}
	}
	return conn, err
}
func (t *TunnelComposer) Compose(tun Tunnel) {
	t.tuns = append(t.tuns, tun)
}

func NewTunnel(scheme string, u *url.Userinfo) (Tunnel, error) {
	switch scheme {
	case "shadow":
		pwd, _ := u.Password()
		conf := shadow.Config{
			Crypt:    u.Username(),
			Password: pwd,
		}
		return shadow.NewTunnel(&conf)
	}
	return nil, ErrTunnelSchemeNotSupported
}
