package moderation

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const itfDiscord = "687565213943332875"
const itfWarroom = "702902280206286858"

func (m *ModerationCommands) sayAlert(s *discordgo.Session, msg *discordgo.MessageCreate) {
	c, err := s.UserChannelCreate(msg.Author.ID)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, "Cannot DM user")
		return
	}

	if c.ID == msg.ChannelID {
		s.ChannelMessageSend(msg.ChannelID, "Cannot alert in DMs")
		return
	}
	s.ChannelMessageDelete(msg.ChannelID, msg.Message.ID)

	if msg.GuildID != itfDiscord {
		s.ChannelMessageSend(c.ID, "Can only alert in ITFactory")
		return
	}

	s.ChannelMessageSend(c.ID, "Alert sent! Thank you.")
	s.ChannelMessageSend(itfWarroom, fmt.Sprintf(":warning: Alert by <@%s> in <#%s>", msg.Author.ID, msg.ChannelID))
}
