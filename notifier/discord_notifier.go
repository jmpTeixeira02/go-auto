package notifier

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type discordNotifier struct {
	Session *discordgo.Session
	Channel *discordgo.Channel
}

func newDiscordNotifier(token string, receiver string) Notifier {
	if token == "" || receiver == "" {
		panic(errors.New("token and receiver must be configured"))
	}
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(fmt.Errorf("error creating Discord session %w", err))
	}

	err = discord.Open()
	if err != nil {
		panic(fmt.Errorf("error opening connection %w", err))
	}
	channel, err := discord.UserChannelCreate(receiver)
	if err != nil {
		panic(fmt.Errorf("error creating DM channel %w", err))
	}

	return discordNotifier{
		Session: discord,
		Channel: channel,
	}
}

func (d discordNotifier) SendMessage(msg string) error {
	_, err := d.Session.ChannelMessageSend(d.Channel.ID, msg)
	if err != nil {
		return fmt.Errorf("error sending message %w", err)
	}
	return nil
}
