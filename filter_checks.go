package main

import (
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func removeLink(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	parts := strings.Split(m.Message.Content, " ")
	for _, part := range parts {
		_, err := url.ParseRequestURI(part)
		if err != nil {
			continue
		}

		u, err := url.Parse(part)
		if err != nil || u.Scheme == "" || u.Host == "" {
			continue
		}

		return true
	}

	return false
}
