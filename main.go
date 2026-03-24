package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	asciiArt()
	printUsage()

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
				go runSender(args[2])
			} else if action == "get" {
				dest := "."
				if len(args) >= 3 {
					dest = args[2]
				}
				runReceiver(dest)
			}
		case "exit", "quit":
			fmt.Println("Bye.")
			return
		default:
			fmt.Println("Unknown command. Type 'exit' to quit.")
		}
	}
}

func asciiArt() {
	fmt.Println("\033[36m" + `
   ______ ______ _____ 
  / ____// ____// ___/ 
 / / __ / /_    \__ \  
/ /_/ // __/   ___/ /  
\____//_/     /____/   

 High Performance File Transfer` + "\033[0m\n")
}

func printUsage() {
	fmt.Println("  gfs send <path>  - Start hosting a file")
	fmt.Println("  gfs get <dest>   - Scan and download a file")
	fmt.Println("  exit             - Close application\n")
}
