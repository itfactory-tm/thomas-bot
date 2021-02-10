package shout

import (
	"log"
	"regexp"

	"github.com/itfactory-tm/thomas-bot/pkg/util/voice"

	discordha "github.com/meyskens/discord-ha"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

// ShoutCommand contains the tm!hello command
type ShoutCommand struct {
	server command.Server
}

// NewShoutCommand gives a new ShoutCommand
func NewShoutCommand() *ShoutCommand {
	return &ShoutCommand{}
}

// Register registers the handlers
func (s *ShoutCommand) Register(registry command.Registry, server command.Server) {
	registry.RegisterMessageCreateHandler("shout", s.shout)

	s.server = server
}

// InstallSlashCommands registers the slash commands
func (s *ShoutCommand) InstallSlashCommands(session *discordgo.Session) error {
	return nil
}

var shoutRegex = regexp.MustCompile(`^tm!shout (.*)$`)

func (s *ShoutCommand) shout(sess *discordgo.Session, m *discordgo.MessageCreate) {
	ch, err := voice.FindVoiceUser(sess, m.GuildID, m.Author.ID)
	if err != nil {
		log.Println(err)
		return
	}
	if ch != "715889803937185812" { // sorry that is the other bot!
		matches := shoutRegex.FindAllStringSubmatch(m.Message.Content, -1)
		if len(matches) > 0 && len(matches[0]) > 1 {
			err := s.server.GetDiscordHA().SendVoiceCommand("thomasbot", discordha.VoiceCommand{
				ChannelID: ch,
				File:      matches[0][1] + ".wav",
				UserID:    m.Author.ID,
			})
			if err != nil {
				log.Printf("Error sending voice command: %q\n", err)
			}
		}

	}
}

// Info return the commands in this package
func (s *ShoutCommand) Info() []command.Command {
	return []command.Command{
		command.Command{
			Name:        "shout",
			Category:    command.CategoryFun,
			Description: "Send an audio message",
			Hidden:      false,
		},
	}
}
