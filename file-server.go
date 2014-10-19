package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/nilslice/color"
)

// LoggedRequest is a type that abstracts the file requested and information
// about it to be logged to the console.
type LoggedRequest struct {
	timestamp string
	method    string
	path      string
	protocol  string
}

// NewLoggedRequest creates and returns a correctly formed LoggedRequest
func NewLoggedRequest(req *http.Request) *LoggedRequest {
	t := time.Now().Format(time.StampMilli)
	lr := &LoggedRequest{
		timestamp: t,
		method:    req.Method,
		path:      req.URL.Path,
		protocol:  req.Proto,
	}
	return lr
}

// Log prints the contents of LoggedRequest to the console in a developer friendly
// format. Implement colors for better readability.
func (lr *LoggedRequest) Log() {
	t := "|" + color.MagentaString(lr.timestamp) + "|"
	p := "(" + lr.protocol + ")"
	fmt.Println(
		t,
		color.CyanString(p),
		color.CyanString(lr.method),
		color.YellowString(lr.path),
	)
}

var (
	port = flag.Int("port", 5000, "port to listen on for incoming requests")
	addr = flag.String("addr", "127.0.0.1", "address to bind the server")
)

func exitFromError(err error) {
	fmt.Errorf("Error: %v\nExiting.", err)
	os.Exit(1)
}

func combinedLogAndFileServer(res http.ResponseWriter, req *http.Request, dir string) {
	lr := NewLoggedRequest(req)
	lr.Log()
	http.ServeFile(res, req, dir+req.URL.Path)
}

func main() {
	flag.Parse()

	directory, err := os.Getwd()
	if err != nil {
		exitFromError(err)
	}

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		combinedLogAndFileServer(res, req, directory)
	})

	serverAddress := net.JoinHostPort(*addr, fmt.Sprint(*port))

	fmt.Println(color.YellowString("⎘"), "Serving files from:", color.YellowString(directory))
	fmt.Println(color.GreenString("⤇"), serverAddress)
	fmt.Println(color.RedString("⨂"), "Press 'ctrl+c' to exit.")

	err = http.ListenAndServe(serverAddress, nil)
	// TODO: implement os signal capture to determine if ERRADDRINUSE, etc
	// EADDRINUSE      = Errno(0x62)
	// http://golang.org/pkg/os/signal/
	// http://golang.org/pkg/syscall/#Errno
	if err != nil {
		exitFromError(err)
	}
}
