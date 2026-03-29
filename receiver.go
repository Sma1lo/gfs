package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

func runReceiver(dest string) {
	addr, _ := net.ResolveUDPAddr("udp", ":"+ProtocolPort)
	pc, err := net.ListenUDP("udp", addr)
	if err != nil {
		return
	}
	defer pc.Close()

	stopAnim := make(chan bool)
	go animateText("Searching for local hosts", stopAnim)

	hosts := make(map[string]HostInfo)
	pc.SetReadDeadline(time.Now().Add(3 * time.Second))
	buf := make([]byte, 2048)

	for {
		n, remoteAddr, err := pc.ReadFrom(buf)
		if err != nil {
			break
		}
		if info, ok := parseMsg(string(buf[:n]), remoteAddr); ok {
			hosts[info.addr] = info
		}
	}
	stopAnim <- true

	if len(hosts) == 0 {
		fmt.Println("No active hosts found.")
		return
	}

	indexed := []HostInfo{}
	fmt.Println("\nAvailable files:")
	for _, h := range hosts {
		fmt.Printf("[%d] %s (%d bytes) - %s\n", len(indexed)+1, h.name, h.size, h.addr)
		indexed = append(indexed, h)
	}

	fmt.Print("\nSelect ID: ")
	var input string
	fmt.Scanln(&input)
	idx, _ := strconv.Atoi(input)

	if idx > 0 && idx <= len(indexed) {
		downloadFile(indexed[idx-1], dest)
	}
}

func downloadFile(h HostInfo, dest string) {
	_ = os.MkdirAll(dest, 0755)
	finalPath := filepath.Join(dest, h.name)

	f, _ := os.OpenFile(finalPath, os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()

	currentSize := int64(0)
	if stat, err := os.Stat(finalPath); err == nil {
		currentSize = stat.Size()
	}
	if currentSize > h.size {
		f.Truncate(h.size)
	} else if currentSize < h.size {
		f.Truncate(h.size)
	}

	numThreads := 4
	chunkSize := h.size / int64(numThreads)
	var wg sync.WaitGroup
	bar := newProgressBar(h.size, "Downloading")

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		start := int64(i) * chunkSize
		end := start + chunkSize
		if i == numThreads-1 {
			end = h.size
		}
		go func(s, e int64) {
			defer wg.Done()
			downloadChunk(h.addr, s, e, f, bar)
		}(start, end)
	}
	wg.Wait()

	stopHash := make(chan bool)
	fmt.Print("\n")
	go animateText("Verifying SHA-256", stopHash)
	newHash, _ := calculateHash(finalPath)
	stopHash <- true

	if newHash == h.hash {
		fmt.Printf("\rSuccess. Hash Match: %s\n", newHash[:8])
	} else {
		fmt.Printf("\rFailure. Hash Mismatch!\n")
	}
}

func downloadChunk(addr string, start, end int64, outF *os.File, bar *progressbar.ProgressBar) {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(2 * time.Minute))

	header := make([]byte, HeaderSize)
	binary.BigEndian.PutUint64(header[:8], uint64(start))
	binary.BigEndian.PutUint64(header[8:], uint64(end-start))
	conn.Write(header)

	buf := make([]byte, BufferSize)
	pos := start
	for pos < end {
		n, err := conn.Read(buf)
		if n > 0 {
			outF.WriteAt(buf[:n], pos)
			pos += int64(n)
			bar.Add(n)
		}
		if err != nil {
			break
		}
	}
}
