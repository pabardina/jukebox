package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	vlc "github.com/adrg/libvlc-go"
	"github.com/otium/ytdl"
)

func main() {

	var playlist []string

	// choose port with flag

	fmt.Println("Launching server...")

	listener, err := net.Listen("tcp", "0.0.0.0:32401")
	if err != nil {
		log.Fatal("Failed to listen:", err.Error())
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Failed to accept connection:", err.Error())
			continue
		}
		youtubeUrl, _ := bufio.NewReader(conn).ReadString('\n')
		playlist = append(playlist, youtubeUrl)

		go playMusic(youtubeUrl)
	}
}

func playMusic(youtubeUrl string) {

	// example
	// vid, _ := ytdl.GetVideoInfo("https://www.youtube.com/watch?v=d27gTrPPAyk")
	vid, _ := ytdl.GetVideoInfo(youtubeUrl)

	formats := vid.Formats

	toto := formats.Extremes(ytdl.FormatAudioBitrateKey, strings.HasPrefix("best-audio", "best"))

	url, err := vid.GetDownloadURL(toto[0])
	if err != nil {
		fmt.Print(err)
	}

	if err := vlc.Init("--no-video", "--quiet"); err != nil {
		log.Fatal(err)
	}
	defer vlc.Release()

	// Create a new player
	player, err := vlc.NewPlayer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		player.Stop()
		player.Release()
	}()

	// Set player media. The second parameter of the method specifies if
	// the media resource is local or remote.

	err = player.SetMedia(url.String(), false)
	if err != nil {
		log.Fatal(err)
	}

	// Play
	err = player.Play()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Start playing video")

	// Wait some amount of time for the media to start playing
	time.Sleep(1 * time.Second)

	// If the media played is a live stream the length will be 0
	length, err := player.MediaLength()
	if err != nil || length == 0 {
		length = 1000 * 60
	}

	time.Sleep(time.Duration(length) * time.Millisecond)

}
