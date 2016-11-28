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

	"os"

	ui "github.com/gizak/termui"
)

var currentSong Song

func main() {

	// choose port with flag

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	// handle key q pressing
	ui.Handle("/sys/kbd/q", func(ui.Event) {
		// press q to quit
		os.Exit(0)
	})

	go func() {
		ui.Loop()
	}()

	fmt.Println("Launching server...")

	listener, err := net.Listen("tcp", "0.0.0.0:32401")
	if err != nil {
		log.Fatal("Failed to listen:", err.Error())
	}

	defer listener.Close()

	playlist := make(chan Song)
	go playMusic(playlist)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Failed to accept connection:", err.Error())
			continue
		}

		youtubeUrl, _ := bufio.NewReader(conn).ReadString('\n')

		// playlist = append(playlist, youtubeUrl)

		newSong := NewSong(youtubeUrl)

		playlist <- newSong

		// time.Sleep(5 * time.Second)

		// stop function
		// newSong.Next()

	}
}

func playMusic(playlist chan Song) {

	// example
	// vid, _ := ytdl.GetVideoInfo("https://www.youtube.com/watch?v=d27gTrPPAyk")

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

	for {
		song := <-playlist

		currentSong = song

		vid, _ := ytdl.GetVideoInfo(song.URL)

		formats := vid.Formats

		streams := formats.Extremes(ytdl.FormatAudioBitrateKey, strings.HasPrefix("best-audio", "best"))

		url, err := vid.GetDownloadURL(streams[0])
		if err != nil {
			fmt.Print(err)
		}

		err = player.SetMedia(url.String(), false)
		if err != nil {
			log.Fatal(err)
		}

		// Play
		err = player.Play()
		if err != nil {
			log.Fatal(err)
		}

		g := ui.NewGauge()
		g.Percent = 0
		g.Width = 50
		g.Height = 3
		g.Y = 11
		g.BorderLabel = vid.Title
		g.BarColor = ui.ColorGreen
		g.BorderFg = ui.ColorWhite
		g.BorderLabelFg = ui.ColorCyan

		ui.Render(g)

		// handle a 1s timer
		ui.Handle("/timer/1s", func(e ui.Event) {
			duration, _ := player.MediaPosition()
			g.Percent = int(duration * 100)
			ui.Render(g)
		})

		// Wait some amount of time for the media to start playing
		time.Sleep(1 * time.Second)

		go func() {
			time.Sleep(vid.Duration)
			song.Next()
		}()
		song.WaitForNext()
		player.Stop()
	}

}

type Song struct {
	URL  string
	next chan bool
}

func NewSong(url string) Song {
	return Song{
		URL:  url,
		next: make(chan bool, 1),
	}
}

func (s Song) Next() {
	s.next <- true
}

func (s Song) WaitForNext() {
	<-s.next
}
