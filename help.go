package main

import (
	"github.com/bwmarrin/discordgo"
)

// map[category][]map[command]description
var helpData = map[string]map[string]string{}

func init() {
	registerCommand("help", sayHelp)
}

func sayHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Bezoek onze website: https://thomasmore.be/opleidingen/professionele-bachelor/it-factory")
}
