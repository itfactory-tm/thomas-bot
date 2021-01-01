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
		requestRegex: regexp.MustCompile(`!hive ([a-zA-Z0-9-_]*) (.*) (.*)`),
	}
}

// NewHiveCommand gives a new HiveCommand
func NewHiveCommandForBob() *HiveCommand {
	return &HiveCommand{
		isBob:        true,
		requestRegex: regexp.MustCompile(`!vc ([a-zA-Z0-9-_]*) (.*) (.*)`),
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
	registry.RegisterMessageCreateHandler("join", h.SayJoin)
	registry.RegisterMessageCreateHandler("leave", h.SayLeave)
}

// SayHive handles the tm!hive command
func (h *HiveCommand) SayHive(s *discordgo.Session, m *discordgo.MessageCreate) {
	hidden := false

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

	if matched[3] == "hidden" {
		hidden = true
	}

	props := discordgo.GuildChannelCreateData{
		Name:      matched[1],
		Bitrate:   128000,
		NSFW:      false,
		ParentID:  catID,
		Type:      chanType,
		UserLimit: int(i),
	}

	if hidden {
		// this is why it is Alpha silly
		j, _ := s.Channel("780775904082395136")
		props.PermissionOverwrites = j.PermissionOverwrites
	}

	newChan, err := s.GuildChannelCreateComplex(m.GuildID, props)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Channel created! Have fun! Reminder: I will delete it when it stays empty for a while")
	if chanType == discordgo.ChannelTypeGuildText && h.isBob {
		s.ChannelMessageSend(newChan.ID, "Welcome to your text channel! If you're finished using this please say `bob!archive`")
	} else if chanType == discordgo.ChannelTypeGuildText {
		s.ChannelMessageSend(newChan.ID, "Welcome to your text channel! If you're finished using this please say `tm!archive`")
	}

	if hidden {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("your channel is hidden (alpha!!!!!!!!!!) say `tm!join %s`", newChan.ID))
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

// SayJoin handles the tm!join ALPHA command
func (h *HiveCommand) SayJoin(s *discordgo.Session, m *discordgo.MessageCreate) {
	matched := regexp.MustCompile(`!join (.*)`).FindStringSubmatch(m.Message.Content)
	if len(matched) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Incorrect syntax, idk why")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "I'm sorry for the terrible UX it is alpha, hase i told you that yet?")

	channel, _ := s.Channel(matched[1])

	allowed := false
	for _, catID := range channelToCategory {
		if channel.ParentID == catID {
			allowed = true
			break
		}
	}

	if !allowed {
		s.ChannelMessageSend(m.ChannelID, "Sorry category not allowed, try privilege escalating otherwise!")
		return
	}

	// target type 1 is user, yes excellent library...
	var allow int
	allow |= discordgo.PermissionReadMessageHistory
	allow |= discordgo.PermissionViewChannel
	allow |= discordgo.PermissionSendMessages

	s.ChannelPermissionSet(channel.ID, m.Author.ID, "1", allow, 0)

	s.ChannelMessageSend(channel.ID, fmt.Sprintf("Welcome <@%s>", m.Author.ID))
}

// SayLeave handles the tm!eave command
func (h *HiveCommand) SayLeave(s *discordgo.Session, m *discordgo.MessageCreate) {
	// check if this is allowed
	channel, err := s.Channel(m.ChannelID)
	ok := false
	for _, category := range channelToCategory {
		if channel.ParentID == category {
			ok = true
		}
	}
	if !ok {
		s.ChannelMessageSend(m.ChannelID, "This command only works in hive created channels, consider using Discord's mute instead")
		return
	}

	info, err := s.Channel(channel.ID)
	if err != nil {
		log.Println(err)
		return
	}

	newPerms := []*discordgo.PermissionOverwrite{}
	for _, perm := range info.PermissionOverwrites {
		if perm.ID != m.Author.ID {
			newPerms = append(newPerms, perm)
		}
	}
	_, err = s.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		PermissionOverwrites: newPerms,
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
