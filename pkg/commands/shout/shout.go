package shout

import (
	"fmt"
	"log"

	"github.com/itfactory-tm/thomas-bot/pkg/util/slash"

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
	registry.RegisterInteractionCreate("shout", s.shout)

	s.server = server
}

// InstallSlashCommands registers the slash commands
func (s *ShoutCommand) InstallSlashCommands(session *discordgo.Session) error {
	return slash.InstallSlashCommand(session, "", discordgo.ApplicationCommand{
		Name:        "shout",
		Description: "Gives a quote in your audio channel",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "number",
				Description: "number of the audio clip",
				Required:    true,
			},
		},
	})
}

func (s *ShoutCommand) shout(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	ch, err := voice.FindVoiceUser(sess, i.GuildID, i.Member.User.ID)
	if err != nil {
		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf("An error happened: %s", err),
				Flags:   64, // ephemeral
			},
		})
		log.Println(err)
		return
	}

	err = s.server.GetDiscordHA().SendVoiceCommand(discordha.VoiceCommand{
		ModuleID:  "thomasbot",
		GuildID:   i.GuildID,
		ChannelID: ch,
		File:      fmt.Sprintf("%d.wav", int(i.Data.Options[0].Value.(float64))),
		UserID:    i.Member.User.ID,
	})
	if err != nil {
		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: fmt.Sprintf("Error sending voice command: %q", err),
				Flags:   64, // ephemeral
			},
		})
	}

	sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "Go Go Go",
		},
	})

}

// Info return the commands in this package
func (s *ShoutCommand) Info() []command.Command {
	return []command.Command{}
}
