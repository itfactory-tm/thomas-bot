package moderation

import (
	"fmt"
	"log"

	"github.com/itfactory-tm/thomas-bot/pkg/sudo"

	"github.com/bwmarrin/discordgo"
)

func (m *ModerationCommands) cleanChannel(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if !sudo.IsAdmin(msg.Author.ID) {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%s is not in the sudoers file. This incident will be reported.", msg.Author.ID))
		return
	}

	hasMessages := true
	for hasMessages {
		messages, err := s.ChannelMessages(msg.ChannelID, 50, "", "", "")
		if err != nil {
			s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error: %v", err))
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

		err = s.ChannelMessagesBulkDelete(msg.ChannelID, ids)
		if err != nil {
			s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error: %v", err))
			return
		}
	}

	log.Printf("%s has cleared massages in %v", msg.Author.ID, msg.ChannelID)
}
