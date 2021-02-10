package hive

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/itfactory-tm/thomas-bot/pkg/db"

	"github.com/itfactory-tm/thomas-bot/pkg/embed"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

// HiveCommand contains the tm!hello command
type HiveCommand struct {
	isBob        bool
	requestRegex *regexp.Regexp
	db           db.Database
}

// NewHiveCommand gives a new HiveCommand
func NewHiveCommand(dbConn db.Database) *HiveCommand {
	return &HiveCommand{
		requestRegex: regexp.MustCompile(`!hive ([a-zA-Z0-9-_]*) ([a-zA-Z0-9]*) ?(.*)$`),
		db:           dbConn,
	}
}

// NewHiveCommand gives a new HiveCommand
func NewHiveCommandForBob(db db.Database) *HiveCommand {
	return &HiveCommand{
		isBob:        true,
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

// InstallSlashCommands registers the slash commands handlers
func (h *HiveCommand) InstallSlashCommands(session *discordgo.Session) error {
	_, err := session.ApplicationCommandCreate("", "", &discordgo.ApplicationCommand{
		Name:        "hive",
		Description: "creates on-remand voice and text channels",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
				Name:        "type",
				Description: "type of channel",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "text",
						Description: "text channel",
						Default:     false,
						Required:    false,
						Choices:     nil,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "name",
								Description: "name of channel",
								Required:    true,
							},
							{
								Type:        discordgo.ApplicationCommandOptionBoolean,
								Name:        "hidden",
								Description: "is channel not visible for everyone",
								Required:    true,
							},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionSubCommand,
						Name:        "voice",
						Description: "voice channel",
						Default:     false,
						Required:    false,
						Choices:     nil,
						Options: []*discordgo.ApplicationCommandOption{
							{
								Type:        discordgo.ApplicationCommandOptionString,
								Name:        "name",
								Description: "name of channel",
								Required:    true,
							},
							{
								Type:        discordgo.ApplicationCommandOptionInteger,
								Name:        "size",
								Description: "number of allowed users (1-99)",
								Required:    true,
							},
						},
					},
				},
			},
		},
	})

	return err
}

// SayHive handles the tm!hive command
func (h *HiveCommand) SayHive(s *discordgo.Session, m *discordgo.MessageCreate) {
	hidden := false

	conf, isHive, err := h.getConfigForRequestChannel(m.GuildID, m.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}
	if !isHive {
		if h.isBob {
			s.ChannelMessageSend(m.ChannelID, "This command only works in the bob-commands channel")
		} else {
			s.ChannelMessageSend(m.ChannelID, "This command only works in the Requests channels")
		}
		return
	}

	matched := h.requestRegex.FindStringSubmatch(m.Content)
	if len(matched) <= 2 {
		if h.isBob {
			s.ChannelMessageSend(m.ChannelID, "Incorrect syntax, syntax is `bob!vc channel-name <number of participants>`")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Incorrect syntax, syntax is `tm!hive channel-name <number of participants>`")
		}
		return
	}

	if matched[3] == "hidden" {
		hidden = true
	}

	var newChan *discordgo.Channel
	isText := false
	if matched[2] == "text" {
		newChan, err = h.createTextChannel(s, m, conf, matched[1], conf.TextCategoryID, hidden)
		isText = true
	} else {
		i, err := strconv.ParseInt(matched[2], 10, 64)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%q is not a number", matched[2]))
			return
		}
		newChan, err = h.createVoiceChannel(s, m, conf, matched[1], conf.VoiceCategoryID, int(i), hidden)
	}

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Channel created! Have fun! Reminder: I will delete it when it stays empty for a while")
	if isText && h.isBob {
		s.ChannelMessageSend(newChan.ID, "Welcome to your text channel! If you're finished using this please say `bob!archive`")
	} else if isText {
		s.ChannelMessageSend(newChan.ID, "Welcome to your text channel! If you're finished using this please say `tm!archive`")
	}

	if hidden {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("your channel is hidden, react 👋 below to join"))
		e := embed.NewEmbed()
		e.SetTitle("Hive Channel")
		e.AddField("name", conf.Prefix+matched[1])
		e.AddField("id", newChan.ID)

		msg, err := s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
		if err != nil {
			log.Println(err)
		}
		s.MessageReactionAdd(m.ChannelID, msg.ID, "👋")
	}
}

func (h *HiveCommand) createTextChannel(s *discordgo.Session, m *discordgo.MessageCreate, conf *db.HiveConfiguration, name, catID string, hidden bool) (*discordgo.Channel, error) {
	props := discordgo.GuildChannelCreateData{
		Name:     conf.Prefix + name,
		NSFW:     false,
		ParentID: catID,
		Type:     discordgo.ChannelTypeGuildText,
	}

	if hidden {
		// admin privileges on channel for creator
		var allow int64
		allow |= discordgo.PermissionReadMessageHistory
		allow |= discordgo.PermissionViewChannel
		allow |= discordgo.PermissionSendMessages
		allow |= discordgo.PermissionVoiceConnect
		allow |= discordgo.PermissionManageMessages
		props.PermissionOverwrites = []*discordgo.PermissionOverwrite{
			{
				ID:    m.Author.ID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Deny:  0,
				Allow: allow,
			},
		}
	}

	return s.GuildChannelCreateComplex(m.GuildID, props)
}

