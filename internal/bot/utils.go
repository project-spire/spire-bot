package bot

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"google.golang.org/protobuf/proto"
)

func post(client *http.Client, url string, req any, resp any, logger *slog.Logger) error {
	data, _ := json.Marshal(req)

	r, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		logger.Error("Error posting", "url", url, "err", err)
		return err
	}
	if r.StatusCode != http.StatusOK {
		logger.Error("Error posting", "url", url, "statusCode", r.StatusCode)
		return errors.New("post error")
	}

	if err := json.NewDecoder(r.Body).Decode(resp); err != nil {
		logger.Error("Error parsing", "url", url, "err", err)
		return err
	}

	return nil
}

func marshalMessage(m proto.Message) ([]byte, error) {
	buf, err := proto.MarshalOptions{}.MarshalAppend(make([]byte, 2), m)
	if err != nil {
		return nil, err
	}

	if 2+len(buf) > 65536 {
		return nil, errors.New("message too large")
	}
	binary.BigEndian.PutUint16(buf[:2], uint16(len(buf)-2))

	return buf, nil
}
