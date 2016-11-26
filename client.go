package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	// get youtube url
	args := os.Args[1:]

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
	conn, _ := net.Dial("tcp", "127.0.0.1:32401")
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
