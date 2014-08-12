package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
)

var (
	port *int    = flag.Int("port", 5000, "port to open for incoming requests.")
	addr *string = flag.String("addr", "127.0.0.1", "address to bind the server")
)

func exitFromError(err error) {
	fmt.Errorf("Error: %v\nExiting.", err)
	os.Exit(1)
}

func main() {
	flag.Parse()

	directory, err := os.Getwd()
	if err != nil {
		exitFromError(err)
	}

	http.Handle("/", http.FileServer(http.Dir(directory)))

	serverAddress := net.JoinHostPort(*addr, fmt.Sprint(*port))

	fmt.Println("Serving files from:", directory)
	fmt.Println("=>", serverAddress)
	fmt.Println("Press ctrl+c to exit.")

	err = http.ListenAndServe(serverAddress, nil)
	// TODO: implement os signal capture to determine if ERRADDRINUSE, etc
	// EADDRINUSE      = Errno(0x62)
	// http://golang.org/pkg/os/signal/
	// http://golang.org/pkg/syscall/#Errno
	if err != nil {
		exitFromError(err)
	}
}
