package main

import (
	"fmt"
	"log"

	"github.com/itfactory-tm/thomas-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
	libgiphy "github.com/sanzaru/go-giphy"
)

const discordTalksVragen = "689915740564095061"

func init() {
	registerCommand(command.Command{
		Name:        "clap",
		Category:    command.CategoryFun,
		Description: "Applaus!",
		Hidden:      false,
		Handler:     clap,
	})
	registerGiphyCommand("hug", "Omdat je altijd een knuffel kunt gebruiken", "hug")
	registerGiphyCommand("cat", "Voor de kattenmensen", "cat")
	registerGiphyCommand("dog", "Voor de honden fans", "dog")
	registerGiphyCommand("bunny", "Voor de niet katten of hondenmensen", "bunny")
}

func registerGiphyCommand(name, description, keyword string) {
	registerCommand(command.Command{
		Name:        name,
		Category:    command.CategoryFun,
		Description: description,
		Hidden:      false,
		Handler: func(s *discordgo.Session, m *discordgo.MessageCreate) {
			postRandomGif(s, m, keyword)
		},
	})
}

func clap(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID == discordTalksVragen {
		connected := make(chan struct{}, 1)
		go connectVoice(s, connected)

		go func(connected chan struct{}) {
			<-connected // wait for audio to connect

			err := ha.SendVoiceCommand(audioChannel, "./sounds/clapping2.wav")
			if err != nil {
				log.Println(err)
			}
		}(connected)
	}
	postRandomGif(s, m, "applause")
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
