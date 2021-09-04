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

// Register registers the handlers
func (m *ModerationCommands) Register(registry command.Registry, server command.Server) {
	registry.RegisterMessageCreateHandler("mute", m.muteUser)
	registry.RegisterMessageCreateHandler("unmute", m.unmuteUser)
	registry.RegisterMessageCreateHandler("alert", m.sayAlert)
	registry.RegisterMessageCreateHandler("choochoo", m.choochoo)
	registry.RegisterMessageCreateHandler("clean", m.cleanChannel)
	registry.RegisterMessageCreateHandler("membercount", m.membercount)

	registry.RegisterMessageCreateHandler("", m.checkMessageCreateAsync)
	registry.RegisterMessageEditHandler("", m.checkMessageUpdateAsync)
	registry.RegisterMessageReactionAddHandler(m.checkMessageReactionAddAsync)
	m.server = server
}

// InstallSlashCommands registers the slash commands
func (m *ModerationCommands) InstallSlashCommands(session *discordgo.Session) error {
	return nil
}

func (m *ModerationCommands) checkMessageCreateAsync(s *discordgo.Session, msg *discordgo.MessageCreate) {
	go m.checkMessage(s, msg)

	// check if dm or not
	c, err := s.UserChannelCreate(msg.Author.ID)
	if err != nil {
		return
	}

	if c.ID == msg.ChannelID {
		s.ChannelMessageSend(msg.ChannelID, "Oh my... human language... let me try... 01101001 01100110 00100000 01111001 01101111 01110101 00100000 01100011 01100001 01101110 00100000 01110010 01100101 01100001 01100100 00100000 01110100 01101000 01101001 01110011 00101100 00100000 01100011 01101111 01101110 01110100 01110010 01101001 01100010 01110101 01110100 01100101 00100000 01101111 01101110 00100000 01100111 01101001 01110100 01101000 01110101 01100010...\n\n Oh no, I still cannot understand humans. I'm sorry if you need me type `/` to get a list of things I can do for you!")
		return
	}
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

// Info return the commands in this package
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
		command.Command{
			Name:        "membercount",
			Category:    command.CategoryModeratie,
			Description: "Count the users in all roles",
			Hidden:      false,
		},
	}
}
