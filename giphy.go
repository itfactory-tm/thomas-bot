package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	libgiphy "github.com/sanzaru/go-giphy"
)

func init() {
	registerCommand("clap", "fun", "Applaus!", clap)
}

func clap(s *discordgo.Session, m *discordgo.MessageCreate) {
	if c.GiphyKey == "" {
		s.ChannelMessageSend(m.ChannelID, "Giphy key is lacking from deployment")
		return
	}
	giphy := libgiphy.NewGiphy(c.GiphyKey)
	data, err := giphy.GetRandom("applause")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	embed := NewEmbed()
	embed.SetImage(data.Data.Fixed_height_downsampled_url)
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}
