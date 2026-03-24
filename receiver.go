package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

type HostInfo struct {
	addr string
	name string
	size int64
}

func runReceiver(dest string) {
	pc, err := net.ListenPacket("udp", protocolPort)
	if err != nil {
		fmt.Println("Error: port 9999 is busy")
		return
	}
	defer pc.Close()

	fmt.Println("Searching for local hosts (3s)...")

	hostsMap := make(map[string]HostInfo)

	deadline := time.Now().Add(3 * time.Second)
	pc.SetReadDeadline(deadline)

	buf := make([]byte, 2048)
	for {
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			break
		}

		msg := string(buf[:n])
		if info, ok := parseMsg(msg, addr); ok {
			hostsMap[info.addr] = info
		}
	}

	if len(hostsMap) == 0 {
		fmt.Println("No active hosts found.")
		return
	}

	indexedHosts := make([]HostInfo, 0, len(hostsMap))
	fmt.Println("\nAvailable files:")
	i := 1
	for _, h := range hostsMap {
		fmt.Printf("[%d] %s (%s) - %d bytes\n", i, h.name, h.addr, h.size)
		indexedHosts = append(indexedHosts, h)
		i++
	}

	fmt.Print("\nEnter ID to download: ")
	var input string
	fmt.Scanln(&input)

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > len(indexedHosts) {
		fmt.Println("Invalid ID. Canceled.")
		return
	}

	target := indexedHosts[idx-1]
	downloadFile(target, dest)
}

func parseMsg(msg string, remote net.Addr) (HostInfo, bool) {
	p := strings.Split(msg, "|")
	if len(p) < 4 || p[0] != signature {
		return HostInfo{}, false
	}

	ip, _, _ := net.SplitHostPort(remote.String())
	targetPort := p[1]
	size, _ := strconv.ParseInt(p[3], 10, 64)

	return HostInfo{
		addr: ip + ":" + targetPort,
		name: p[2],
		size: size,
	}, true
}

func downloadFile(h HostInfo, dest string) {
	conn, err := net.DialTimeout("tcp", h.addr, 5*time.Second)
	if err != nil {
		fmt.Println("Connection failed:", err)
		return
	}
	defer conn.Close()

	r := bufio.NewReader(conn)
	_, _ = r.ReadString('\n')

	_ = os.MkdirAll(dest, 0755)

	finalPath := filepath.Join(dest, h.name)

	if _, err := os.Stat(finalPath); err == nil {
		finalPath = filepath.Join(dest, "new_"+h.name)
	}

	f, err := os.Create(finalPath)
	if err != nil {
		fmt.Println("File error:", err)
		return
	}
	defer f.Close()

	bar := progressbar.NewOptions64(
		h.size,
		progressbar.OptionSetDescription("Downloading "+h.name),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	_, err = io.Copy(io.MultiWriter(f, bar), r)

	if err == nil {
		conn.Write([]byte("ACK\n"))
		fmt.Printf("\nSuccess: %s saved\n", h.name)
	} else {
		fmt.Printf("\nError during download: %v\n", err)
	}
}
