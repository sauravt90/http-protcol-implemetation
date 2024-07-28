package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func HadleTCPConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	str := string(buffer[:n])

	arr := strings.Split(str, "\r\n")
	reqLine := arr[0]
	reqLineArr := strings.Split(reqLine, " ")

	if strings.HasPrefix(reqLineArr[1], "/files") {
		dir := os.Args[2]
		fileName := strings.TrimPrefix(reqLineArr[1], "/files/")
		if reqLineArr[0] == "POST" {
			file, _ := os.Create(dir + fileName)
			file.WriteString(arr[len(arr)-1])
			conn.Write([]byte("HTTP/1.1 201 Created\r\n\r\n"))
		}

		fmt.Print(fileName)
		data, err := os.ReadFile(dir + fileName)
		response := ""
		if err != nil {
			response = "HTTP/1.1 404 Not Found\r\n\r\n"
		} else {
			response = fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\r\n%s", len(data), data)
		}

		conn.Write([]byte(response))
	}

	if strings.HasPrefix(reqLineArr[1], "/user-agent") {
		fmt.Println(arr[2])
		h2 := strings.Split(arr[2], ": ")
		fmt.Println(h2)
		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(h2[1])) + "\r\n\r\n" + h2[1]))
	}

	if strings.HasPrefix(reqLineArr[1], "/echo/") {
		echo := strings.Split(reqLineArr[1], "/")
		for i := 1; i < len(arr); i++ {
			h2 := strings.Split(arr[i], ": ")

			fmt.Println(h2)

			if h2[0] == "Accept-Encoding" && strings.Contains(h2[1], "gzip") {
				var b bytes.Buffer
				gz := gzip.NewWriter(&b)
				if _, err := gz.Write([]byte(echo[2])); err != nil {
					log.Fatal(err)
				}
				if err := gz.Close(); err != nil {
					log.Fatal(err)
				}
				fmt.Println(b.Bytes())
				conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Encoding: gzip\r\nContent-Length: " + strconv.Itoa(b.Len()) + "\r\n\r\n" + string(b.Bytes())))
			} else if h2[0] == "Accept-Encoding" {
				conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n"))
			}
		}
		fmt.Println(echo)
		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: " + strconv.Itoa(len(echo[2])) + "\r\n\r\n" + echo[2]))
	}

	if reqLineArr[1] != "/" {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	} else {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	}

}
