package peers

import (
	"encoding/binary"
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
	return sendMessage(conn, 2, nil)
}

func SendRequest(conn net.Conn, index uint32, begin uint32, length uint32) error {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], index)
	binary.BigEndian.PutUint32(payload[4:8], begin)
	binary.BigEndian.PutUint32(payload[8:12], length)

	return sendMessage(conn, 6, payload)
}

// (generic helper fn)

func sendMessage(conn net.Conn, id byte, payload []byte) error {

	message := make([]byte, 5+len(payload)) // length = 4 bytes, ID = 1 byte.
	binary.BigEndian.PutUint32(message[0:4], uint32(1+len(payload)))
	message[4] = id // first 4 bytes for length

	copy(message[5:], payload)

	// in case of req

	for len(message) > 0 {
		n, err := conn.Write(message)
		if err != nil {
			return err
		}
		message = message[n:]
	}
	return nil
}
