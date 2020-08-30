package members

import (
	"fmt"
	"log"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/embed"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

// TODO: replace me
const itfDiscord = "687565213943332875"
const itfWelcome = "687588438886842373"
const guestRole = "687568536356257890"

// MemberCommands contains the tm!role command and welcome messages
type MemberCommands struct{}

// NewMemberCommands gives a new MemberCommands
func NewMemberCommand() *MemberCommands {
	return &MemberCommands{}
}

func (m *MemberCommands) Register(registry command.Registry, server command.Server) {
	registry.RegisterMessageCreateHandler("role", m.SayRole)
	registry.RegisterGuildMemberAddHandler(m.onGuildMemberAdd)
	registry.RegisterMessageReactionAddHandler(m.HandleRolePermissionReaction)
	registry.RegisterMessageReactionAddHandler(m.HandleRoleReaction)
}

func (m *MemberCommands) onGuildMemberAdd(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	if g.GuildID != itfDiscord {
		return
	}

	err := s.GuildMemberRoleAdd(g.GuildID, g.Member.User.ID, guestRole) // gast role
	if err != nil {
		log.Printf("Cannot set role for user %s: %q\n", g.Member.User.ID, err)
	}

	s.ChannelMessageSend(itfWelcome, fmt.Sprintf("Welcone <@%s> to the **IT Factory Official** Discord server. We will send you a DM in a moment to get you set up!", g.User.ID))

	c, err := s.UserChannelCreate(g.Member.User.ID)
	if err != nil {
		log.Printf("Cannot DM user %s\n", g.Member.User.ID)
		return
	}

	s.ChannelMessageSend(c.ID, fmt.Sprintf("Hello %s", g.User.Username))
	time.Sleep(time.Second)
	s.ChannelMessageSend(c.ID, "Welcome to the ITFactory Discord!")
	time.Sleep(time.Second)
	s.ChannelMessageSend(c.ID, "My name is Thomas Bot, i am a bot who can help you!")
	time.Sleep(time.Second)
	s.ChannelMessageSend(c.ID, "New to Discord? No problem we got a manual for you: https://itf.to/discord-help")
	embed := embed.NewEmbed()
	embed.SetImage("https://static.eyskens.me/thomas-bot/opendeurdag-1.png")
	embed.SetURL("https://itf.to/discord-help")
	s.ChannelMessageSendEmbed(c.ID, embed.MessageEmbed)

	time.Sleep(time.Second)
	s.ChannelMessageSend(c.ID, "If you need help just type tm!help")
	time.Sleep(time.Second)
	s.ChannelMessageSend(c.ID, "Warning, i am only able to reply to messages starting with `tm!`, not to normal questions.")
	time.Sleep(5 * time.Second)
	s.ChannelMessageSend(c.ID, "")

	m.SendRoleDM(s, g.Member.User.ID)
}

func (h *MemberCommands) Info() []command.Command {
	return []command.Command{
		command.Command{
			Name:        "role",
			Category:    command.CategoryAlgemeen,
			Description: "Modify your ITFactory Discord role",
			Hidden:      false,
		},
	}
}
