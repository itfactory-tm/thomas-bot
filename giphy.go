package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	libgiphy "github.com/sanzaru/go-giphy"
)

const discordTalksVragen = "689915740564095061"

func init() {
	registerCommand("clap", "fun", "Applaus!", clap)
	registerCommand("hug", "fun", "Omdat je altijd een knuffel kunt gebruiken", hug)
	registerCommand("cat", "fun", "Voor de kattenmensen", cat)
	registerCommand("dog", "fun", "Voor de honden fans", dog)
	registerCommand("bunny", "fun", "Voor de niet katten of hondenmensen", bunny)
}

func clap(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID == discordTalksVragen {
		go func() {
			voiceQueueChan <- "./sounds/clapping2.wav"
		}()
		if !audioConnected {
			go connectVoice(s)
		}
	}
	postRandomGif(s, m, "applause")
}

func hug(s *discordgo.Session, m *discordgo.MessageCreate) {
	postRandomGif(s, m, "hug")
}

func cat(s *discordgo.Session, m *discordgo.MessageCreate) {
	postRandomGif(s, m, "cat")
}

func dog(s *discordgo.Session, m *discordgo.MessageCreate) {
	postRandomGif(s, m, "dog")
}

func bunny(s *discordgo.Session, m *discordgo.MessageCreate) {
	postRandomGif(s, m, "bunny")
}

func postRandomGif(s *discordgo.Session, m *discordgo.MessageCreate, subject string) {
	if c.GiphyKey == "" {
		s.ChannelMessageSend(m.ChannelID, "Giphy key is lacking from deployment")
		return
	}
	giphy := libgiphy.NewGiphy(c.GiphyKey)
	data, err := giphy.GetRandom(subject)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	embed := NewEmbed()
	embed.SetImage(data.Data.Fixed_height_downsampled_url)
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}
