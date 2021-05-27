package hello

import (
	"log"

	"github.com/itfactory-tm/thomas-bot/pkg/util/slash"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

// HelloCommand contains the tm!hello command
type HelloCommand struct{}

// NewHelloCommand gives a new HelloCommand
func NewHelloCommand() *HelloCommand {
	return &HelloCommand{}
}

// Register registers the handlers
func (h *HelloCommand) Register(registry command.Registry, server command.Server) {
	registry.RegisterInteractionCreate("hello", h.SayHello)
}

// InstallSlashCommands registers the slash commands
func (h *HelloCommand) InstallSlashCommands(session *discordgo.Session) error {
	return slash.InstallSlashCommand(session, "", discordgo.ApplicationCommand{
		Name:        "hello",
		Description: "Thomas Bot will say hello",
		Options:     []*discordgo.ApplicationCommandOption{},
	})
}

// SayHello sends an hello message
func (h *HelloCommand) SayHello(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Beep bop boop! I am Thomas Bot, fork me on GitHub!",
		},
	})

	if err != nil {
		log.Println(err)
	}
}

// Info return the commands in this package
func (h *HelloCommand) Info() []command.Command {
	return []command.Command{}
}
