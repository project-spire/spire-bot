package bot

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func (b *Bot) RequestAccount(lobbyAddress string) error {
	url := fmt.Sprintf("https://%s", lobbyAddress)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}

	type MeRequest struct {
		DevID string `json:"dev_id"`
	}

	type MeResponse struct {
		Found     bool   `json:"found"`
		AccountID uint64 `json:"account_id"`
	}

	devID := fmt.Sprintf("bot_%05d", b.BotId)
	var account MeResponse
	if err := post(client, url+"/account/dev/me",
		MeRequest{DevID: devID}, &account, b.logger); err != nil {
		return err
	}

	if !account.Found {
		type AccountCreateRequest struct {
			DevID string `json:"dev_id"`
		}

		type AccountCreateResponse struct {
			AccountID uint64 `json:"account_id"`
		}

		var create AccountCreateResponse
		if err := post(client, url+"/account/dev/create",
			AccountCreateRequest{DevID: devID}, &create, b.logger); err != nil {
			return err
		}

		b.Account.AccountID = create.AccountID
	} else {
		b.Account.AccountID = account.AccountID
	}

	type CharacterListRequest struct {
		AccountID uint64 `json:"account_id"`
	}

	type CharacterListResponse struct {
		Characters []uint64 `json:"characters"`
	}

	var characters CharacterListResponse
	if err := post(client, url+"/character/list",
		CharacterListRequest{AccountID: b.Account.AccountID}, &characters, b.logger); err != nil {
		return err
	}

	if len(characters.Characters) == 0 {
		type CharacterCreateRequest struct {
			AccountId     uint64 `json:"account_id"`
			CharacterName string `json:"character_name"`
			Race          string `json:"race"`
		}

		type CharacterCreateResponse struct {
			CharacterId uint64 `json:"character_id"`
		}

		characterName := fmt.Sprintf("bot_%05d_%d", b.BotId, 1)
		var characterCreate CharacterCreateResponse
		if err := post(client, url+"/character/create",
			CharacterCreateRequest{
				AccountId:     b.Account.AccountID,
				CharacterName: characterName,
				Race:          "Barbarian",
			}, &characterCreate, b.logger); err != nil {
			return err
		}

		b.Account.CharacterId = characterCreate.CharacterId
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
	if err := post(client, url+"/account/auth",
		AuthRequest{AccountId: b.Account.AccountID, CharacterId: b.Account.CharacterId}, &auth, b.logger); err != nil {
		return err
	}

	b.Account.AuthToken = auth.Token
	b.logger.Debug("AccountID", b.Account.AccountID, "Token", b.Account.AuthToken)

	return nil
}
