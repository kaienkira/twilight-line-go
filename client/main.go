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
	LocalAddr  string
	ServerAddr string
}

var Config Configuration

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
	serverConn, err := net.Dial("tcp4", Config.ServerAddr)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer serverConn.Close()

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
	l, err := net.Listen("tcp4", Config.LocalAddr)
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
	if Config.ServerAddr == "" {
		fmt.Fprintf(os.Stderr, "config.serverAddr is required\n")
		return false
	}

	return true
}

func main() {
	configFile := flag.String("e", "",
		"config file path")
	flag.StringVar(&Config.LocalAddr, "l", "",
		"REQUIRED: localAddr(local listen addr)")
	flag.StringVar(&Config.ServerAddr, "s", "",
		"REQUIRED: serverAddr(proxy server addr)")
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
