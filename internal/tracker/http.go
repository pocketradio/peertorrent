package tracker

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func SendTrackerRequest(
	announceURL string,
	infoHash [20]byte,
	port uint64,
	uploaded uint64,
	downloaded uint64,
	left uint64,
) (TrackerResponse, error) {

	peerID := GeneratePeerID()

	u, err := url.Parse(announceURL)
	if err != nil {
		return TrackerResponse{}, err
	}

	q := u.Query()
	q.Set("info_hash", string(infoHash[:]))
	q.Set("peer_id", peerID)
	q.Set("port", strconv.FormatUint(port, 10))
	q.Set("uploaded", strconv.FormatUint(uploaded, 10))
	q.Set("downloaded", strconv.FormatUint(downloaded, 10))
	q.Set("left", strconv.FormatUint(left, 10))

	u.RawQuery = q.Encode() // rawquery is "" since no query params in the announceURL

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		// log.Fatalf("Error creating request: %v", err)
		return TrackerResponse{}, err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return TrackerResponse{}, err
	}

	inst, err := Decode(resp)
	if err != nil {
		return TrackerResponse{}, err
	}

	if inst.FailureReason != "" {
		return TrackerResponse{}, fmt.Errorf("tracker failure : %s", inst.FailureReason)
	}

	return inst, nil
}

func GeneratePeerID() string {

	key := make([]byte, 20)
	rand.Read(key)
	return string(key)
}
