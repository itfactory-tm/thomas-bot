package game

import (
	"fmt"

	"github.com/itfactory-tm/thomas-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/sudo"
)

// adduser contains the bob!adduser and bob!remuser command
type UserCommand struct{}

// NewUserCommand gives a new UserCommand
func NewUserCommand() *UserCommand {
	return &UserCommand{}
}

// Register registers the handlers
func (u *UserCommand) Register(registry command.Registry, server command.Server) {
	registry.RegisterMessageCreateHandler("adduser", u.addUser)
	registry.RegisterMessageCreateHandler("remuser", u.remUser)
}

func (u *UserCommand) addUser(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if !sudo.IsItfAdmin(msg.Author.ID) {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%s is not in the sudoers file. This incident will be reported.", msg.Author.ID))
		return
	}

	mentions := msg.Message.Mentions
	if len(mentions) < 1 {
		s.ChannelMessageSend(msg.ChannelID, "You need to specify a user")
		return
	}

	roles, err := s.GuildRoles(msg.GuildID)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error getting guild roles: %v", err))
		return
	}

	gameRoleID := ""
	for _, role := range roles {
		if role.Name == "ITF Gamer" {
			gameRoleID = role.ID
		}
	}

	//Adds role of single or multiple users
	affectedUsers := ""
	for _, user := range mentions {
		userID := user.ID
		err = s.GuildMemberRoleAdd(msg.GuildID, userID, gameRoleID)
		if err != nil {
			s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error adding role: %v", err))
			return
		}
		affectedUsers += fmt.Sprintf("<@%s> ", userID)
	}

	s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("User added! %s", affectedUsers))
}

func (u *UserCommand) remUser(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if !sudo.IsItfAdmin(msg.Author.ID) {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%s is not in the sudoers file. This incident will be reported.", msg.Author.ID))
		return
	}

	mentions := msg.Message.Mentions
	if len(mentions) < 1 {
		s.ChannelMessageSend(msg.ChannelID, "You need to specify a user")
		return
	}

	roles, err := s.GuildRoles(msg.GuildID)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error getting guild roles: %v", err))
		return
	}

	gameRoleID := ""
	for _, role := range roles {
		if role.Name == "ITF Gamer" {
			gameRoleID = role.ID
		}
	}

	//Removes role of single or multiple users
	affectedUsers := ""
	for _, user := range mentions {
		userID := user.ID
		err = s.GuildMemberRoleRemove(msg.GuildID, userID, gameRoleID)
		if err != nil {
			s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error removing role: %v", err))
			return
		}
		affectedUsers += fmt.Sprintf("<@%s> ", userID)
	}

	s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("User removed! %s", affectedUsers))
}

// Info return the commands in this package
func (u *UserCommand) Info() []command.Command {
	return []command.Command{
		command.Command{
			Name:        "adduser",
			Category:    command.CategoryModeratie,
			Description: "Add a user to the ITF Gamer role (ITF Game admin only)",
			Hidden:      false,
		},
		command.Command{
			Name:        "remuser",
			Category:    command.CategoryModeratie,
			Description: "Remove a user from the ITF Gamer role (ITF Game admin only)",
			Hidden:      false,
		}}
}
