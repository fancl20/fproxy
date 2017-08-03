package kcp

type Config struct {
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
}

var DefaultConfig = Config{
	Key:          "",
	Crypt:        "aes",
	MTU:          1350,
	SndWnd:       128,
	RcvWnd:       512,
	DataShard:    10,
	ParityShard:  3,
	DSCP:         0,
	AckNodelay:   false,
	NoDelay:      0,
	Interval:     30,
	Resend:       2,
	NoCongestion: 1,
	SockBuf:      4194304,
}
