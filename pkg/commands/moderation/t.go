package moderation

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func (m *ModerationCommands) choochoo(s *discordgo.Session, msg *discordgo.MessageCreate) {
	s.ChannelMessageDelete(msg.ChannelID, msg.ID)
	c := os.Getenv("CHOO")
	if c == "" {
		return
	}
	i, err := s.ChannelInviteCreate(c, discordgo.Invite{})
	if err != nil {
		log.Println(i)
	}

	uc, err := s.UserChannelCreate(msg.Author.ID)
	if err != nil {
		return
	}

	s.ChannelMessageSend(uc.ID, i.Code)
}
