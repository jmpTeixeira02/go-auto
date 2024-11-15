package notifier

import (
	"errors"
	"go-auto/config"
)

type notifyService string

const (
	notifyDiscord  notifyService = "discord"
	notifyTerminal notifyService = "terminal"
)

type Notifier interface {
	SendMessage(string) error
}

func NewNotifier(notifier config.Notifier) (Notifier, error) {
	switch notifyService(notifier.Service) {
	case notifyDiscord:
		return newDiscordNotifier(notifier.Config.Token, notifier.Config.Receiver), nil
	case notifyTerminal:
		return newTerminalNotifier(), nil
	default:
		return nil, errors.New("invalid notifier")
	}
}
