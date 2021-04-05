package command

import (
	"github.com/bwmarrin/discordgo"
	discordha "github.com/meyskens/discord-ha"
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

// Registry is the interface of a command registry
type Registry interface {
	// if command is "" all messages will be sent
	RegisterMessageCreateHandler(command string, fn func(*discordgo.Session, *discordgo.MessageCreate))
	RegisterMessageEditHandler(command string, fn func(*discordgo.Session, *discordgo.MessageUpdate))
	RegisterMessageReactionAddHandler(fn func(*discordgo.Session, *discordgo.MessageReactionAdd))
	RegisterGuildMemberAddHandler(fn func(*discordgo.Session, *discordgo.GuildMemberAdd))
	RegisterMessageReactionRemoveHandler(fn func(*discordgo.Session, *discordgo.MessageReactionRemove))
	RegisterInteractionCreate(command string, fn func(*discordgo.Session, *discordgo.InteractionCreate))
}

// Interface defines how a command should be structured
type Interface interface {
	Info() []Command
	Register(registry Registry, server Server)
	InstallSlashCommands(session *discordgo.Session) error
}

// Server represents a discord bot server
type Server interface {
	GetDiscordHA() discordha.HA
	GetAllCommandInfos() []Command
}
