package tracker

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
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

// tracker resps can represent peers in two different formats, so we decode
// them generically first and normalize both forms into their internal peer type.

func Decode(resp *http.Response) (TrackerResponse, error) {
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return TrackerResponse{}, fmt.Errorf("tracker returned http status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TrackerResponse{}, err
	}

	value, err := bencode.Decode(bytes.NewReader(body))
	if err != nil {
		return TrackerResponse{}, err
	}

	root, ok := value.(map[string]any)
	if !ok {
		return TrackerResponse{}, fmt.Errorf("tracker response is not a dictionary")
	}

	response := TrackerResponse{}
	if failure, ok := root["failure reason"].(string); ok {
		response.FailureReason = failure
		return response, nil
	}

	interval, ok := root["interval"]
	if !ok {
		return TrackerResponse{}, fmt.Errorf("tracker response has no interval")
	}
	response.Interval, err = asUint64(interval)
	if err != nil {
		return TrackerResponse{}, fmt.Errorf("invalid tracker interval: %w", err)
	}

	peers, ok := root["peers"]
	if !ok {
		return TrackerResponse{}, fmt.Errorf("tracker response has no peers")
	}
	response.Peers, err = decodePeers(peers)
	if err != nil {
		return TrackerResponse{}, err
	}

	return response, nil
}

func decodePeers(value any) ([]Peer, error) {
	switch peers := value.(type) {

	case string:
		data := []byte(peers)
		if len(data)%6 != 0 { // 6 because each peer is 6 bytes ( 4 for ip, 2 for port )
			return nil, fmt.Errorf("compact peer list length is not divisible by 6")
		}

		result := make([]Peer, 0, len(data)/6)

		for position := 0; position < len(data); position += 6 {
			ip := net.IPv4(data[position], data[position+1], data[position+2], data[position+3])
			port := binary.BigEndian.Uint16(data[position+4 : position+6])
			result = append(result, Peer{IP: ip.String(), Port: uint64(port)})
		}
		return result, nil
	case []any:
		result := make([]Peer, 0, len(peers))
		for _, value := range peers {
			peerMap, ok := value.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("peer entry is not a dictionary")
			}

			ip, ok := peerMap["ip"].(string)
			if !ok {
				return nil, fmt.Errorf("peer entry has invalid ip")
			}
			port, err := asUint64(peerMap["port"])
			if err != nil {
				return nil, fmt.Errorf("peer entry has invalid port: %w", err)
			}
			peer := Peer{IP: ip, Port: port}
			if id, ok := peerMap["peer id"].(string); ok {
				peer.ID = id
			}
			result = append(result, peer)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("tracker peers have unsupported format")
	}
}

func asUint64(value any) (uint64, error) {
	switch number := value.(type) {
	case uint64:
		return number, nil
	case int64:
		if number < 0 {
			return 0, fmt.Errorf("negative value")
		}
		return uint64(number), nil
	default:
		return 0, fmt.Errorf("expected integer")
	}
}
