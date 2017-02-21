package main

import (
	"io"
	"log"
	"net"
)

var Config struct {
	localAddr  string
	serverAddr string
}

func proxy(clientConn net.Conn) {
	defer clientConn.Close()

	// create socks5 server
	s := NewSocks5Server(clientConn)
	err := s.MethodSelect()
	if err != nil {
		return
	}
	dstAddr, err := s.ReceiveDstAddr()
	if err != nil {
		return
	}
	log.Printf("proxy_request: [%s] => [%s]\n",
		clientConn.RemoteAddr(), dstAddr)

	// create tl client
	serverConn, err := net.Dial("tcp4", Config.serverAddr)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	c := NewTlClient(serverConn)

	err = c.Connect(dstAddr)
	if err != nil {
		return
	}

	err = s.NotifyConnectSuccess()
	if err != nil {
		return
	}

	go io.Copy(c, clientConn)
	io.Copy(clientConn, c)
}

func handleProxy() {
	l, err := net.Listen("tcp4", Config.localAddr)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("%v", err)
		}
		go proxy(conn)
	}
}

func main() {
	Config.localAddr = "0.0.0.0:8000"
	Config.serverAddr = "127.0.0.1:8001"

	handleProxy()
}
