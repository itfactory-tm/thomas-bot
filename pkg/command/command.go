package command

import "github.com/bwmarrin/discordgo"

type Command struct {
	Name        string
	Category    Category
	Description string
	Hidden      bool
	Handler     func(*discordgo.Session, *discordgo.MessageCreate)
}
