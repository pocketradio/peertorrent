package main

import (
	"log"
	"os"
	"peertorrent/internal/torrent"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: client <torrent-file>")
	}

	torrentFile := os.Args[1]
	file, err := os.ReadFile(torrentFile)
	if err != nil {
		log.Fatalf("failed to read torrent file: %v", err)
	}

	_, err = torrent.FileParser(file)
	if err != nil {
		log.Fatalf("failed to parse torrent file: %v", err)
	}
}
