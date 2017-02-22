package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"io"
	"math/rand"
	"net"
)

type TlServer struct {
	conn    net.Conn
	secKey  string
	commKey []byte
	encoder cipher.Stream
	decoder cipher.Stream
}

func NewTlServer(conn net.Conn, secKey string) *TlServer {
	aesKey := sha256.Sum256([]byte(secKey))
	iv := make([]byte, aes.BlockSize)

	encoderCipher, _ := aes.NewCipher(aesKey[:])
	encoder := cipher.NewCFBEncrypter(encoderCipher, iv)
	decoderCipher, _ := aes.NewCipher(aesKey[:])
	decoder := cipher.NewCFBDecrypter(decoderCipher, iv)

	s := new(TlServer)
	s.conn = conn
	s.secKey = secKey
	s.encoder = encoder
	s.decoder = decoder
	return s
}

func (s *TlServer) Read(buf []byte) (int, error) {
	n, err := s.conn.Read(buf)
	if err != nil {
		return n, err
	}

	s.decoder.XORKeyStream(buf[:n], buf[:n])

	return n, nil
}

func (s *TlServer) Write(buf []byte) (int, error) {
	s.encoder.XORKeyStream(buf, buf)

	n, err := s.conn.Write(buf)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (s *TlServer) Accept() (net.Conn, error) {
	b := make([]byte, 512)

	// read request addr
	_, err := io.ReadFull(s, b[:2])
	if err != nil {
		return nil, err
	}
	addrLen := int(b[0])<<8 + int(b[1])
	if addrLen > 260 {
		return nil, ErrProtocolError
	}

	_, err = io.ReadFull(s, b[:addrLen])
	if err != nil {
		return nil, err
	}
	addr := string(b[:addrLen])

	// check addr sign
	_, err = io.ReadFull(s, b[:32])
	if err != nil {
		return nil, err
	}
	sign := sha256.Sum256([]byte(addr + s.secKey))
	if bytes.Equal(sign[:], b[:32]) == false {
		return nil, ErrProtocolError
	}

	// connect to request addr
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		return nil, err
	}

	// create communication key
	commKeyLen := 32 + rand.Intn(128-32)
	s.commKey = make([]byte, commKeyLen)
	for i := 0; i < commKeyLen; i++ {
		s.commKey[i] = byte(rand.Intn(256))
	}

	// send communication key
	b[0] = byte(commKeyLen)
	copy(b[1:], s.commKey)
	s.Write(b[:commKeyLen+1])

	return conn, nil
}
