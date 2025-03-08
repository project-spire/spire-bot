package bot

import (
	"spire/bot/gen/protocol"
	"spire/bot/gen/protocol/auth"
)

func (b *Bot) RequestLogin() {
	login := &auth.Login{
		Role:        auth.Role_Player,
		AccountId:   b.Account.AccountId,
		CharacterId: b.Account.CharacterId,
		Token:       b.Account.AuthToken,
	}

	p := &protocol.AuthProtocol{
		Protocol: &protocol.AuthProtocol_Login{Login: login},
	}

	buf, err := marshalMessage(p)
	if err != nil {
		b.logger.Error("%v", err)
		b.Stop()
	}

	b.conn.Sender <- buf
}
