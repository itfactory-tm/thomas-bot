package main

import (
	"fmt"
	"log"

	"github.com/itfactory-tm/thomas-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
)

func init() {
	registerCommand(command.Command{
		Name:        "clean",
		Category:    command.CategoryModeratie,
		Description: "Een channel leeghalen (admin only)",
		Hidden:      false,
		Handler:     cleanChannel,
	})
}

func cleanChannel(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !isAdmin(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is not in the sudoers file. This incident will be reported.", m.Author.ID))
		return
	}

	hasMessages := true
	for hasMessages {
		messages, err := s.ChannelMessages(m.ChannelID, 50, "", "", "")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
			return
		}

		ids := []string{}
		for _, message := range messages {
			ids = append(ids, message.ID)
		}

		if len(ids) == 0 {
			hasMessages = false
			break
		}

		err = s.ChannelMessagesBulkDelete(m.ChannelID, ids)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
			return
		}
	}

	log.Printf("%s has cleared massages in %v", m.Author.ID, m.ChannelID)
}
