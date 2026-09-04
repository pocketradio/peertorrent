package piece

type ClientManager struct {
	ClientInterested bool
	ClientChoked     bool
	ClientPieces     map[uint32]bool
}

func NewClientManager(pieceCount uint32) *ClientManager {
	return &ClientManager{
		ClientChoked: true,
		ClientPieces: make(map[uint32]bool, pieceCount),
	}
}
