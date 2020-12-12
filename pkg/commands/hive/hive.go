package hive

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

var channelToCategory = map[string]string{
	"775453791801049119": "775436992136871957", // the hive
}

var requestRegex = regexp.MustCompile(`!hive ([a-zA-Z0-9-_]*) (.*)`)

// HiveCommand contains the tm!hello command
type HiveCommand struct{}

// NewHiveCommand gives a new HiveCommand
func NewHiveCommand() *HiveCommand {
	return &HiveCommand{}
}

// Register registers the handlers
func (h *HiveCommand) Register(registry command.Registry, server command.Server) {
	registry.RegisterMessageCreateHandler("hive", h.SayHive)
}

// SayHive handles the tm!hive command
func (h *HiveCommand) SayHive(s *discordgo.Session, m *discordgo.MessageCreate) {

	// check of in the request channel to apply limits
	catID, ok := channelToCategory[m.ChannelID]
	if !ok {
		s.ChannelMessageSend(m.ChannelID, "This command only works in the Requests channels")
		return
	}

	matched := requestRegex.FindStringSubmatch(m.Content)
	if len(matched) <= 2 {
		s.ChannelMessageSend(m.ChannelID, "Incorrect syntax, syntax is `tm!hive channel-name <number of participants>`")
		return
	}

	var i int64
	chanType := discordgo.ChannelTypeGuildVoice
	if matched[2] == "text" {
		chanType = discordgo.ChannelTypeGuildText
	} else {
		var err error
		i, err = strconv.ParseInt(matched[2], 10, 64)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%q is not a number", matched[2]))
			return
		}
	}

	_, err := s.GuildChannelCreateComplex(m.GuildID, discordgo.GuildChannelCreateData{
		Name:      matched[1],
		Bitrate:   128000,
		NSFW:      false,
		ParentID:  catID,
		Type:      chanType,
		UserLimit: int(i),
	})

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Channel created! Have fun! Reminder: I will delete it when it stays empty for a while")
}

// Info return the commands in this package
func (h *HiveCommand) Info() []command.Command {
	return []command.Command{
		command.Command{
			Name:        "hive",
			Category:    command.CategoryFun,
			Description: "Set up temporary meeting rooms",
			Hidden:      false,
		},
	}
}
