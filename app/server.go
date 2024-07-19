package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	con, err := l.Accept()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	req := make([]byte, 1024)
	n, err := con.Read(req)
	if err != nil {
		return
	}
	s := string(req[:n])
	path := strings.Split(s, " ")[1]
	if path == "/" {
		con.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		con.Close()
		return
	}

	if strings.Split(path, "/")[1] == "echo" {
		message := strings.Split(path, "/")[2]
		res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length:%v\r\n\r\n%s", len(message), message)
		con.Write([]byte(res))
		return
	}

	con.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	con.Close()
}
