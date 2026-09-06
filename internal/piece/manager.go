package piece

type ClientManager struct {
	ClientInterested bool
	ClientChoked     bool
	ClientPieces     map[uint32]bool
	SelectedPiece    uint32
	HasSelectedPiece bool
}

func NewClientManager(pieceCount uint32) *ClientManager {
	return &ClientManager{
		ClientChoked: true,
		ClientPieces: make(map[uint32]bool, pieceCount),
	}
}

func (cm *ClientManager) SelectNeededPiece(peerPieces map[uint32]bool) bool {
	for pieceIndex := range peerPieces {
		if !cm.ClientPieces[pieceIndex] {
			cm.SelectedPiece = pieceIndex
			cm.HasSelectedPiece = true
			return true
		}
	}

	cm.HasSelectedPiece = false
	return false
}
