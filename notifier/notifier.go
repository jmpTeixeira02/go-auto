package notifier

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type DiscordNotifier struct {
	Session *discordgo.Session
	Channel *discordgo.Channel
}

func New(token string, receiver string) DiscordNotifier {

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

	return DiscordNotifier{
		Session: discord,
		Channel: channel,
	}
}

func (d *DiscordNotifier) sendMessage(str string) {
	_, err := d.Session.ChannelMessageSend(d.Channel.ID, str)
	if err != nil {
		panic(fmt.Errorf("Error sending message %w", err))
	}
}
