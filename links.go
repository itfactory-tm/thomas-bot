package main

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	registerCommand("website", sayWebsite)
}

func sayWebsite(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Bezoek onze website: https://thomasmore.be/opleidingen/professionele-bachelor/it-factory")
}
