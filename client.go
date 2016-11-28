package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var serverAddr string

func init() {
	flag.StringVar(&serverAddr, "server-addr", "", "The remote server address")
	flag.Parse()

	if serverAddr == "" {
		log.Fatal("Server addr must be set")
	}
}

func main() {

	// get youtube url
	args := flag.Args()

	if len(args) == 0 {
		log.Fatal("Need a youtube url in args")
	}

	// check if valid youtube url

	// send video url to the server
	err := SendVideo(args[0])

	if err != nil {
		return
	}

}

func SendVideo(url string) error {

	// connect to this socket
	conn, _ := net.Dial("tcp", serverAddr)
	for {
		// fmt.Print("Text to send: ")
		// text, _ := reader.ReadString('\n')

		// send url
		fmt.Fprintf(conn, url+"\n")
		// listen for reply
		// message, _ := bufio.NewReader(conn).ReadString('\n')
		// fmt.Print("Message from server: " + message)
		os.Exit(0)
	}

	return nil
}
