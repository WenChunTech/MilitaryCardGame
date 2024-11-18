package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

//go:embed MilitaryCardGameFrontend/dist
var distFs embed.FS

var distDir = "MilitaryCardGameFrontend/dist"
var port = 8080
var host = "localhost"
var url = fmt.Sprintf("http://%s:%d", host, port)

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGTERM)

	go OpenServer()
	go OpenBrowser()

	<-signalChan
}

func OpenServer() {
	_, err := distFs.Open(distDir)
	if err != nil {
		log.Panicf("Error open %s: %v\n", distDir, err)
	}

	log.Printf("Serving %s on HTTP port: %d\n", distDir, port)

	http.Handle("/", http.FileServer(http.Dir(distDir)))
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Panicf("Error starting server: %v\n", err)
	}
}

func OpenBrowser() {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}

	args = append(args, url)
	err := exec.Command(cmd, args...).Start()
	if err != nil {
		log.Panicf("Error open browser: %v\n", err)
	}

	log.Printf("Open browser at %s\n", url)
}
