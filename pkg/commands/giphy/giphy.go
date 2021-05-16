package giphy

import (
	"fmt"
	"log"
	"os"

	"github.com/itfactory-tm/thomas-bot/pkg/util/slash"

	"github.com/itfactory-tm/thomas-bot/pkg/util/voice"

	discordha "github.com/meyskens/discord-ha"

	"github.com/itfactory-tm/thomas-bot/pkg/embed"

	libgiphy "github.com/sanzaru/go-giphy"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

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
	registry.RegisterInteractionCreate("gif", g.slashCommand)
	registry.RegisterInteractionCreate("clap", g.clap)

	g.server = server
}

// InstallSlashCommands registers the slash commands
func (g *GiphyCommands) InstallSlashCommands(session *discordgo.Session) error {
	err := slash.InstallSlashCommand(session, "", discordgo.ApplicationCommand{
		Name:        "gif",
		Description: "Posts a GIF",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "name of the GIF",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "hug",
						Value: "hug",
					},
					{
						Name:  "cat",
						Value: "cat",
					},
					{
						Name:  "dog",
						Value: "dog",
					},
					{
						Name:  "bunny",
						Value: "bunny",
					},
					{
						Name:  "bunny",
						Value: "bunny",
					},
					{
						Name:  "thumbsup",
						Value: "thumbsup",
					},
				},
			},
		},
	})

	if err != nil {
		return fmt.Errorf("error installing /gif %w", err)
	}

	err = slash.InstallSlashCommand(session, "", discordgo.ApplicationCommand{
		Name:        "clap",
		Description: "Applause!",
		Options:     []*discordgo.ApplicationCommandOption{},
	})

	if err != nil {
		return fmt.Errorf("error installing /clap %w", err)
	}

	return nil
}

func (g *GiphyCommands) slashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	subject := ""

	if len(i.Data.Options) > 0 {
		if r, ok := i.Data.Options[0].Value.(string); ok {
			subject = r
		}
	}

	if subject == "" {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: "Something went wrong",
				Flags:   64,
			},
		})

		if err != nil {
			log.Println(err)
		}

		return
	}

	g.postRandomGif(s, i, subject)
}

func (g *GiphyCommands) clap(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ch, _ := voice.FindVoiceUser(s, i.GuildID, i.Member.User.ID)

	if ch != "" {
		err := g.server.GetDiscordHA().SendVoiceCommand(discordha.VoiceCommand{
			ModuleID:  "thomasbot",
			GuildID:   i.GuildID,
			ChannelID: ch,
			File:      "clappingmono.wav",
			UserID:    i.Member.User.ID,
		})
		if err != nil {
			log.Printf("Error sending voice command: %q\n", err)
		}
	}
	g.postRandomGif(s, i, "applause")
}

func (g *GiphyCommands) postRandomGif(s *discordgo.Session, i *discordgo.InteractionCreate, subject string) {
	// TODO: share config with commands
	if os.Getenv("THOMASBOT_GIPHYKEY") == "" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: "Giphy key is lacking from deployment",
				Flags:   64,
			},
		})
		return
	}
	giphy := libgiphy.NewGiphy(os.Getenv("THOMASBOT_GIPHYKEY"))
	data, err := giphy.GetRandom(subject)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf("Error: %v", err),
				Flags:   64,
			},
		})
		return
	}

	embed := embed.NewEmbed()
	embed.SetImage(data.Data.Fixed_height_downsampled_url)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "",
			Embeds:  []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}

// Info return the commands in this package
func (g *GiphyCommands) Info() []command.Command {
	return []command.Command{}
}
