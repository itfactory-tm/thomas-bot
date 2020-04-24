package main

import (
	"fmt"

	"github.com/itfactory-tm/thomas-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
)

const itfWarroom = "702902280206286858"

func init() {
	registerCommand(command.Command{
		Name:        "alert",
		Category:    command.CategoryModeratie,
		Description: "Verwittig het Discord moderatie team (use with care!)",
		Hidden:      false,
		Handler:     sayAlert,
	})
}

func sayAlert(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Cannot DM user")
		return
	}

	if c.ID == m.ChannelID {
		s.ChannelMessageSend(m.ChannelID, "Cannot alert in DMs")
		return
	}
	s.ChannelMessageDelete(m.ChannelID, m.Message.ID)

	if m.GuildID != itfDiscord {
		s.ChannelMessageSend(c.ID, "Can only alert in ITFactory")
		return
	}

	s.ChannelMessageSend(c.ID, "Alert sent! Thank you.")
	s.ChannelMessageSend(itfWarroom, fmt.Sprintf(":warning: Alert by <@%s> in <#%s>", m.Author.ID, m.ChannelID))
}
