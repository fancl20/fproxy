package main

import (
	"io"
	"log"
	"net"
	"net/url"
	"os"

	"github.com/fancl20/fproxy/transport"
)

func forward(c1, c2 net.Conn) {
	defer c1.Close()
	defer c2.Close()

	ch := make(chan struct{}, 2)
	go func() { io.Copy(c1, c2); ch <- struct{}{} }()
	go func() { io.Copy(c2, c1); ch <- struct{}{} }()
	<-ch

	log.Println("connection closed:", c1.RemoteAddr(), c2.RemoteAddr())
}

func parse(u string) *url.URL {
	uu, err := url.Parse(u)
	if err != nil {
		log.Fatalln("invalid addr:", u)
	}
	return uu
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("usage: fproxy tcp://localAddr kcp://aes:12345@remoteAddr")
	}
	local, remote := parse(os.Args[1]), parse(os.Args[2])
	l, err := transport.NewListener(local.Scheme, local.Host, local.User)
	if err != nil {
		log.Fatalln("listen error:", err)
	}
	dial := transport.NewDialFunc(remote.Scheme, remote.User)
	if dial == nil {
		log.Fatalln("invalid dialFunc")
	}
	for {
		c1, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		log.Println("connection accepted from:", c1.RemoteAddr())

		go func() {
			c2, err := dial(remote.Host)
			if err != nil {
				c1.Close()
				log.Println("dial error:", err)
				return
			}
			log.Println("connection dialed to:", c2.RemoteAddr())

			forward(c1, c2)
		}()
	}
}
