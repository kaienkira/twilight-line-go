package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"net"
)

type TlClient struct {
	conn    net.Conn
	secKey  string
	commKey []byte
}

func NewTlClient(conn net.Conn, secKey string) *TlClient {
	c := new(TlClient)
	c.conn = conn
	c.secKey = secKey
	return c
}

func (c *TlClient) Read(buf []byte) (int, error) {
	// TODO: decode here

	n, err := c.conn.Read(buf)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (c *TlClient) Write(buf []byte) (int, error) {
	// TODO: encode here

	n, err := c.conn.Write(buf)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (c *TlClient) Connect(dstAddr string) error {
	{
		// send request addr
		b := new(bytes.Buffer)
		binary.Write(b, binary.BigEndian, int16(len(dstAddr)))
		b.WriteString(dstAddr)
		sign := sha256.Sum256([]byte(dstAddr + c.secKey))
		b.Write(sign[:])
		b.WriteTo(c)
	}
	{
		// create communication key
		b := make([]byte, 256)
		_, err := io.ReadFull(c, b[:1])
		if err != nil {
			return err
		}
		commKeyLen := int(b[0])
		c.commKey = make([]byte, commKeyLen)
		_, err = io.ReadFull(c, c.commKey)
		if err != nil {
			return err
		}
	}

	return nil
}
