// This is a simple TCP server that listens on a port and responds to any
// connection with "Hello, World!" after sleeping for a few seconds. It
// is used to test the TCP stats receiver by allowing us to simulate
// connections that are open for a long time and grow the TCP queue.

package main

import (
	"flag"
	"io"
	"log"
	"net"
	"time"
)

var listenAddr = flag.String("listenAddr", "127.0.0.1:8005", "server listen address")
var sleepSeconds = flag.Int("sleepSeconds", 5, "number of seconds to sleep before responding")

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp4", *listenAddr)
	if err != nil {
		log.Fatalf("Failed to bind to port: %v", err)
	}

	log.Println("Listening on " + *listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		log.Println("Message received from " + conn.RemoteAddr().String())

		// Sleep for a few seconds before responding
		log.Printf("Sleeping for %d seconds before responding", *sleepSeconds)
		time.Sleep(time.Duration(*sleepSeconds) * time.Second)

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	_, err := io.WriteString(conn, "Hello, World!\n")
	if err != nil {
		log.Printf("Failed to write to connection: %v", err)
		return
	}

	err = conn.Close()
	if err != nil {
		log.Printf("Failed to close connection: %v", err)
		return
	}

	log.Println("Connection handled successfully")
}
