package shadow

import (
	"net"

	"github.com/shadowsocks/go-shadowsocks2/core"
)

type Config struct {
	Crypt    string `json:"crypt"`
	Password string `json:"password"`
}

type Tunnel struct {
	cipher core.StreamConnCipher
}

func (t *Tunnel) With(conn net.Conn) (net.Conn, error) {
	return t.cipher.StreamConn(conn), nil
}

func NewTunnel(conf *Config) (*Tunnel, error) {
	cipher, err := core.PickCipher(conf.Crypt, nil, conf.Password)
	if err != nil {
		return nil, err
	}
	return &Tunnel{
		cipher: cipher,
	}, nil
}
