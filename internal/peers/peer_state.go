package peers

type PeerState struct {
	Interested bool
	Choked     bool
	Have       map[uint32]bool
}
