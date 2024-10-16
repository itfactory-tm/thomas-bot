package hive

import (
	"fmt"
	"log"
	"strings"

	"github.com/itfactory-tm/thomas-bot/pkg/util/slash"

	"github.com/itfactory-tm/thomas-bot/pkg/db"

	"github.com/itfactory-tm/thomas-bot/pkg/embed"

	"github.com/bwmarrin/discordgo"
)

var defaultAllows int64

func init() {
	defaultAllows |= discordgo.PermissionReadMessageHistory
	defaultAllows |= discordgo.PermissionViewChannel
	defaultAllows |= discordgo.PermissionSendMessages
	defaultAllows |= discordgo.PermissionVoiceConnect
	defaultAllows |= discordgo.PermissionAddReactions
	defaultAllows |= discordgo.PermissionAttachFiles
	defaultAllows |= discordgo.PermissionEmbedLinks
}

// HiveCommand contains the handlers for the Hive commandos
type HiveCommand struct {
	db    db.Database
	isBob bool // deprecated
}

// NewHiveCommand gives a new HiveCommand
func NewHiveCommand(dbConn db.Database) *HiveCommand {
	return &HiveCommand{
		db: dbConn,
	}
}

// deprecated
func NewHiveCommandForBob(dbConn db.Database) *HiveCommand {
	return &HiveCommand{
		isBob: true,
		db:    dbConn,
	}
}

// InstallSlashCommands registers the slash commands handlers
func (h *HiveCommand) InstallSlashCommands(session *discordgo.Session) error {
	if session == nil {
		return nil
	}

	if err := h.installArchive(session); err != nil {
		return err
	}

	if err := h.installLeave(session); err != nil {
		return err
	}

	if err := h.installHive(session); err != nil {
		return err
	}

	return nil
}

func (h *HiveCommand) installLeave(session *discordgo.Session) error {
	return slash.InstallSlashCommand(session, "", discordgo.ApplicationCommand{
		Name:        "leave",
		Description: "Leave an on-remand text channel",
		Options:     []*discordgo.ApplicationCommandOption{},
	})
}

func (h *HiveCommand) installArchive(session *discordgo.Session) error {
	return slash.InstallSlashCommand(session, "", discordgo.ApplicationCommand{
		Name:        "archive",
		Description: "Archives an on-remand text channel",
		Options:     []*discordgo.ApplicationCommandOption{},
	})
}

func (h *HiveCommand) installHive(session *discordgo.Session) error {
	return slash.InstallSlashCommand(session, "", discordgo.ApplicationCommand{
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
}

func (h *HiveCommand) HiveCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return // not the app command
	}

	hidden := false
	istext := false
	var size int
	var name string
	if len(i.ApplicationCommandData().Options) < 1 || len(i.ApplicationCommandData().Options[0].Options) < 1 || len(i.ApplicationCommandData().Options[0].Options[0].Options) < 2 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Invalid command options",
				Flags:   64,
			},
		})
		return // invalid
	}

	if i.ApplicationCommandData().Options[0].Options[0].Name == "text" {
		istext = true
	}

	conf, ok := h.precheck(s, i)
	if !ok {
		return
	}

	for _, option := range i.ApplicationCommandData().Options[0].Options[0].Options {
		switch option.Name {
		case "name":
			name, ok = option.Value.(string)
			if !ok {
				return
			}
		case "size":
			fVal, ok := option.Value.(float64)
			if !ok {
				return
			}
			size = int(fVal)
		case "hidden":
			hidden, ok = option.Value.(bool)
			if !ok {
				return
			}
		}
	}

	h.createChannel(s, i, name, istext, hidden, conf, size)
}

func (h *HiveCommand) precheck(s *discordgo.Session, i *discordgo.InteractionCreate) (*db.HiveConfiguration, bool) {
	if i.Member == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command does not work in DMs",
				Flags:   64,
			},
		})
		return nil, false
	}

	conf, isHive, err := h.getConfigForRequestChannel(i.GuildID, i.ChannelID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("An error happened: %v", err),
				Flags:   64,
			},
		})
		return nil, false
	}
	if !isHive {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command only works in the Requests channels",
				Flags:   64,
			},
		})
		return nil, false
	}

	return conf, true
}

