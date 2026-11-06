package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var addr string = "localhost:8888"

func main() {
	var conn net.Conn

	config := readConfig("client.cfg")
	for _, kv := range config {
		key, value := kv[0], kv[1]
		switch key {
		case "addr":
			addr = value
		}	
	}

	for {
		conn = nil
		for conn == nil {
			conn = dialServer(addr)
			if conn == nil {
				time.Sleep(5 * time.Second)
			}
		}
		fmt.Printf("Host: %v\n", conn.LocalAddr().String())
		readServer(conn)
	}
}

func dialServer(addr string) (conn net.Conn) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Dail error: %v\n", err)
		return nil
	}

	return conn
}

func readServer(conn net.Conn) {
	var buf []byte

	for {
		buf = make([]byte, 256)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Read error: %v", err)
			break
		}
		fmt.Printf("Msg: %s\n", buf)
	}
}

func readConfig(name string) (config [][]string) {
	data, err := os.ReadFile(name)
	if err != nil {
		log.Fatalf("read error: %v\n", err)
	}

	lines := strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n")

	for _, line := range lines {
		keyValuePair := strings.Split(line, "=")
		if len(keyValuePair) == 2 {
			config = append(config, keyValuePair)
		}
	}

	return config
}
