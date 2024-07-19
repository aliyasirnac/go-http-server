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
	fmt.Println(path)
	if path == "/" {
		con.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		con.Close()
		return
	}
	p := strings.Split(path, "/")[1]
	if p == "echo" {
		message := strings.Split(path, "/")[2]
		res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length:%v\r\n\r\n%s", len(message), message)
		con.Write([]byte(res))
		return
	}
	if p == "user-agent" {
		userAgent := strings.Split(s, "\r\n")[2]
		userAgent = strings.Split(userAgent, ": ")[1]
		fmt.Println("user-agent: ", userAgent)
		res := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length:%v\r\n\r\n%s", len(userAgent), userAgent)
		con.Write([]byte(res))
		return
	}
	con.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	con.Close()
}
