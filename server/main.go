package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type Configuration struct {
	LocalAddr string
	SecKey    string
}

var Config Configuration

func proxy(clientConn net.Conn) {
	defer clientConn.Close()

	// create tl server
	s := NewTlServer(clientConn, Config.SecKey)
	serverConn, err := s.Accept()
	if err != nil {
		return
	}
	defer serverConn.Close()

	go io.Copy(serverConn, s)
	io.Copy(s, serverConn)
}

func handleProxy() {
	l, err := net.Listen("tcp4", Config.LocalAddr)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("%v", err)
			continue
		}
		go proxy(conn)
	}
}

func loadConfig(configFile string) bool {
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return false
	}

	err = json.Unmarshal(b, &Config)
	if err != nil {
		return false
	}

	return true
}

func checkConfig() bool {
	if Config.LocalAddr == "" {
		fmt.Fprintf(os.Stderr, "config.localAddr is required\n")
		return false
	}

	return true
}

func main() {
	configFile := flag.String("e", "",
		"config file path")
	flag.StringVar(&Config.LocalAddr, "l", "",
		"REQUIRED: localAddr(local listen addr)")
	flag.StringVar(&Config.SecKey, "k", "",
		"secKey(secure key)")
	flag.Parse()

	if *configFile != "" {
		if loadConfig(*configFile) == false {
			fmt.Fprintf(os.Stderr,
				"load config file %s failed\n", *configFile)
			return
		}
	}
	if checkConfig() == false {
		fmt.Fprintf(os.Stderr, "please read help with -h or --help\n")
		return
	}

	handleProxy()
}
