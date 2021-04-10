package hive

import (
	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

// Register registers the handlers
func (h *HiveCommand) Register(registry command.Registry, server command.Server) {
	if h.isBob {
		registry.RegisterMessageCreateHandler("vc", func(session *discordgo.Session, create *discordgo.MessageCreate) {
			session.ChannelMessageSend(create.ChannelID, "Thank you for your interest in `bob!vc`, due to new Discord capabilities this command has been replaced by `/hive`. We hope you enjoy the new experience!")
		})
	} else {
		registry.RegisterMessageCreateHandler("hive", func(session *discordgo.Session, create *discordgo.MessageCreate) {
			session.ChannelMessageSend(create.ChannelID, "Thank you for your interest in `tm!hive`, due to new Discord capabilities this command has been replaced by `/hive`. We hope you enjoy the new experience!")
		})
		registry.RegisterMessageCreateHandler("attendance", h.SayAttendance)
		registry.RegisterMessageCreateHandler("verify", h.SayVerify)
	}

	registry.RegisterMessageReactionAddHandler(h.handleReaction)
	registry.RegisterMessageCreateHandler("archive", h.SayArchive)
	registry.RegisterMessageCreateHandler("leave", h.SayLeave)

	if !h.isBob {
		registry.RegisterInteractionCreate("hive", h.HiveCommand)
	}
}

// Info return the commands in this package
func (h *HiveCommand) Info() []command.Command {
	if h.isBob {
		return []command.Command{
			command.Command{
				Name:        "vc",
				Category:    command.CategoryFun,
				Description: "Set up temporary gaming rooms",
				Hidden:      false,
			},
			command.Command{
				Name:        "archive",
				Category:    command.CategoryFun,
				Description: "Archive temporary text gaming rooms",
				Hidden:      false,
			},
		}
	}

	return []command.Command{
		command.Command{
			Name:        "hive",
			Category:    command.CategoryFun,
			Description: "Set up temporary meeting rooms",
			Hidden:      false,
		},
		command.Command{
			Name:        "archive",
			Category:    command.CategoryFun,
			Description: "Archive temporary text meeting rooms",
			Hidden:      false,
		},
	}
}
