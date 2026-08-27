package tracker

import (
	"bytes"
	"io"
	"net/http"

	"github.com/jackpal/bencode-go"
)

type TrackerResponse struct {
	FailureReason string
	Interval      uint64
	Peers         []Peer
}

type Peer struct {
	ID   string
	IP   string
	Port uint64
}

func Decode(resp *http.Response) (TrackerResponse, error) {

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.Fatalf("error reading tracker response : %v", err)
		return TrackerResponse{}, err
	}

	defer resp.Body.Close()

	inst := TrackerResponse{}
	reader := bytes.NewReader(body)

	err = bencode.Unmarshal(reader, &inst)
	if err != nil {
		return TrackerResponse{}, err
	}

	return inst, nil
}
