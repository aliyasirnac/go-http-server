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
	for {
		con, err := l.Accept()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(con)
	}
}

func handleConnection(con net.Conn) {
	defer con.Close()
	for {

		req := make([]byte, 1024)
		n, err := con.Read(req)
		if err != nil {
			return
		}
		s := string(req[:n])
		splitted := strings.Split(s, "\r\n")
		path := strings.Split(s, " ")[1]
		requestLine := strings.Fields(s)
		method := requestLine[0]
		fmt.Println("requestLine", splitted[len(splitted)-1])
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
		if p == "files" {
			dir := os.Args[2]
			fileName := strings.TrimPrefix(path, "/files/")

			if method == "POST" {
				content := strings.Trim(splitted[len(splitted)-1], "\x00")
				err := os.WriteFile(dir+fileName, []byte(content), 0644)

				if err != nil {
					con.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\n"))
					return
				}
				con.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
				return
			}
			fmt.Println("buraya girdi")
			data, err := os.ReadFile(dir + fileName)

			if err != nil {
				con.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
				return
			}

			response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(data), data)

			con.Write([]byte(response))
			return
		}
		con.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))

	}
}
