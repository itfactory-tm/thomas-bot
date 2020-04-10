package main

import (
	"fmt"
	"regexp"

	"github.com/itfactory-tm/thomas-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
)

var userRegex = regexp.MustCompile(`!u?n?mute <(.*)>`)

func init() {
	registerCommand(command.Command{
		Name:        "mute",
		Category:    command.CategoryModeratie,
		Description: "Een gebruiker muten (admin only)",
		Hidden:      false,
		Handler:     muteUser,
	})
	registerCommand(command.Command{
		Name:        "unmute",
		Category:    command.CategoryModeratie,
		Description: "Een gebruiker unmuten (admin only)",
		Hidden:      false,
		Handler:     unmuteUser,
	})
}

func muteUser(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !isAdmin(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is not in the sudoers file. This incident will be reported.", m.Author.ID))
		return
	}

	matched := userRegex.FindStringSubmatch(m.Content)
	if len(matched) <= 1 {
		s.ChannelMessageSend(m.ChannelID, "Je moet een gebruiker opgeven")
		return
	}
	user := matched[1]

	roles, err := s.GuildRoles(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	mutedID := ""
	for _, role := range roles {
		if role.Name == "Muted" {
			mutedID = role.ID
		}
	}

	err = s.GuildMemberRoleAdd(m.GuildID, user[2:], mutedID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, ":mute:")
}

func unmuteUser(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !isAdmin(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is not in the sudoers file. This incident will be reported.", m.Author.ID))
		return
	}

	matched := userRegex.FindStringSubmatch(m.Content)
	if len(matched) <= 1 {
		s.ChannelMessageSend(m.ChannelID, "Je moet een gebruiker opgeven")
		return
	}
	user := matched[1]

	roles, err := s.GuildRoles(m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	mutedID := ""
	for _, role := range roles {
		if role.Name == "Muted" {
			mutedID = role.ID
		}
	}

	err = s.GuildMemberRoleRemove(m.GuildID, user[2:], mutedID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, ":speaking_head:")
}
