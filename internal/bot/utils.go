package bot

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
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
