package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"
)

const (
	signature      = "GFS_V1"
	protocolPort   = ":9999"
	maxConnections = 5
)

func runSender(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Invalid path:", err)
		return
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		fmt.Println("File error:", err)
		return
	}

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Println("Listen error:", err)
		return
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	fmt.Printf("[Host] Sharing %s on port %d\n", fileInfo.Name(), port)

	go broadcast(port, fileInfo.Name(), fileInfo.Size())

	sem := make(chan struct{}, maxConnections)
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		sem <- struct{}{}
		go func(c net.Conn) {
			defer func() { <-sem; c.Close() }()
			handleClient(c, absPath, fileInfo.Name(), fileInfo.Size())
		}(conn)
	}
}

func handleClient(conn net.Conn, fullPath, name string, size int64) {
	header := fmt.Sprintf("%s|%d\n", name, size)
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if _, err := conn.Write([]byte(header)); err != nil {
		return
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return
	}
	defer f.Close()

	fmt.Printf("[Host] Sending to %s...\n", conn.RemoteAddr())
	_, err = io.Copy(conn, f)
	if err != nil {
		return
	}

	buf := make([]byte, 4)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n, _ := conn.Read(buf)
	if string(buf[:n]) == "ACK\n" {
		fmt.Printf("[Host] Success: %s received the file\n", conn.RemoteAddr())
	}
}

func broadcast(port int, name string, size int64) {
	addr, _ := net.ResolveUDPAddr("udp", "255.255.255.255"+protocolPort)
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	msg := fmt.Sprintf("%s|%d|%s|%d", signature, port, name, size)
	for {
		conn.Write([]byte(msg))
		time.Sleep(3 * time.Second)
	}
}
