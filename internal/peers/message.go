package peers

import (
	"encoding/binary"
	"io"
	"net"
)

type Message struct {

	// based on bep3 spec
	Length  uint32
	ID      byte
	Payload []byte
}

func ReadMessage(conn net.Conn, pState *PeerState) error {

	lengthBuffer := make([]byte, 4)

	_, err := io.ReadFull(conn, lengthBuffer)
	if err != nil {
		return err
	}

	msgLength := binary.BigEndian.Uint32(lengthBuffer)

	if msgLength == 0 { // keepalive ; no ID or payload
		return nil // tcp conn does not close
	}

	payloadAndIDBuffer := make([]byte, msgLength) // payload = length - 1 (+ 1byte for the ID)
	_, err = io.ReadFull(conn, payloadAndIDBuffer)

	if err != nil {
		return err
	}

	ID := payloadAndIDBuffer[0]
	payload := payloadAndIDBuffer[1:]
	HandleMessage(conn, ID, payload, pState)

	return nil
}

func HandleMessage(conn net.Conn, ID byte, payload []byte, pState *PeerState) {

	switch ID {
	case 0: // choke
		pState.Choked = true

	case 1:
		pState.Choked = false //

	case 2:
		pState.Interested = true // peer is interested

	case 3:
		pState.Interested = false // peer is not interested

	case 4: // 'have' ; announces a single newly completed piece.
		// eg. "i just got piece xyz"

		if len(payload) != 4 { // 4 since
			return
		}
		if pState.Have == nil {
			pState.Have = make(map[uint32]bool)
		}
		pieceIndex := binary.BigEndian.Uint32(payload)
		pState.Have[pieceIndex] = true

	case 5:
		for byteIndex, byteValue := range payload {
			for bit := 0; bit < 8; bit++ {
				bitValue := (byteValue >> (7 - bit)) & 1
				pieceIndex := byteIndex*8 + bit

				if bitValue == 1 {
					pState.Have[uint32(pieceIndex)] = true
				}
				// absence of key means peer doesnt have the piece
			}
		}
	}

}
