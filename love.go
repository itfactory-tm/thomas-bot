package main

import (
	"github.com/bwmarrin/discordgo"
)



func init() {
	registerCommand("love", "fun", "Toon wat liefde aan elkaar <3", sayLove)
}

func sayLove(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := NewEmbed()
	embed.SetTitle("<3 IT-Factory <3")
	embed.SetImage("https://static.eyskens.me/thomas-bot/love.gif")
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}