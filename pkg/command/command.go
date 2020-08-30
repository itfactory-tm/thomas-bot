package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/discordha"
)

// Command is a struct of a bot command
type Command struct {
	Name        string
	Category    Category
	Description string
	Hidden      bool
	// deprecated
	Handler func(*discordgo.Session, *discordgo.MessageCreate)
}

type Registry interface {
	// if command is "" all messages will be sent
	RegisterMessageCreateHandler(command string, fn func(*discordgo.Session, *discordgo.MessageCreate))
	RegisterMessageEditHandler(command string, fn func(*discordgo.Session, *discordgo.MessageUpdate))
	RegisterMessageReactionAddHandler(fn func(*discordgo.Session, *discordgo.MessageReactionAdd))
	RegisterGuildMemberAddHandler(fn func(*discordgo.Session, *discordgo.GuildMemberAdd))
}

type Interface interface {
	Info() []Command
	Register(registry Registry, server Server)
}

type Server interface {
	GetDiscordHA() *discordha.HA
	GetAllCommandInfos() []Command
}
