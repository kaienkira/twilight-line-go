package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type Configuration struct {
	LocalAddr         string
	ServerAddr        string
	SecKey            string
	FakeRequest       []string
	FakeResponse      []string
	fakeRequestBytes  []byte
	fakeResponseBytes []byte
}

var Config Configuration

func copyData(src io.Reader, dest io.Writer, quitSignal chan bool) {
	b := make([]byte, 32*1024)
	for {
		n, err := src.Read(b)
		if err != nil {
			quitSignal <- true
			return
		}
		dest.Write(b[:n])
	}
}

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

	c := NewTlClient(serverConn, Config.SecKey,
		Config.fakeRequestBytes, Config.fakeResponseBytes)

	err = c.Connect(dstAddr)
	if err != nil {
		return
	}

	err = s.NotifyConnectSuccess()
	if err != nil {
		return
	}

	quitSignal := make(chan bool, 2)
	go copyData(c, clientConn, quitSignal)
	go copyData(clientConn, c, quitSignal)
	<-quitSignal
	clientConn.Close()
	serverConn.Close()
	<-quitSignal
	close(quitSignal)
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
	b, err := os.ReadFile(configFile)
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
	Config.fakeRequestBytes =
		[]byte(strings.Join(Config.FakeRequest, ""))
	Config.fakeResponseBytes =
		[]byte(strings.Join(Config.FakeResponse, ""))

	if checkConfig() == false {
		fmt.Fprintf(os.Stderr, "please read help with -h or --help\n")
		return
	}

	handleProxy()
}
