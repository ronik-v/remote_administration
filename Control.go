package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func copy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s HOST:PORT ", os.Args[0])
		os.Exit(1)
	}
	server := os.Args[1]
	conn, _ := net.Dial("tcp", server)
	go copy(os.Stdout, conn)
	copy(conn, os.Stdin)
}
