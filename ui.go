package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

func asciiArt() {
	fmt.Println("\033[36m" + `
   ______ ______ _____ 
  / ____// ____// ___/ 
 / / __ / /_    \__ \  
/ /_/ // __/   ___/ /  
\____//_/     /____/ 

High Performance File Transfer
` + "\033[0m")
}

func printUsage() {
	fmt.Println("  gfs send <path>  - Host a file")
	fmt.Println("  gfs get <dest>   - Scan and download")
	fmt.Println("  exit             - Close\n")
}

func animateText(text string, stop chan bool) {
	frames := []string{".  ", ".. ", "...", "   "}
	i := 0
	for {
		select {
		case <-stop:
			fmt.Print("\r" + strings.Repeat(" ", len(text)+10) + "\r")
			return
		default:
			fmt.Printf("\r%s %s", text, frames[i%len(frames)])
			i++
			time.Sleep(400 * time.Millisecond)
		}
	}
}

func newProgressBar(size int64, description string) *progressbar.ProgressBar {
	return progressbar.NewOptions64(
		size,
		progressbar.OptionSetDescription(description),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(20),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "█",
			SaucerPadding: "░",
			BarStart:      "|",
			BarEnd:        "|",
		}),
	)
}
