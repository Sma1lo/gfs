package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	Signature    = "GFS_V2"
	ProtocolPort = "9999"
	BufferSize   = 32 * 1024
	HeaderSize   = 16
)

type HostInfo struct {
	addr string
	name string
	size int64
	hash string
}

func parseMsg(msg string, remote net.Addr) (HostInfo, bool) {
	p := strings.Split(msg, "|")
	if len(p) < 5 || p[0] != Signature {
		return HostInfo{}, false
	}
	ip, _, _ := net.SplitHostPort(remote.String())
	size, _ := strconv.ParseInt(p[3], 10, 64)
	return HostInfo{
		addr: ip + ":" + p[1],
		name: p[2],
		size: size,
		hash: p[4],
	}, true
}

func calculateHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
