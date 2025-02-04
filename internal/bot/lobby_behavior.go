package bot

func (b *Bot) RequestAuthTokenAsync(lobbyAddress string) <-chan error {
	errCh := make(chan error, 1)

	go func() {
		defer close(errCh)

	}()

	return errCh
}
