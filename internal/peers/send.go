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
	message := []byte{0, 0, 0, 1, 2} // 4 byte length, 1 byte ID ( = 2), no payload
	_, err := conn.Write(message)
	return err
}

func SendRequest(conn net.Conn, index uint32, begin uint32, length uint32) error {
	message := make([]byte, 17)
	binary.BigEndian.PutUint32(message[0:4], 13)
	message[4] = 6
	binary.BigEndian.PutUint32(message[5:9], index)
	binary.BigEndian.PutUint32(message[9:13], begin)
	binary.BigEndian.PutUint32(message[13:17], length)

	_, err := conn.Write(message)
	return err
}
