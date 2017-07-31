package main

import (
	"net"

	"github.com/xtaci/kcp-go"
)

type KCPConfig struct {
	Key          string `json:"key"`
	Crypt        string `json:"crypt"`
	MTU          int    `json:"mtu"`
	SndWnd       int    `json:"sndwnd"`
	RcvWnd       int    `json:"rcvwnd"`
	DataShard    int    `json:"datashard"`
	ParityShard  int    `json:"parityshard"`
	DSCP         int    `json:"dscp"`
	AckNodelay   bool   `json:"acknodelay"`
	NoDelay      int    `json:"nodelay"`
	Interval     int    `json:"interval"`
	Resend       int    `json:"resend"`
	NoCongestion int    `json:"nc"`
	SockBuf      int    `json:"sockbuf"`
	KeepAlive    int    `json:"keepalive"`
}

var (
	DefaultKCPConf = &KCPConfig{
		Key:          "flyinhigh",
		Crypt:        "aes",
		MTU:          1350,
		SndWnd:       1024,
		RcvWnd:       1024,
		DataShard:    10,
		ParityShard:  3,
		DSCP:         0,
		AckNodelay:   false,
		NoDelay:      1,
		Interval:     20,
		Resend:       2,
		NoCongestion: 1,
		SockBuf:      4194304,
		KeepAlive:    10,
	}
)

type KCPListener struct {
	*kcp.Listener
	conf *KCPConfig
}

func SetKCPSession(conn *kcp.UDPSession, conf *KCPConfig) {
	conn.SetStreamMode(true)
	conn.SetNoDelay(conf.NoDelay, conf.Interval, conf.Resend, conf.NoCongestion)
	conn.SetMtu(conf.MTU)
	conn.SetWindowSize(conf.SndWnd, conf.RcvWnd)
	conn.SetACKNoDelay(conf.AckNodelay)
}

func NewKCPListener(laddr string, conf *KCPConfig) (net.Listener, error) {
	blk, _ := kcp.NewAESBlockCrypt([]byte(conf.Key))
	l, err := kcp.ListenWithOptions(laddr, blk, conf.DataShard, conf.ParityShard)
	if err != nil {
		return nil, err
	}
	return &KCPListener{
		Listener: l,
		conf:     conf,
	}, nil
}

func (l *KCPListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.AcceptKCP()
	if err != nil {
		return nil, err
	}
	SetKCPSession(conn, l.conf)
	return conn, nil
}

func DialKCP(raddr string, conf *KCPConfig) (net.Conn, error) {
	blk, _ := kcp.NewAESBlockCrypt([]byte(conf.Key))
	conn, err := kcp.DialWithOptions(raddr, blk, conf.DataShard, conf.ParityShard)
	if err != nil {
		return nil, err
	}
	SetKCPSession(conn, conf)
	return conn, nil
}
