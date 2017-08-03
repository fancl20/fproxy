package kcp

import (
	"errors"
	"net"

	"github.com/xtaci/kcp-go"
)

var (
	ErrInvalidKey = errors.New("kcp: crypt key not set")
)

func Dial(addr string, conf *Config) (net.Conn, error) {
	if conf.Key == "" {
		return nil, ErrInvalidKey
	}
	blk := newCrypt(conf.Crypt, conf.Key)
	conn, err := kcp.DialWithOptions(addr, blk, conf.DataShard, conf.ParityShard)
	if err != nil {
		return nil, err
	}
	return conn, setKCPSession(conn, conf)
}

type listener struct {
	*kcp.Listener
	conf *Config
}

func Listen(addr string, conf *Config) (net.Listener, error) {
	if conf.Key == "" {
		return nil, ErrInvalidKey
	}
	blk := newCrypt(conf.Crypt, conf.Key)
	l, err := kcp.ListenWithOptions(addr, blk, conf.DataShard, conf.ParityShard)
	if err != nil {
		return nil, err
	}
	return &listener{
		Listener: l,
		conf:     conf,
	}, nil
}

func (l *listener) Accept() (net.Conn, error) {
	conn, err := l.Listener.AcceptKCP()
	if err != nil {
		return nil, err
	}
	return conn, setKCPSession(conn, l.conf)
}

func setKCPSession(conn *kcp.UDPSession, conf *Config) error {
	conn.SetStreamMode(true)
	conn.SetWriteDelay(true)
	conn.SetMtu(conf.MTU)
	conn.SetNoDelay(conf.NoDelay, conf.Interval, conf.Resend, conf.NoCongestion)
	conn.SetWindowSize(conf.SndWnd, conf.RcvWnd)
	conn.SetACKNoDelay(conf.AckNodelay)

	// ignore returned error
	conn.SetDSCP(conf.DSCP)
	conn.SetReadBuffer(conf.SockBuf)
	conn.SetWriteBuffer(conf.SockBuf)
	return nil
}
