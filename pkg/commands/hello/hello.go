package hello

import (
	"log"
	"reflect"

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
	app := discordgo.ApplicationCommand{
		Name:        "hello",
		Description: "Thomas Bot will say hello",
		Options:     []*discordgo.ApplicationCommandOption{},
	}

	cmds, err := session.ApplicationCommands(session.State.User.ID, "") // ITF only for now till links are moved to a DB
	if err != nil {
		return err
	}
	exists := false
	for _, cmd := range cmds {
		if cmd.Name == "links" {
			exists = reflect.DeepEqual(app.Options, cmd.Options)
		}
	}

	if !exists {
		_, err = session.ApplicationCommandCreate(session.State.User.ID, "", &app)
	}

	return err
}

// SayHello sends an hello message
func (h *HelloCommand) SayHello(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
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
