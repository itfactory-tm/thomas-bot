package command

import "github.com/bwmarrin/discordgo"

// Command is a struct of a bot command
type Command struct {
	Name        string
	Category    Category
	Description string
	Hidden      bool
	Handler     func(*discordgo.Session, *discordgo.MessageCreate)
}
