package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	asciiArt()
	rootCtx, cancelAll := context.WithCancel(context.Background())
	defer cancelAll()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		args := strings.Fields(line)
		cmd := strings.ToLower(args[0])
		switch cmd {
		case "help":
			printUsage()
		case "gfs":
			if len(args) < 2 {
				fmt.Println("Usage: gfs [send|get]")
				continue
			}
			action := strings.ToLower(args[1])
			if action == "send" {
				if len(args) < 3 {
					fmt.Println("Error: specify file path")
					continue
				}
				go runSender(rootCtx, args[2])
			} else if action == "get" {
				dest := "."
				if len(args) >= 3 {
					dest = args[2]
				}
				runReceiver(dest)
			}
		case "exit", "quit":
			cancelAll()
			time.Sleep(200 * time.Millisecond)
			return
		default:
			fmt.Printf("Unknown command: %s\n", cmd)
		}
	}
}