func (h *HiveCommand) createChannel(s *discordgo.Session, i *discordgo.InteractionCreate, name string, isText, hidden bool, conf *db.HiveConfiguration, size int) {
	var newChan *discordgo.Channel
	var err error
	if isText {
		newChan, err = h.createTextChannel(s, conf, name, conf.TextCategoryID, i, hidden)
	} else {
		newChan, err = h.CreateVoiceChannel(s, conf, name, conf.VoiceCategoryID, i.GuildID, size, hidden)
	}

	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("An error happened: %v", err),
				Flags:   64,
			},
		})
		return
	}

	if !isText {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Channel <#%s> has been created!  Reminder: I will delete it when it stays empty for a while", newChan.ID),
				Flags:   64,
			},
		})
	} else if !hidden {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Channel <#%s> has been created!", newChan.ID),
				Flags:   64,
			},
		})
		s.ChannelMessageSend(newChan.ID, "Welcome to your text channel! If you're finished using this please say `/archive`")
	} else {
		e := embed.NewEmbed()
		e.SetTitle("Hive Channel")
		e.AddField("name", conf.Prefix+name)
		e.AddField("id", newChan.ID)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Channel has been created! Your channel is hidden, click 👋 below to join",
				Embeds:  []*discordgo.MessageEmbed{e.MessageEmbed},
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label:    "Join Channel",
								Style:    discordgo.SuccessButton,
								CustomID: "hive_join",
								Emoji: &discordgo.ComponentEmoji{
									Name: "👋",
								},
							},
						},
					},
				},
			},
		})
	}
}

func (h *HiveCommand) createTextChannel(s *discordgo.Session, conf *db.HiveConfiguration, name, catID string, i *discordgo.InteractionCreate, hidden bool) (*discordgo.Channel, error) {
	cat, err := s.Channel(catID)
	if err != nil {
		return nil, err
	}
	props := discordgo.GuildChannelCreateData{
		Name:                 conf.Prefix + name,
		Type:                 discordgo.ChannelTypeGuildText,
		Position:             99,
		ParentID:             catID,
		NSFW:                 false,
		PermissionOverwrites: cat.PermissionOverwrites,
	}

	if hidden {
		// admin privileges on channel for creator
		var allow = defaultAllows
		allow |= discordgo.PermissionManageMessages
		props.PermissionOverwrites = []*discordgo.PermissionOverwrite{
			{
				ID:    i.Member.User.ID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Deny:  0,
				Allow: allow,
			},
			{
				ID:   i.GuildID,
				Type: discordgo.PermissionOverwriteTypeRole,
				Deny: discordgo.PermissionAll,
			},
		}
	}

	return s.GuildChannelCreateComplex(i.GuildID, props)
}

// we filled up on junk quickly, we should recycle a voice channel from junkjard
func (h *HiveCommand) recycleVoiceChannel(s *discordgo.Session, conf *db.HiveConfiguration, name, catID, guildID string, limit int, hidden bool) (*discordgo.Channel, error, bool) {
	channels, err := s.GuildChannels(guildID)
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
		Bitrate:              conf.VoiceBitrate,
	}

	newChan, err := s.ChannelEditComplex(toRecycle.ID, edit)

	return newChan, err, true
}

func (h *HiveCommand) CreateVoiceChannel(s *discordgo.Session, conf *db.HiveConfiguration, name, catID, guildID string, limit int, hidden bool) (*discordgo.Channel, error) {
	newChan, err, ok := h.recycleVoiceChannel(s, conf, name, catID, guildID, limit, hidden)
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

	return s.GuildChannelCreateComplex(guildID, props)
}

// SayArchive handles the archive command
func (h *HiveCommand) SayArchive(s *discordgo.Session, i *discordgo.InteractionCreate) {
	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error getting channel info",
				Flags:   64,
			},
		})
		log.Println(err)
		return
	}

	conf, isHive, err := h.getConfigForRequestCategory(s, i.GuildID, i.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}
	if !isHive || h.isPrivilegedChannel(channel.ID, conf) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command only works in hive created channels",
				Flags:   64,
			},
		})
		return
	}

	if conf.Prefix != "" {
		if !strings.HasPrefix(channel.Name, conf.Prefix) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This command only works in hive created channels with correct prefix",
					Flags:   64,
				},
			})

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

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Channel is archived",
		},
	})
}

// SayLeave handles the tm!eave command
func (h *HiveCommand) SayLeave(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// check if this is allowed
	conf, isHive, err := h.getConfigForRequestCategory(s, i.GuildID, i.ChannelID)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
				Flags:   64,
			},
		})
		log.Println(err)
		return
	}
	if !isHive || h.isPrivilegedChannel(i.ChannelID, conf) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command only works in hive created channels, consider using Discord's mute instead",
				Flags:   64,
			},
		})

		return
	}

	channel, err := s.Channel(i.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}

	newPerms := []*discordgo.PermissionOverwrite{}
	for _, perm := range channel.PermissionOverwrites {
		if perm.ID != i.Member.User.ID {
			newPerms = append(newPerms, perm)
		}
	}
	_, err = s.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
		PermissionOverwrites: newPerms,
	})
	if err != nil {
		log.Println(err)
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("<@%s> has left the chat", i.Member.User.ID),
		},
	})
}
