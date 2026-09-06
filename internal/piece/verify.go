package piece

import (
	"bytes"
	"crypto/sha1"
)

func VerifyPiece(pieceData []byte, expectedHash []byte) bool {
	calculatedHash := sha1.Sum(pieceData)
	return bytes.Equal(calculatedHash[:], expectedHash)
}
