package hello

import (
	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

// HelloCommand contains the tm!hello command
type HelloCommand struct{}

// NewHelloCommand gives a new HelloCommand
func NewHelloCommand() *HelloCommand {
	return &HelloCommand{}
}

func (h *HelloCommand) Register(registry command.Registry, server command.Server) {
	registry.RegisterMessageCreateHandler("hello", h.SayHello)
}

func (h *HelloCommand) SayHello(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Beep bop boop! I am Thomas Bot, fork me on GitHub!")
}

func (h *HelloCommand) Info() []command.Command {
	return []command.Command{
		command.Command{
			Name:        "hello",
			Category:    command.CategoryFun,
			Description: "Say hello world",
			Hidden:      false,
		},
	}
}
