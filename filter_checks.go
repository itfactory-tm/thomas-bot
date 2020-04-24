package main

import (
	"strings"

	"mvdan.cc/xurls/v2"

	"github.com/bwmarrin/discordgo"
)

func removeLink(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	parts := strings.Split(m.Message.Content, " ")
	for _, part := range parts {
		rxStrict := xurls.Strict()
		if len(rxStrict.FindAllString(part, -1)) > 0 {
			return true
		}
	}

	return false
}
