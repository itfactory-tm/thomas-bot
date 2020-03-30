package main

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	registerCommand("hello", sayHello)
}

func sayHello(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Beep bop boop! Ik ben Thomas Bot, fork me on GitHub!")
}
