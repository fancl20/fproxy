package main

import (
	"io"
	"log"
	"net"
	"net/url"
	"os"
)

func Listen(typ, laddr string) (net.Listener, error) {
	switch typ {
	case "tcp":
		return net.Listen("tcp", laddr)
	case "kcp":
		return NewKCPListener(laddr, DefaultKCPConf)
	case "kcpmux":
		return NewKCPMuxListener(laddr, DefaultKCPConf)
	default:
		log.Fatalln("unsupported protocol:", typ)
	}
	return nil, nil
}

func Dial(typ, raddr string) (net.Conn, error) {
	switch typ {
	case "tcp":
		return net.Dial("tcp", raddr)
	case "kcp":
		return DialKCP(raddr, DefaultKCPConf)
	case "kcpmux":
		return DialKCPMux(raddr, DefaultKCPConf)
	default:
		log.Fatalln("unsupported protocol:", typ)
	}
	return nil, nil
}

func Forward(c1, c2 net.Conn) {
	defer c1.Close()
	defer c2.Close()

	ch1 := make(chan struct{})
	ch2 := make(chan struct{})
	go func() { io.Copy(c1, c2); close(ch1) }()
	go func() { io.Copy(c2, c1); close(ch2) }()

	select {
	case <-ch1:
	case <-ch2:
	}

	log.Println("connection closed:", c1.RemoteAddr(), c2.RemoteAddr())
}

func Parse(u string) (typ, addr string) {
	uu, err := url.Parse(u)
	if err != nil {
		log.Fatalln("invalid addr:", u)
	}
	return uu.Scheme, uu.Host
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("usage: kcproxy tcp://localAddr kcp://remoteAddr")
	}
	ltyp, laddr := Parse(os.Args[1])
	rtyp, raddr := Parse(os.Args[2])
	l, err := Listen(ltyp, laddr)
	if err != nil {
		log.Fatalln("listen error:", err)
	}
	for {
		c1, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		log.Println("connection accepted from:", c1.RemoteAddr())

		c2, err := Dial(rtyp, raddr)
		if err != nil {
			c1.Close()
			log.Println("dial error:", err)
			return
		}
		log.Println("connection dialed to:", c2.RemoteAddr())

		go Forward(c1, c2)
	}
}
