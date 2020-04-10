package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

func init() {
	registerCommand(command.Command{
		Name:        "hello",
		Category:    command.CategoryFun,
		Description: "Zeg hallo",
		Hidden:      false,
		Handler:     sayHello,
	})
}

func sayHello(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Beep bop boop! Ik ben Thomas Bot, fork me on GitHub!")
}
