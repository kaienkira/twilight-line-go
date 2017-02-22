package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
)

type TlClient struct {
	conn net.Conn
}

func NewTlClient(conn net.Conn) *TlClient {
	c := new(TlClient)
	c.conn = conn
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
		b := new(bytes.Buffer)
		binary.Write(b, binary.BigEndian, int16(len(dstAddr)))
		b.WriteString(dstAddr)
		b.WriteTo(c)
	}
	{
		b := make([]byte, 1)
		_, err := io.ReadFull(c, b[:1])
		if err != nil {
			return nil
		}
	}

	return nil
}