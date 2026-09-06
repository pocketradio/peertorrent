package peers

import (
	"fmt"
	"io"
	"net"
	"peertorrent/internal/piece"
	"peertorrent/internal/torrent"
	"peertorrent/internal/tracker"
	"slices"
	"strconv"
)

func TCPHandshake(tf torrent.TorrentFile, tr tracker.TrackerResponse, clientPeerID string) error {

	address := net.JoinHostPort(tr.Peers[0].IP,
		strconv.FormatUint(tr.Peers[0].Port, 10))

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("Connection error : %s", err)
	}

	defer conn.Close()

	pState := PeerState{
		Choked: true,
	}

	ClientManager := piece.ClientManager{}
	err = PerformHandshake(tf, clientPeerID, tr.Peers[0].ID, conn)

	for pState.Choked {
		err = ReadMessage(conn, &pState, &ClientManager, tf.Info.Pieces)
	}

	return err
}

func PerformHandshake(tf torrent.TorrentFile, clientPeerID string, peerID string, conn net.Conn) error {

	// BEP structure needs 68 bytes
	handshake := []byte{19}
	handshake = append(handshake, []byte("BitTorrent protocol")...) // 19 bytes
	handshake = append(handshake, make([]byte, 8)...)
	handshake = append(handshake, tf.InfoHash[:]...)
	handshake = append(handshake, []byte(clientPeerID)...)

	err := writeAll(conn, handshake)
	if err != nil {
		return fmt.Errorf("error establishing handshake with peer : %s", err)
	}

	bytesBuffer := make([]byte, 68)

	_, err = io.ReadFull(conn, bytesBuffer)
	if err != nil {
		return fmt.Errorf("failed to connect to peer : %s", err)
	}

	if bytesBuffer[0] != 19 || string(bytesBuffer[1:20]) != "BitTorrent protocol" {
		return fmt.Errorf("invalid bittorrent handshake")
	}

	peer_info_hash := bytesBuffer[28:48]
	response_peer_ID := string(bytesBuffer[48:68])

	if !slices.Equal(peer_info_hash, tf.InfoHash[:]) {
		return fmt.Errorf("Info hash did not match : ")
	}

	if response_peer_ID != peerID {
		return fmt.Errorf("peerID did not match.")
	}

	return nil
}

// replacing conn.write with this to ensure partial writes complete
func writeAll(conn net.Conn, data []byte) error {
	for len(data) > 0 {
		n, err := conn.Write(data)
		if err != nil {
			return err
		}
		data = data[n:]
	}
	return nil
}
