package main

import (
	"fmt"
	"io"
	"net"
)

type Socks5Server struct {
	conn net.Conn
}

func NewSocks5Server(conn net.Conn) *Socks5Server {
	s := new(Socks5Server)
	s.conn = conn
	return s
}

func (s *Socks5Server) MethodSelect() error {
	b := make([]byte, 256)

	_, err := io.ReadFull(s.conn, b[:2])
	if err != nil {
		return err
	}
	version := b[0]
	nMethods := b[1]

	// check version
	if version != 0x05 {
		return ErrProtocolError
	}

	// discard methods
	_, err = io.ReadFull(s.conn, b[:nMethods])
	if err != nil {
		return err
	}

	// answer server accepted method
	_, err = s.conn.Write([]byte{0x05, 0x00})
	if err != nil {
		return err
	}

	return nil
}

func (s *Socks5Server) ReceiveDstAddr() (string, error) {
	b := make([]byte, 256)

	_, err := io.ReadFull(s.conn, b[:5])
	if err != nil {
		return "", err
	}
	version := b[0]
	cmd := b[1]
	addrType := b[3]
	domainLength := b[4]

	// check version
	if version != 0x05 {
		return "", ErrProtocolError
	}
	// only support connect command
	if cmd != 0x01 {
		return "", ErrProtocolError
	}
	// only support domain
	if addrType != 0x03 {
		return "", ErrProtocolError
	}

	_, err = io.ReadFull(s.conn, b[:domainLength+2])
	if err != nil {
		return "", err
	}
	domain := string(b[:domainLength])
	port := int(b[domainLength])<<8 + int(b[domainLength+1])

	return fmt.Sprintf("%s:%d", domain, port), nil
}

func (s *Socks5Server) NotifyConnectSuccess() error {
	_, err := s.conn.Write([]byte{
		0x05, 0x00, 0x00, 0x01, 0x00,
	    0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		return err
	}

	return nil
}
