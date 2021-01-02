package hive

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/itfactory-tm/thomas-bot/pkg/embed"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

const junkyard = "780775904082395136"

var channelToCategory = map[string]string{
	"775453791801049119": "775436992136871957", // the hive
	"787346218304274483": "760860082241142790", // ITF Gaming
}

// cats with prefixes
var categoryPrefixes = map[string]string{
	"775436992136871957": "",     // the hive
	"760860082241142790": "BOB-", // ITF gaming
}

// HiveCommand contains the tm!hello command
type HiveCommand struct {
	isBob        bool
	prefix       string
	requestRegex *regexp.Regexp
}

// NewHiveCommand gives a new HiveCommand
func NewHiveCommand() *HiveCommand {
	return &HiveCommand{
		requestRegex: regexp.MustCompile(`!hive ([a-zA-Z0-9-_]*) ([a-zA-Z0-9]*) ?(.*)$`),
	}
}

// NewHiveCommand gives a new HiveCommand
func NewHiveCommandForBob() *HiveCommand {
	return &HiveCommand{
		isBob:        true,
		prefix:       "BOB-",
		requestRegex: regexp.MustCompile(`!vc ([a-zA-Z0-9-_]*) ([a-zA-Z0-9]*) ?(.*)$`),
	}
}

// Register registers the handlers
func (h *HiveCommand) Register(registry command.Registry, server command.Server) {
	if h.isBob {
		registry.RegisterMessageCreateHandler("vc", h.SayHive)
	} else {
		registry.RegisterMessageCreateHandler("hive", h.SayHive)
		registry.RegisterMessageCreateHandler("attendance", h.SayAttendance)
		registry.RegisterMessageCreateHandler("verify", h.SayVerify)
	}

	registry.RegisterMessageReactionAddHandler(h.handleReaction)
	registry.RegisterMessageCreateHandler("archive", h.SayArchive)
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
		Name:      h.prefix + matched[1],
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
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("your channel is hidden (alpha!!!!!!!!!!), react 👋 below to join"))
		e := embed.NewEmbed()
		e.SetTitle("Hive Channel")
		e.AddField("name", h.prefix+matched[1])
		e.AddField("id", newChan.ID)

		msg, err := s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
		if err != nil {
			log.Println(err)
		}
		s.MessageReactionAdd(m.ChannelID, msg.ID, "👋")
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

	if categoryPrefixes[channel.ParentID] != "" {
		if !strings.HasPrefix(channel.Name, categoryPrefixes[channel.ParentID]) {
			s.ChannelMessageSend(m.ChannelID, "This command only works in hive created channels with correct prefix")
			return
		}
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

	newPerms := []*discordgo.PermissionOverwrite{}
	for _, perm := range channel.PermissionOverwrites {
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

func (h *HiveCommand) handleReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Println("Cannot get message of reaction", r.ChannelID)
		return
	}

	if message.Author.ID != s.State.User.ID {
		return // not the bot user
	}

	if len(message.Embeds) < 1 {
		return // not an embed
	}

	if len(message.Embeds[0].Fields) < 2 {
		return // not the correct embed
	}

	if message.Embeds[0].Title != "Hive Channel" {
		return // not the hive message
	}

	channel, err := s.Channel(message.Embeds[0].Fields[1].Value)
	if err != nil {
		log.Println(err)
	}

	allowed := false
	for _, catID := range channelToCategory {
		if channel.ParentID == catID {
			allowed = true
			break
		}
	}

	if !allowed {
		//s.ChannelMessageSend(r.ChannelID, "Sorry category not allowed, try privilege escalating otherwise!")
		return
	}

	// target type 1 is user, yes excellent library...
	var allow int
	allow |= discordgo.PermissionReadMessageHistory
	allow |= discordgo.PermissionViewChannel
	allow |= discordgo.PermissionSendMessages
	allow |= discordgo.PermissionVoiceConnect

	s.ChannelPermissionSet(channel.ID, r.UserID, "1", allow, 0)

	s.ChannelMessageSend(channel.ID, fmt.Sprintf("Welcome <@%s>, you can leave any time by saying `tm!leave`", r.UserID))
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
