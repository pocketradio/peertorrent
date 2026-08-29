package main

import (
	"log"
	"os"
	"peertorrent/internal/peers"
	"peertorrent/internal/torrent"
	"peertorrent/internal/tracker"
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

	f, err := torrent.FileParser(file)
	if err != nil {
		log.Fatalf("failed to parse torrent file: %v", err)
	}

	peerID := tracker.GeneratePeerID()
	trackerResponse, err := tracker.SendTrackerRequest(
		f.Announce,
		f.InfoHash,
		6881,
		0,
		0,
		f.Info.Length,
		peerID,
	)

	if err != nil {
		log.Fatalf("failed to announce to tracker: %v", err)
	}

	peers.TCPHandshake(f, trackerResponse, peerID)

}