// we filled up on junk quickly, we should recycle a voice channel from junkjard
func (h *HiveCommand) recycleVoiceChannel(s *discordgo.Session, m *discordgo.MessageCreate, conf *db.HiveConfiguration, name, catID string, limit int, hidden bool) (*discordgo.Channel, error, bool) {
	channels, err := s.GuildChannels(m.GuildID)
	if err != nil {
		return nil, err, true
	}

	var toRecycle *discordgo.Channel

	for _, channel := range channels {
		if channel.ParentID == conf.JunkyardCategoryID && channel.Type == discordgo.ChannelTypeGuildVoice {
			toRecycle = channel
			break
		}
	}

	// no junk found we can recycle, buy a new one :(
	if toRecycle == nil {
		return nil, nil, false
	}

	cat, err := s.Channel(catID)
	if err != nil {
		return nil, err, true
	}

	edit := &discordgo.ChannelEdit{
		ParentID:             catID,
		PermissionOverwrites: cat.PermissionOverwrites,
		UserLimit:            limit,
		Name:                 conf.Prefix + name,
		Bitrate:              128000,
	}

	if hidden {
		// admin privileges on channel for creator
		var allow int64
		allow |= discordgo.PermissionReadMessageHistory
		allow |= discordgo.PermissionViewChannel
		allow |= discordgo.PermissionSendMessages
		allow |= discordgo.PermissionVoiceConnect
		allow |= discordgo.PermissionManageMessages
		edit.PermissionOverwrites = []*discordgo.PermissionOverwrite{
			{
				ID:    m.Author.ID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Deny:  0,
				Allow: allow,
			},
		}
	}
	newChan, err := s.ChannelEditComplex(toRecycle.ID, edit)

	return newChan, err, true
}

func (h *HiveCommand) createVoiceChannel(s *discordgo.Session, m *discordgo.MessageCreate, conf *db.HiveConfiguration, name, catID string, limit int, hidden bool) (*discordgo.Channel, error) {
	newChan, err, ok := h.recycleVoiceChannel(s, m, conf, name, catID, limit, hidden)
	if ok {
		return newChan, err
	}
	props := discordgo.GuildChannelCreateData{
		Name:      conf.Prefix + name,
		Bitrate:   conf.VoiceBitrate,
		NSFW:      false,
		ParentID:  catID,
		Type:      discordgo.ChannelTypeGuildVoice,
		UserLimit: limit,
	}

	if hidden {
		// admin privileges on channel for creator
		var allow int64
		allow |= discordgo.PermissionReadMessageHistory
		allow |= discordgo.PermissionViewChannel
		allow |= discordgo.PermissionSendMessages
		allow |= discordgo.PermissionVoiceConnect
		allow |= discordgo.PermissionManageMessages
		props.PermissionOverwrites = []*discordgo.PermissionOverwrite{
			{
				ID:    m.Author.ID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Deny:  0,
				Allow: allow,
			},
		}
	}

	return s.GuildChannelCreateComplex(m.GuildID, props)
}

// SayArchive handles the tm!archive command
func (h *HiveCommand) SayArchive(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}

	conf, isHive, err := h.getConfigForRequestCategory(s, m.GuildID, m.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}
	if !isHive || h.isPrivilegedChannel(channel.ID, conf) {
		s.ChannelMessageSend(m.ChannelID, "This command only works in hive created channels")
		return
	}

	if conf.Prefix != "" {
		if !strings.HasPrefix(channel.Name, conf.Prefix) {
			s.ChannelMessageSend(m.ChannelID, "This command only works in hive created channels with correct prefix")
			return
		}
	}

	j, err := s.Channel(conf.JunkyardCategoryID)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = s.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		ParentID:             conf.JunkyardCategoryID,
		PermissionOverwrites: j.PermissionOverwrites,
	})
	if err != nil {
		log.Println(err)
	}
}

// SayLeave handles the tm!eave command
func (h *HiveCommand) SayLeave(s *discordgo.Session, m *discordgo.MessageCreate) {
	// check if this is allowed
	conf, isHive, err := h.getConfigForRequestCategory(s, m.GuildID, m.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}
	if !isHive || h.isPrivilegedChannel(m.ChannelID, conf) {
		s.ChannelMessageSend(m.ChannelID, "This command only works in hive created channels, consider using Discord's mute instead\"")
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Println(err)
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

	channel, err := s.Channel(message.Embeds[0].Fields[len(message.Embeds[0].Fields)-1].Value)
	if err != nil {
		log.Println(err)
		return
	}

	_, isHive, err := h.getConfigForRequestChannel(r.GuildID, r.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}

	if !isHive {
		//s.ChannelMessageSend(r.ChannelID, "Sorry category not allowed, try privilege escalating otherwise!")
		return
	}

	var allow int64
	allow |= discordgo.PermissionReadMessageHistory
	allow |= discordgo.PermissionViewChannel
	allow |= discordgo.PermissionSendMessages
	allow |= discordgo.PermissionVoiceConnect

	err = s.ChannelPermissionSet(channel.ID, r.UserID, discordgo.PermissionOverwriteTypeMember, allow, 0)
	if err != nil {
		log.Println("Cannot set permissions", err)
		return
	}

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
