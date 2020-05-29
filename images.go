package main

import (
	"fmt"
	"math/rand"

	"github.com/itfactory-tm/thomas-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
)

func init() {
	registerCommandDEPRECATED("erasmus", "links", "Link naar erasmus", sayErasmus)
	registerCommandDEPRECATED("partners", "links", "Link naar partners", sayPartners)
	registerCommandDEPRECATED("love", "fun", "Toon wat liefde aan elkaar <3", sayLove)
	registerCommandDEPRECATED("loesje", "fun", "'Een fan van loesje' heeft wijze spreuken", sayLoesje)

	registerCommand(command.Command{
		Name:        "geit",
		Category:    command.CategoryFun,
		Description: "De E-F-blok geiten nu ook online",
		Hidden:      false,
		Handler:     sayGeit,
	})

	registerCommand(command.Command{
		Name:        "paard",
		Category:    command.CategoryFun,
		Description: "De E-F-blok paardjes nu ook online",
		Hidden:      false,
		Handler:     sayPaard,
	})

	registerCommand(command.Command{
		Name:        "schaap",
		Category:    command.CategoryFun,
		Description: "De E-F-blok schapen nu ook online",
		Hidden:      false,
		Handler:     saySchaap,
	})
	
	registerCommand(command.Command{
		Name:        "steun",
		Category:    command.CategoryFun,
		Description: "Examens? We komen er samen door!",
		Hidden:      false,
		Handler:     saySteun,
	})
}

func sayErasmus(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := NewEmbed()
	embed.SetTitle("Erasmus @ ITfactory")
	embed.SetImage("https://static.eyskens.me/thomas-bot/sem_2_2020.gif")
	embed.SetURL("https://thomasmore365.sharepoint.com/sites/james/NL/international?tmbaseCampus=Geel")
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}

func sayPartners(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := NewEmbed()
	embed.SetTitle("Partners in education")
	embed.SetImage("https://static.eyskens.me/thomas-bot/voorstelling_partners_in_education.png")
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}

func sayLove(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := NewEmbed()
	embed.SetTitle("<3 IT-Factory <3")
	embed.SetImage("https://static.eyskens.me/thomas-bot/love.gif")
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}

func sayLoesje(s *discordgo.Session, m *discordgo.MessageCreate) {
	i := rand.Intn(7)
	embed := NewEmbed()
	embed.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/loesje%d.png", i+1))
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}

func sayGeit(s *discordgo.Session, m *discordgo.MessageCreate) {
	i := rand.Intn(9)
	embed := NewEmbed()
	embed.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/geit%d.png", i+1))
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}

func sayPaard(s *discordgo.Session, m *discordgo.MessageCreate) {
	i := rand.Intn(3)
	embed := NewEmbed()
	embed.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/paard%d.png", i+1))
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}

func saySchaap(s *discordgo.Session, m *discordgo.MessageCreate) {
	i := rand.Intn(9)
	embed := NewEmbed()
	embed.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/schaap%d.png", i+1))
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}


func saySteun(s *discordgo.Session, m *discordgo.MessageCreate) {
	i := rand.Intn(40)
	embed := NewEmbed()
	embed.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/examensteun/%02d.png", i+1))
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}
