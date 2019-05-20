package main

import (
	"net"

	"github.com/fancl20/fproxy/transport"
	"github.com/fancl20/fproxy/tunnel"
)

type listener struct {
	l   net.Listener
	tun tunnel.Tunnel
}

func (l *listener) Accept() (net.Conn, error) {
	c, err := l.l.Accept()
	if err != nil {
		return nil, err
	}
	return l.tun.With(c)
}

func (l *listener) Close() error {
	return l.l.Close()
}

func (l *listener) Addr() net.Addr {
	return l.l.Addr()
}

func parseListener(args []string) (net.Listener, error) {
	u := parse(args[1])
	l, err := transport.NewListener(u.Scheme, u.Host, u.User)
	if err != nil {
		return nil, err
	}

	tuns := &tunnel.TunnelComposer{}
	for _, v := range args[2:] {
		u := parse(v)
		if u.Scheme != "tun" {
			break
		}
		t, err := tunnel.NewTunnel(u.Host, u.User)
		if err != nil {
			return nil, err
		}
		tuns.Compose(t)
	}
	return &listener{
		l:   l,
		tun: tuns,
	}, nil
}
