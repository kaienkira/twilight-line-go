package main

import (
	"log"
	"net"
)

func proxy(conn net.Conn) {
	defer conn.Close()

	log.Printf("new connection from: %s\n", conn.RemoteAddr())
}

func runServer(addr string) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go proxy(conn)
	}
}

func main() {
	runServer("0.0.0.0:8001")
}
