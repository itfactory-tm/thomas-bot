package members

import (
	"bytes"
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/db"
	"github.com/itfactory-tm/thomas-bot/pkg/embed"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

// TODO: replace me
const itfDiscord = "687565213943332875"
const guestRole = "687568536356257890"

// MemberCommands contains the tm!role command and welcome messages
type MemberCommands struct {
	db db.Database
}

// NewMemberCommand gives a new MemberCommands
func NewMemberCommand(conn db.Database) *MemberCommands {
	return &MemberCommands{
		db: conn,
	}
}

// Register registers the handlers
func (m *MemberCommands) Register(registry command.Registry, server command.Server) {
	registry.RegisterMessageCreateHandler("role", m.sayRole)
	registry.RegisterGuildMemberAddHandler(m.onGuildMemberAdd)
	registry.RegisterMessageReactionAddHandler(m.handleRolePermissionReaction)
	registry.RegisterMessageReactionAddHandler(m.handleRoleReaction)
}

func (m *MemberCommands) onGuildMemberAdd(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	conf, err := m.db.ConfigForGuild(g.GuildID)
	if err != nil {
		log.Println(err)
		return
	}

	if conf.WelcomeChannelID == "" {
		// no welcome channel set!
		return
	}

	t, err := template.New("welcome").Parse(conf.WelcomeText)
	if err != nil {
		log.Println(err)
		return
	}
	var welcomeText bytes.Buffer
	err = t.Execute(&welcomeText, g)
	if err != nil {
		log.Println(err)
		return
	}

	welcome, _ := s.ChannelMessageSend(conf.WelcomeChannelID, welcomeText.String())
	go func() {
		// waving back is not essential and should not delay other actions
		// plus the students want to race against the bot in waving at new users so let's give a head start
		time.Sleep(5 * time.Minute)
		err = s.MessageReactionAdd(conf.WelcomeChannelID, welcome.ID, "👋")
		if err != nil {
			log.Println(err)
		}
		err = s.MessageReactionAdd(conf.WelcomeChannelID, welcome.ID, "💗")
		if err != nil {
			log.Println(err)
		}
	}()

	if g.GuildID == itfDiscord {
		m.superITFSpecificStuffWeShouldPutIntoAGeneralThing(s, g)
	}

}

func (m *MemberCommands) superITFSpecificStuffWeShouldPutIntoAGeneralThing(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	err := s.GuildMemberRoleAdd(g.GuildID, g.Member.User.ID, guestRole) // gast role
	if err != nil {
		log.Printf("Cannot set role for user %s: %q\n", g.Member.User.ID, err)
	}

	c, err := s.UserChannelCreate(g.Member.User.ID)
	if err != nil {
		log.Printf("Cannot DM user %s\n", g.Member.User.ID)
		return
	}

	s.ChannelMessageSend(c.ID, fmt.Sprintf("Hello %s", g.User.Username))
	time.Sleep(time.Second)
	s.ChannelMessageSend(c.ID, "Welcome to the ITFactory Discord!")
	time.Sleep(time.Second)
	s.ChannelMessageSend(c.ID, "My name is Thomas Bot, I am a bot who can help you!")
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
	s.ChannelMessageSend(c.ID, "Please set your name for our Discord server to your actual name, this will help us to identify you and let you in! Thank you!")
	time.Sleep(3 * time.Second)

	m.SendRoleDM(s, g.Member.User.ID)
}

// Info return the commands in this package
func (m *MemberCommands) Info() []command.Command {
	return []command.Command{
		command.Command{
			Name:        "role",
			Category:    command.CategoryAlgemeen,
			Description: "Modify your ITFactory Discord role",
			Hidden:      false,
		},
	}
}
