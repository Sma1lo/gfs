package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"
)

func runSender(ctx context.Context, path string) {
	absPath, _ := filepath.Abs(path)
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		fmt.Println("File error:", err)
		return
	}

	stopAnim := make(chan bool)
	go animateText("Calculating SHA-256", stopAnim)
	fileHash, _ := calculateHash(absPath)
	stopAnim <- true

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return
	}
	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	port := ln.Addr().(*net.TCPAddr).Port
	fmt.Printf("\r[Host] %s (Port: %d, Hash: %s[:8])\n", fileInfo.Name(), port, fileHash[:8])

	go startBroadcaster(ctx, port, fileInfo.Name(), fileInfo.Size(), fileHash)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		go handleClient(conn, absPath)
	}
}

func handleClient(conn net.Conn, fullPath string) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(1 * time.Minute))

	header := make([]byte, HeaderSize)
	if _, err := io.ReadFull(conn, header); err != nil {
		return
	}

	offset := int64(binary.BigEndian.Uint64(header[:8]))
	length := int64(binary.BigEndian.Uint64(header[8:]))

	f, err := os.Open(fullPath)
	if err != nil {
		return
	}
	defer f.Close()

	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return
	}

	io.CopyN(conn, f, length)
}
