package bot

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
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
		type ErrorData struct {
			Error string `json:"error"`
		}

		var errorData ErrorData
		if err := json.NewDecoder(r.Body).Decode(&errorData); err != nil {
			errorData.Error = "Parsing error"
		}

		logger.Error("Error posting", "url", url, "statusCode", r.StatusCode, "error", errorData.Error)
		return fmt.Errorf("error posting: status code %d", r.StatusCode)
	}

	if err := json.NewDecoder(r.Body).Decode(resp); err != nil {
		logger.Error("Error parsing", "url", url, "err", err)
		return err
	}

	return nil
}

func (b *Bot) RequestAccount(lobbyAddress string) error {
	url := fmt.Sprintf("https://%s", lobbyAddress)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}

	type AccountRequest struct {
		BotId uint64 `json:"bot_id"`
	}

	type AccountResponse struct {
		AccountId uint64 `json:"account_id"`
	}

	var account AccountResponse
	if err := post(client, url+"/bot/account",
		AccountRequest{BotId: b.BotId}, &account, b.logger); err != nil {
		return err
	}

	if account.AccountId == 0 {
		type RegisterRequest struct {
			BotId uint64 `json:"bot_id"`
		}

		type RegisterResponse struct {
			AccountId uint64 `json:"account_id"`
		}

		var register RegisterResponse
		if err := post(client, url+"/bot/register",
			RegisterRequest{BotId: b.BotId}, &register, b.logger); err != nil {
			return err
		}

		b.Account.AccountId = register.AccountId
	} else {
		b.Account.AccountId = account.AccountId
	}

	type CharacterListRequest struct {
		AccountId uint64 `json:"account_id"`
	}

	type CharacterListResponse struct {
		Characters []uint64 `json:"characters"`
	}

	var characters CharacterListResponse
	if err := post(client, url+"/bot/character/list",
		CharacterListRequest{AccountId: b.Account.AccountId}, &characters, b.logger); err != nil {
		return err
	}

	if len(characters.Characters) == 0 {
		//TODO: Create character
	} else {
		b.Account.CharacterId = characters.Characters[0]
	}

	type AuthRequest struct {
		AccountId   uint64 `json:"account_id"`
		CharacterId uint64 `json:"character_id"`
	}

	type AuthResponse struct {
		Token string `json:"token"`
	}

	var auth AuthResponse
	if err := post(client, url+"/bot/auth",
		AuthRequest{AccountId: b.Account.AccountId, CharacterId: b.Account.CharacterId}, &auth, b.logger); err != nil {
		return err
	}

	b.Account.AuthToken = auth.Token
	b.logger.Debug("AccountId", b.Account.AccountId, "Token", b.Account.AuthToken)

	return nil
}
