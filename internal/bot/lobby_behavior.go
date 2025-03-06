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
		BotId uint64 `json:"bot_id"`
	}

	type MeResponse struct {
		Found     bool   `json:"found"`
		AccountId uint64 `json:"account_id"`
	}

	var account MeResponse
	if err := post(client, url+"/bot/account/me",
		MeRequest{BotId: b.BotId}, &account, b.logger); err != nil {
		return err
	}

	if !account.Found {
		type AccountCreateRequest struct {
			BotId uint64 `json:"bot_id"`
		}

		type AccountCreateResponse struct {
			AccountId uint64 `json:"account_id"`
		}

		var create AccountCreateResponse
		if err := post(client, url+"/bot/account/create",
			AccountCreateRequest{BotId: b.BotId}, &create, b.logger); err != nil {
			return err
		}

		b.Account.AccountId = create.AccountId
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
		type CharacterCreateRequest struct {
			AccountId     uint64 `json:"account_id"`
			CharacterName string `json:"character_name"`
		}

		type CharacterCreateResponse struct {
			CharacterId uint64 `json:"character_id"`
		}

		characterName := fmt.Sprintf("bot_%05d_%d", b.BotId, 1)
		var characterCreate CharacterCreateResponse
		if err := post(client, url+"/bot/character/create",
			CharacterCreateRequest{AccountId: b.Account.AccountId, CharacterName: characterName}, &characterCreate, b.logger); err != nil {
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
	if err := post(client, url+"/bot/account/auth",
		AuthRequest{AccountId: b.Account.AccountId, CharacterId: b.Account.CharacterId}, &auth, b.logger); err != nil {
		return err
	}

	b.Account.AuthToken = auth.Token
	b.logger.Debug("AccountId", b.Account.AccountId, "Token", b.Account.AuthToken)

	return nil
}
