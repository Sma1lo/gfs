package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

func startBroadcaster(ctx context.Context, port int, fileName string, fileSize int64, fileHash string) {
	msg := []byte(fmt.Sprintf("%s|%d|%s|%d|%s", Signature, port, fileName, fileSize, fileHash))
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ifaces, _ := net.Interfaces()
			for _, iface := range ifaces {
				if iface.Flags&net.FlagBroadcast != 0 && iface.Flags&net.FlagUp != 0 {
					addrs, _ := iface.Addrs()
					for _, addr := range addrs {
						if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
							broadcast := make(net.IP, len(ipnet.IP.To4()))
							for i := range broadcast {
								broadcast[i] = ipnet.IP.To4()[i] | ^ipnet.Mask[i]
							}
							udpAddr, _ := net.ResolveUDPAddr("udp", broadcast.String()+":"+ProtocolPort)
							conn, err := net.DialUDP("udp", nil, udpAddr)
							if err == nil {
								conn.Write(msg)
								conn.Close()
							}
						}
					}
				}
			}
		}
	}
}
