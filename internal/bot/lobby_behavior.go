package bot

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

func (b *Bot) RequestAccount(lobbyAddress string) error {
	type AccountRequest struct {
		BotId uint64 `json:"bot_id"`
	}

	type AccountResponse struct {
		AccountId uint64 `json:"account_id"`
	}

	type RegisterRequest struct {
		BotId uint64 `json:"bot_id"`
	}

	type RegisterResponse struct {
		AccountId uint64 `json:"account_id"`
	}

	type AuthRequest struct {
		AccountId uint64 `json:"account_id"`
	}

	type AuthResponse struct {
		Token string `json:"token"`
	}

	url := fmt.Sprintf("https://%s", lobbyAddress)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}

	req, _ := json.Marshal(AccountRequest{BotId: b.BotId})
	resp, err := client.Post(url+"/account/bot", "application/json", bytes.NewBuffer(req))
	if err != nil {
		b.logger.Error("Error requesting account: %v", err)
		return err
	}

	var account AccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		b.logger.Error("Error parsing account response: %v", err)
		return err
	}

	if account.AccountId == 0 {
		req, _ := json.Marshal(RegisterRequest{BotId: b.BotId})
		resp, err := client.Post(url+"/register/bot", "application/json", bytes.NewBuffer(req))
		if err != nil {
			b.logger.Error("Error requesting register: %v", err)
			return err
		}

		var register RegisterResponse
		if err := json.NewDecoder(resp.Body).Decode(&register); err != nil {
			b.logger.Error("Error parsing register response: %v", err)
			return err
		}

		b.Account.AccountId = register.AccountId
	} else {
		b.Account.AccountId = account.AccountId
	}

	req, _ = json.Marshal(AuthRequest{AccountId: b.Account.AccountId})
	resp, err = client.Post(url+"/auth/bot", "application/json", bytes.NewBuffer(req))
	if err != nil {
		b.logger.Error("Error requesting auth: %v", err)
		return err
	}

	var auth AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&auth); err != nil {
		b.logger.Error("Error parsing auth response: %v", err)
		return err
	}

	b.Account.AuthToken = auth.Token
	b.logger.Debug("AccountId: %d, Token: %s", b.Account.AccountId, b.Account.AuthToken)

	return nil
}
