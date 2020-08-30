package moderation

import (
	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

// ModerationCommands contains the tm!hello command
type ModerationCommands struct {
	server command.Server
}

// NewModerationCommands gives a new ModerationCommands
func NewModerationCommands() *ModerationCommands {
	return &ModerationCommands{}
}

func (m *ModerationCommands) Register(registry command.Registry, server command.Server) {
	registry.RegisterMessageCreateHandler("mute", m.muteUser)
	registry.RegisterMessageCreateHandler("unmute", m.unmuteUser)
	registry.RegisterMessageCreateHandler("alert", m.sayAlert)
	registry.RegisterMessageCreateHandler("choochoo", m.choochoo)
	registry.RegisterMessageCreateHandler("clean", m.cleanChannel)

	registry.RegisterMessageCreateHandler("", m.checkMessageCreateAsync)
	registry.RegisterMessageEditHandler("", m.checkMessageUpdateAsync)
	registry.RegisterMessageReactionAddHandler(m.checkMessageReactionAddAsync)
	m.server = server
}

func (m *ModerationCommands) checkMessageCreateAsync(s *discordgo.Session, msg *discordgo.MessageCreate) {
	go m.checkMessage(s, msg)
}

func (m *ModerationCommands) checkMessageUpdateAsync(s *discordgo.Session, msg *discordgo.MessageUpdate) {
	m2 := &discordgo.MessageCreate{
		msg.Message,
	}

	go m.checkMessage(s, m2)
}

func (m *ModerationCommands) checkMessageReactionAddAsync(s *discordgo.Session, msg *discordgo.MessageReactionAdd) {
	go m.checkReaction(s, msg)
}

func (m *ModerationCommands) Info() []command.Command {
	return []command.Command{
		command.Command{
			Name:        "mute",
			Category:    command.CategoryModeratie,
			Description: "Mute a user (admin only)",
			Hidden:      false,
		},
		command.Command{
			Name:        "unmute",
			Category:    command.CategoryModeratie,
			Description: "Unmute a user (admin only)",
			Hidden:      false,
		},
		command.Command{
			Name:        "choochoo",
			Category:    command.CategoryFun,
			Description: "now you see me, now you don't",
			Hidden:      true,
		},
		command.Command{
			Name:        "alert",
			Category:    command.CategoryModeratie,
			Description: "Send an alert to the moderators (use with care!)",
			Hidden:      false,
		},
		command.Command{
			Name:        "clean",
			Category:    command.CategoryModeratie,
			Description: "Delete all messages in a channel (admin only)",
			Hidden:      false,
		},
	}
}
