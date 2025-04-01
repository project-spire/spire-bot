package bot

import (
	"spire/protocol"
	"spire/protocol/auth"
)

func (b *Bot) RequestLogin() {
	login := &auth.Login{
		Token: b.Account.AuthToken,
	}
	p := auth.AuthClientProtocol{
		Protocol: &auth.AuthClientProtocol_Login{Login: login},
	}

	buf, err := protocol.SerializeProtocol(protocol.Auth, &p)
	if err != nil {
		b.logger.Error("%v", err)
		b.Stop()
	}

	b.conn.Sender <- buf
}
