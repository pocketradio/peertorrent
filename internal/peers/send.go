package peers

import (
	"net"
	"peertorrent/internal/piece"
)

func CheckInterest(conn net.Conn, pState *PeerState, clientMgr *piece.ClientManager) {

	for pieceIndex := range pState.Have {
		if !clientMgr.ClientPieces[pieceIndex] {
			clientMgr.ClientInterested = true
			_ = SendInterested(conn)
			return
		}
	}

}

func SendInterested(conn net.Conn) error {
	message := []byte{0, 0, 0, 1, 2} // 4 byte length, 1 byte ID ( = 2), no payload
	_, err := conn.Write(message)
	return err
}