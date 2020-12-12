package hive

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

const junkyard = "780775904082395136"

var channelToCategory = map[string]string{
	"775453791801049119": "775436992136871957", // the hive
	"787346218304274483": "787345995105173524", // ITF Gaming
}

// HiveCommand contains the tm!hello command
type HiveCommand struct {
	isBob        bool
	requestRegex *regexp.Regexp
}

// NewHiveCommand gives a new HiveCommand
func NewHiveCommand() *HiveCommand {
	return &HiveCommand{
		requestRegex: regexp.MustCompile(`!hive ([a-zA-Z0-9-_]*) (.*)`),
	}
}

// NewHiveCommand gives a new HiveCommand
func NewHiveCommandForBob() *HiveCommand {
	return &HiveCommand{
		isBob:        true,
		requestRegex: regexp.MustCompile(`!vc ([a-zA-Z0-9-_]*) (.*)`),
	}
}

// Register registers the handlers
func (h *HiveCommand) Register(registry command.Registry, server command.Server) {
	if h.isBob {
		registry.RegisterMessageCreateHandler("vc", h.SayHive)
	} else {
		registry.RegisterMessageCreateHandler("hive", h.SayHive)
	}

	registry.RegisterMessageCreateHandler("archive", h.SayArchive)
}

// SayHive handles the tm!hive command
func (h *HiveCommand) SayHive(s *discordgo.Session, m *discordgo.MessageCreate) {

	// check of in the request channel to apply limits
	catID, ok := channelToCategory[m.ChannelID]
	if !ok {
		s.ChannelMessageSend(m.ChannelID, "This command only works in the Requests channels")
		return
	}

	matched := h.requestRegex.FindStringSubmatch(m.Content)
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

	newChan, err := s.GuildChannelCreateComplex(m.GuildID, discordgo.GuildChannelCreateData{
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
	if chanType == discordgo.ChannelTypeGuildText {
		s.ChannelMessageSend(newChan.ID, "Welcome to your text channel! If you're finished using this please say `tm!archive`")
	}
}

// SayArchive handles the tm!archive command
func (h *HiveCommand) SayArchive(s *discordgo.Session, m *discordgo.MessageCreate) {

	// check if this is allowed to be archived
	channel, err := s.Channel(m.ChannelID)
	ok := false
	for _, category := range channelToCategory {
		if channel.ParentID == category {
			ok = true
		}
	}
	if !ok {
		s.ChannelMessageSend(m.ChannelID, "This command only works in hive created channels")
		return
	}

	j, err := s.Channel(junkyard)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = s.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		ParentID:             junkyard,
		PermissionOverwrites: j.PermissionOverwrites,
	})
	if err != nil {
		log.Println(err)
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
