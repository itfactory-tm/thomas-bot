package moderation

import (
	"fmt"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/sudo"
)

var userRegex = regexp.MustCompile(`!u?n?mute <(.*)>`)

func (m *ModerationCommands) muteUser(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if !sudo.IsAdmin(msg.Author.ID) {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%s is not in the sudoers file. This incident will be reported.", msg.Author.ID))
		return
	}

	matched := userRegex.FindStringSubmatch(msg.Content)
	if len(matched) <= 1 {
		s.ChannelMessageSend(msg.ChannelID, "You need to specify a user")
		return
	}
	user := matched[1]

	roles, err := s.GuildRoles(msg.GuildID)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	mutedID := ""
	for _, role := range roles {
		if role.Name == "Muted" {
			mutedID = role.ID
		}
	}

	err = s.GuildMemberRoleAdd(msg.GuildID, user[2:], mutedID)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	s.ChannelMessageSend(msg.ChannelID, ":mute:")
}

func (m *ModerationCommands) unmuteUser(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if !sudo.IsAdmin(msg.Author.ID) {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%s is not in the sudoers file. This incident will be reported.", msg.Author.ID))
		return
	}

	matched := userRegex.FindStringSubmatch(msg.Content)
	if len(matched) <= 1 {
		s.ChannelMessageSend(msg.ChannelID, "You need to specify a user")
		return
	}
	user := matched[1]

	roles, err := s.GuildRoles(msg.GuildID)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	mutedID := ""
	for _, role := range roles {
		if role.Name == "Muted" {
			mutedID = role.ID
		}
	}

	err = s.GuildMemberRoleRemove(msg.GuildID, user[2:], mutedID)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	s.ChannelMessageSend(msg.ChannelID, ":speaking_head:")
}
