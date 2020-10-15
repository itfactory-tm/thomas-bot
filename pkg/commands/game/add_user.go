package game

import (
	"fmt"

	"github.com/itfactory-tm/thomas-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/sudo"
)

// adduser contains the bob!adduser and bob!remuser command
type AddUserCommand struct{}

// NewAddUserCommand gives a new AddUserCommand
func NewAddUserCommand() *AddUserCommand {
	return &AddUserCommand{}
}

// Register registers the handlers
func (m *AddUserCommand) Register(registry command.Registry, server command.Server) {
	registry.RegisterMessageCreateHandler("adduser", m.addUser)
	registry.RegisterMessageCreateHandler("remuser", m.remUser)
}

func (m *AddUserCommand) addUser(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if !(sudo.IsItfAdmin(msg.Author.ID) || sudo.IsAdmin(msg.Author.ID)) {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%s is not in the sudoers file. This incident will be reported.", msg.Author.ID))
		return
	}

	matched := msg.Message.Mentions
	if len(matched) < 1 {
		s.ChannelMessageSend(msg.ChannelID, "You need to specify a user")
		return
	}
	user := matched[0].ID

	roles, err := s.GuildRoles(msg.GuildID)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	gameRoleID := ""
	for _, role := range roles {
		if role.Name == "ITF Gamer" {
			gameRoleID = role.ID
		}
	}

	err = s.GuildMemberRoleAdd(msg.GuildID, user, gameRoleID)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	s.ChannelMessageSend(msg.ChannelID, ("User added! <@" + user + ">"))
}

func (m *AddUserCommand) remUser(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if !(sudo.IsItfAdmin(msg.Author.ID) || sudo.IsAdmin(msg.Author.ID)) {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%s is not in the sudoers file. This incident will be reported.", msg.Author.ID))
		return
	}

	matched := msg.Message.Mentions
	if len(matched) < 1 {
		s.ChannelMessageSend(msg.ChannelID, "You need to specify a user")
		return
	}

	user := matched[0].ID

	roles, err := s.GuildRoles(msg.GuildID)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	gameRoleID := ""
	for _, role := range roles {
		if role.Name == "ITF Gamer" {
			gameRoleID = role.ID
		}
	}

	err = s.GuildMemberRoleRemove(msg.GuildID, user, gameRoleID)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	s.ChannelMessageSend(msg.ChannelID, ("User removed! <@" + user + ">"))
}

// Info return the commands in this package
func (m *AddUserCommand) Info() []command.Command {
	return []command.Command{
		command.Command{
			Name:        "adduser",
			Category:    command.CategoryModeratie,
			Description: "Add a user to the ITF Gamer role (ITF admin only)",
			Hidden:      false,
		},
		command.Command{
			Name:        "remuser",
			Category:    command.CategoryModeratie,
			Description: "Remove a user from the ITF Gamer role (ITF admin only)",
			Hidden:      false,
		}}
}
