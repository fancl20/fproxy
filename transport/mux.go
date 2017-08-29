package transport

import (
	"net"
	"sync"

	"github.com/golang/groupcache/singleflight"
	"github.com/xtaci/smux"
)

// A plugable mux for transport layer

type DialFunc func(addr string) (net.Conn, error)
type ListenFunc func(addr string) (net.Listener, error)

func WithMuxDial(f DialFunc) DialFunc {
	conns := new(sync.Map)
	var g singleflight.Group
	return func(addr string) (c net.Conn, err error) {
		ret, ok := conns.Load(addr)
		if !ok || ret.(*smux.Session).IsClosed() {
			if ret, err = g.Do("addr", genDialFunc(addr, f, conns)); err != nil {
				return nil, err
			}
		}
		return ret.(*smux.Session).OpenStream()
	}
}

func genDialFunc(addr string, f DialFunc, conns *sync.Map) func() (interface{}, error) {
	return func() (interface{}, error) {
		c, err := f(addr)
		if err != nil {
			return nil, err
		}
		s, err := smux.Client(c, smux.DefaultConfig())
		if err != nil {
			return nil, err
		}
		conns.Store(addr, s)
		return s, nil
	}
}

type muxListener struct {
	net.Listener
	conns        *sync.Map
	connChannel  chan net.Conn
	errorChannel chan error
}

func WithMuxListen(f ListenFunc) ListenFunc {
	return func(addr string) (net.Listener, error) {
		l, err := f(addr)
		if err != nil {
			return nil, err
		}
		ret := &muxListener{
			Listener:     l,
			conns:        new(sync.Map),
			connChannel:  make(chan net.Conn, 10),
			errorChannel: make(chan error, 10),
		}
		go ret.serve()
		return ret, nil
	}
}

func (l *muxListener) Accept() (net.Conn, error) {
	select {
	case c := <-l.connChannel:
		return c, nil
	case err := <-l.errorChannel:
		return nil, err
	}
}

func (l *muxListener) serve() {
	for {
		ss, err := l.acceptSession()
		if err != nil {
			l.errorChannel <- err
			continue
		}
		go l.serveSession(ss)
	}
}

func (l *muxListener) serveSession(ss *smux.Session) {
	for !ss.IsClosed() {
		s, err := ss.AcceptStream()
		if err != nil {
			l.errorChannel <- err
			continue
		}
		l.connChannel <- s
	}
}

func (l *muxListener) acceptSession() (*smux.Session, error) {
	c, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	ss, err := smux.Server(c, smux.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return ss, nil
}
