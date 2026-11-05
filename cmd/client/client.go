package main

import (
	"fmt"
	"net"
	"time"
)

const addr string = "localhost:8888"

func main() {
	var conn net.Conn

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
