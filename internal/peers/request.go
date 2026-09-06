package peers

import (
	"encoding/binary"
	"fmt"
	"net"

	"peertorrent/internal/piece"
)

func RequestNextBlock(conn net.Conn, pState *PeerState, clientMgr *piece.ClientManager, pieceHashes string) error {

	if !clientMgr.HasSelectedPiece {
		return fmt.Errorf("piece doesn't exist. ")
	}

	reqByte := make([]byte, 17)

	binary.BigEndian.PutUint32(reqByte[0:4], 13)
	reqByte[4] = 6 // ID

	// payload
	binary.BigEndian.PutUint32(reqByte[5:9], clientMgr.SelectedPiece)
	binary.BigEndian.PutUint32(reqByte[9:13], 0)
	binary.BigEndian.PutUint32(reqByte[13:17], 16384)

	err := writeAll(conn, reqByte)
	if err != nil {
		return err
	}

	err = ReadMessage(conn, pState, clientMgr, pieceHashes)
	if err != nil {
		return err
	}

	return nil
}
