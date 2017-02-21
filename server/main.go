package main

import (
	"io"
	"log"
	"net"
)

var Config struct {
	localAddr string
}

func proxy(clientConn net.Conn) {
	defer clientConn.Close()

	// create tl server
	s := NewTlServer(clientConn)
	serverConn, err := s.Accept()
	if err != nil {
		return
	}

	go io.Copy(serverConn, s)
	io.Copy(s, serverConn)
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
	Config.localAddr = "0.0.0.0:8001"

	handleProxy()
}
