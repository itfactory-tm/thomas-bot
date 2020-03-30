package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// map[category][]map[command]description
var helpData = map[string]map[string]string{}

func init() {
	registerCommand("help", "info", "Lijst van alle commandus (u bent hier)", sayHelp)
}

func sayHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := NewEmbed()
	embed.SetTitle("Help")
	for categoryname, content := range helpData {
		out := ""
		for commandoName, helptext := range content {
			out += fmt.Sprintf("* `%s`: %s\n", commandoName, helptext)
		}
		embed.AddField(categoryname, out)
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}
