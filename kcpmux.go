package main

import (
	"log"
	"net"
	"sync"

	"github.com/xtaci/smux"
)

type KCPMuxListener struct {
	net.Listener
	ch chan interface{}
}

func NewKCPMuxListener(laddr string, conf *KCPConfig) (net.Listener, error) {
	l, err := NewKCPListener(laddr, conf)
	if err != nil {
		return nil, err
	}

	ret := &KCPMuxListener{
		Listener: l,
		ch:       make(chan interface{}, 10),
	}
	go ret.start()

	return ret, nil
}

func (l *KCPMuxListener) start() {
	for {
		c, err := l.Listener.Accept()
		if err != nil {
			l.ch <- err
			continue
		}
		ss, err := smux.Server(c, smux.DefaultConfig())
		if err != nil {
			l.ch <- err
			continue
		}

		go func() {
			for !ss.IsClosed() {
				s, err := ss.AcceptStream()
				if err != nil {
					l.ch <- err
					continue
				}
				l.ch <- s
			}
		}()
	}
}

func (l *KCPMuxListener) Accept() (net.Conn, error) {
	switch ret := (<-l.ch).(type) {
	case net.Conn:
		return ret, nil
	case error:
		return nil, ret
	default:
		log.Fatal("KCPMuxListener invalid accept result", ret)
	}
	return nil, nil
}

var kcpMuxConns = make(map[string]*smux.Session)
var lock sync.Mutex

func DialKCPMux(raddr string, conf *KCPConfig) (net.Conn, error) {
	lock.Lock()
	s, ok := kcpMuxConns[raddr]
	if !ok || s.IsClosed() {
		conn, err := DialKCP(raddr, conf)
		if err != nil {
			return nil, err
		}
		s, err = smux.Client(conn, smux.DefaultConfig())
		if err != nil {
			return nil, err
		}
		kcpMuxConns[raddr] = s
	}
	lock.Unlock()
	return s.OpenStream()
}
