package main

import (
	"io"
	"net"
)

type TlServer struct {
	conn net.Conn
}

func NewTlServer(conn net.Conn) *TlServer {
	s := new(TlServer)
	s.conn = conn
	return s
}

func (s *TlServer) Read(buf []byte) (int, error) {
	// TODO: decode here

	n, err := s.conn.Read(buf)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (s *TlServer) Write(buf []byte) (int, error) {
	// TODO: encode here

	n, err := s.conn.Write(buf)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (s *TlServer) Accept() (net.Conn, error) {
	b := make([]byte, 512)

	_, err := io.ReadFull(s, b[:2])
	if err != nil {
		return nil, err
	}
	addrLen := int(b[0])<<8 + int(b[1])

	_, err = io.ReadFull(s, b[:addrLen])
	if err != nil {
		return nil, err
	}
	addr := string(b[:addrLen])

	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		return nil, err
	}

	s.Write([]byte{0x00})

	return conn, nil
}
