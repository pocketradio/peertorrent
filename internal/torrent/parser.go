package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"strconv"

	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce string `bencode:"announce"`
	Info     Info   `bencode:"info"`
	InfoHash [20]byte
}

type Info struct {
	Name        string `bencode:"name"`
	Pieces      string `bencode:"pieces"`
	PieceLength uint32 `bencode:"piece length"`
	Length      uint64 `bencode:"length"`
}

func FileParser(file []byte) (TorrentFile, error) {

	// getting raw bytes of info alone dict to build the SHA ( for info_hash ). rest of it is unmarshalled and stored in the struct

	infoBytes, err := extractInfoBytes(file)
	if err != nil {
		return TorrentFile{}, err
	}

	reader := bytes.NewReader(file)
	var instance TorrentFile
	if err := bencode.Unmarshal(reader, &instance); err != nil {
		return TorrentFile{}, fmt.Errorf("unmarshal torrent file: %w", err)
	}

	instance.InfoHash = sha1.Sum(infoBytes)

	return instance, nil

}

func extractInfoBytes(data []byte) ([]byte, error) {
	if len(data) == 0 || data[0] != 'd' {
		return nil, fmt.Errorf("torrent must start with a dictionary")
	}

	position := 1
	for position < len(data) && data[position] != 'e' {
		key, next, err := readBencodeString(data, position)
		if err != nil {
			return nil, fmt.Errorf("read torrent key: %w", err)
		}
		position = next

		valueStart := position
		valueEnd, err := scanBencodeValue(data, valueStart)
		if err != nil {
			return nil, fmt.Errorf("read torrent value: %w", err)
		}
		if string(key) == "info" {
			return data[valueStart:valueEnd], nil
		}
		position = valueEnd
	}

	return nil, fmt.Errorf("torrent has no info dictionary")
}

func scanBencodeValue(data []byte, position int) (int, error) {
	if position >= len(data) {
		return 0, fmt.Errorf("unexpected end of data")
	}

	switch data[position] {
	case 'i':
		end := bytes.IndexByte(data[position+1:], 'e')
		if end < 0 {
			return 0, fmt.Errorf("unterminated integer")
		}
		return position + 1 + end + 1, nil
	case 'l':
		position++
		for position < len(data) && data[position] != 'e' {
			var err error
			position, err = scanBencodeValue(data, position)
			if err != nil {
				return 0, err
			}
		}
		if position >= len(data) {
			return 0, fmt.Errorf("unterminated list")
		}
		return position + 1, nil
	case 'd':
		position++
		for position < len(data) && data[position] != 'e' {
			var err error
			_, position, err = readBencodeString(data, position)
			if err != nil {
				return 0, err
			}
			position, err = scanBencodeValue(data, position)
			if err != nil {
				return 0, err
			}
		}
		if position >= len(data) {
			return 0, fmt.Errorf("unterminated dictionary")
		}
		return position + 1, nil
	default:
		_, end, err := readBencodeString(data, position)
		return end, err
	}
}

func readBencodeString(data []byte, position int) ([]byte, int, error) {
	colon := bytes.IndexByte(data[position:], ':')
	if colon < 0 {
		return nil, 0, fmt.Errorf("string has no length separator")
	}
	colon += position
	length, err := strconv.Atoi(string(data[position:colon]))
	if err != nil || length < 0 {
		return nil, 0, fmt.Errorf("invalid string length")
	}
	end := colon + 1 + length
	if end > len(data) {
		return nil, 0, fmt.Errorf("string exceeds input")
	}
	return data[colon+1 : end], end, nil
}
