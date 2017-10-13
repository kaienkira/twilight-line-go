package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"net"
)

type TlClient struct {
	conn              net.Conn
	secKey            string
	fakeRequestBytes  []byte
	fakeResponseBytes []byte
	commKey           []byte
	encoder           cipher.Stream
	decoder           cipher.Stream
}

func NewTlClient(
	conn net.Conn,
	secKey string,
	fakeRequestBytes []byte,
	fakeResponseBytes []byte) *TlClient {

	c := new(TlClient)
	c.conn = conn
	c.secKey = secKey
	c.fakeRequestBytes = fakeRequestBytes
	c.fakeResponseBytes = fakeResponseBytes
	c.ResetCipher([]byte(secKey))
	return c
}

func (c *TlClient) ResetCipher(key []byte) {
	aesKey := sha256.Sum256(key)
	iv := make([]byte, aes.BlockSize)

	encoderCipher, _ := aes.NewCipher(aesKey[:])
	encoder := cipher.NewCFBEncrypter(encoderCipher, iv)
	decoderCipher, _ := aes.NewCipher(aesKey[:])
	decoder := cipher.NewCFBDecrypter(decoderCipher, iv)

	c.encoder = encoder
	c.decoder = decoder
}

func (c *TlClient) Read(buf []byte) (int, error) {
	n, err := c.conn.Read(buf)
	if err != nil {
		return n, err
	}

	c.decoder.XORKeyStream(buf[:n], buf[:n])

	return n, nil
}

func (c *TlClient) Write(buf []byte) (int, error) {
	c.encoder.XORKeyStream(buf, buf)

	n, err := c.conn.Write(buf)
	if err != nil {
		return n, err
	}

	return n, nil
}

func (c *TlClient) Connect(dstAddr string) error {
	{
		b2 := new(bytes.Buffer)
		// write fake request
		b2.Write(c.fakeRequestBytes[:])
		{
			b3 := new(bytes.Buffer)
			// write request addr
			binary.Write(b3, binary.BigEndian, int16(len(dstAddr)))
			b3.WriteString(dstAddr)
			sign := sha256.Sum256([]byte(dstAddr + c.secKey))
			b3.Write(sign[:])
			// encrypt data
			c.encoder.XORKeyStream(b3.Bytes(), b3.Bytes())
			b3.WriteTo(b2)
		}
		b2.WriteTo(c.conn)
	}
	{
		// read fake response
		b2 := make([]byte, len(c.fakeResponseBytes))
		_, err := io.ReadFull(c.conn, b2[:])
		if err != nil {
			return err
		}
	}

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

	// reset cipher 
	c.ResetCipher(c.commKey)

	return nil
}
