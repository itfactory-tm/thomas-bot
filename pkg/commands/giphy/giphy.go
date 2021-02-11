package giphy

import (
	"fmt"
	"log"
	"os"

	discordha "github.com/meyskens/discord-ha"

	"github.com/itfactory-tm/thomas-bot/pkg/embed"

	libgiphy "github.com/sanzaru/go-giphy"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

//TODO: replace me
const discordTalksVragen = "689915740564095061"
const audioChannel = "688370622228725848"

// GiphyCommands contains the tm!hello command
type GiphyCommands struct {
	server command.Server
}

// NewGiphyCommands gives a new GiphyCommands
func NewGiphyCommands() *GiphyCommands {
	return &GiphyCommands{}
}

// Register registers the handlers
func (g *GiphyCommands) Register(registry command.Registry, server command.Server) {
	registry.RegisterMessageCreateHandler("clap", g.clap)

	registry.RegisterMessageCreateHandler("hug", func(s *discordgo.Session, m *discordgo.MessageCreate) {
		g.postRandomGif(s, m, "hug")
	})
	registry.RegisterMessageCreateHandler("cat", func(s *discordgo.Session, m *discordgo.MessageCreate) {
		g.postRandomGif(s, m, "cat")
	})
	registry.RegisterMessageCreateHandler("dog", func(s *discordgo.Session, m *discordgo.MessageCreate) {
		g.postRandomGif(s, m, "dog")
	})
	registry.RegisterMessageCreateHandler("bunny", func(s *discordgo.Session, m *discordgo.MessageCreate) {
		g.postRandomGif(s, m, "bunny")
	})
	registry.RegisterMessageCreateHandler("honk", func(s *discordgo.Session, m *discordgo.MessageCreate) {
		g.postRandomGif(s, m, "untitled goose game")
	})

	g.server = server
}

// InstallSlashCommands registers the slash commands
func (g *GiphyCommands) InstallSlashCommands(session *discordgo.Session) error {
	return nil
}

func (g *GiphyCommands) clap(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID == discordTalksVragen {
		err := g.server.GetDiscordHA().SendVoiceCommand("thomasbot", discordha.VoiceCommand{
			ChannelID: audioChannel,
			File:      "clappingmono.wav",
			UserID:    m.Author.ID,
		})
		if err != nil {
			log.Printf("Error sending voice command: %q\n", err)
		}
	}
	g.postRandomGif(s, m, "applause")
}

func (g *GiphyCommands) postRandomGif(s *discordgo.Session, m *discordgo.MessageCreate, subject string) {
	// TODO: share config with commands
	if os.Getenv("THOMASBOT_GIPHYKEY") == "" {
		s.ChannelMessageSend(m.ChannelID, "Giphy key is lacking from deployment")
		return
	}
	giphy := libgiphy.NewGiphy(os.Getenv("THOMASBOT_GIPHYKEY"))
	data, err := giphy.GetRandom(subject)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", err))
		return
	}

	embed := embed.NewEmbed()
	embed.SetImage(data.Data.Fixed_height_downsampled_url)
	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}

// Info return the commands in this package
func (g *GiphyCommands) Info() []command.Command {
	return []command.Command{
		command.Command{
			Name:        "clap",
			Category:    command.CategoryFun,
			Description: "Applause!",
			Hidden:      false,
		},
		command.Command{
			Name:        "hug",
			Category:    command.CategoryFun,
			Description: "You can always use a hug",
			Hidden:      false,
		},
		command.Command{
			Name:        "cat",
			Category:    command.CategoryFun,
			Description: "Purrrrfect",
			Hidden:      false,
		},
		command.Command{
			Name:        "dog",
			Category:    command.CategoryFun,
			Description: "Not everyone likes cats",
			Hidden:      false,
		},
		command.Command{
			Name:        "bunny",
			Category:    command.CategoryFun,
			Description: "For those who neither like cats nor dogs",
			Hidden:      false,
		},
		command.Command{
			Name:        "honk",
			Category:    command.CategoryFun,
			Description: "Peace was never an option",
			Hidden:      false,
		},
	}
}
