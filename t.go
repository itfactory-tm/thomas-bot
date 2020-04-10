package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

func init() {
	registerCommand(command.Command{
		Name:        "choochoo",
		Category:    command.CategoryFun,
		Description: "now you see me, now you don't",
		Hidden:      true,
		Handler:     choochoo,
	})
}

func choochoo(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageDelete(m.ChannelID, m.ID)
	c := os.Getenv("CHOO")
	if c == "" {
		return
	}
	i, err := s.ChannelInviteCreate(c, discordgo.Invite{})
	if err != nil {
		log.Println(i)
	}

	uc, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		return
	}

	s.ChannelMessageSend(uc.ID, i.Code)
}
