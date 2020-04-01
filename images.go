package main

import (
	"github.com/bwmarrin/discordgo"
)



func init() {
	registerCommand("erasmus", "links", "Link naar erasmus", sayErasmus)
}

func sayErasmus(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := NewEmbed()
	embed.SetTitle("Erasmus @ ITfactory")
	embed.SetImage("https://static.eyskens.me/thomas-bot/sem_2_2020.gif")
	embed.SetURL("https://thomasmore365.sharepoint.com/sites/james/NL/international?tmbaseCampus=Geel")
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}

func init() {
	registerCommand("partners", "links", "Link naar partners", sayPartners)
}

func sayPartners(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := NewEmbed()
	embed.SetTitle("Partners in education")
	embed.SetImage("https://static.eyskens.me/thomas-bot/voorstelling_partners_in_education.png")
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}

func init() {
	registerCommand("love", "fun", "Toon wat liefde aan elkaar <3", sayLove)
}

func sayLove(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := NewEmbed()
	embed.SetTitle("<3 IT-Factory <3")
	embed.SetImage("https://static.eyskens.me/thomas-bot/love.gif")
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}