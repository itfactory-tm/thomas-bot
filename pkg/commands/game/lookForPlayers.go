package game

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/util/slash"

	"github.com/itfactory-tm/thomas-bot/pkg/commands/hive"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
	"github.com/itfactory-tm/thomas-bot/pkg/db"
)

//GuildID for init of slash commands
const tmGaming = "773847927910432789"
const itf = "687565213943332875"

// LookCommand contains the /lookforplayers command
type LookCommand struct {
	db db.Database
}

// NewLookCommand gives a new LookCommand
func NewLookCommand(dbConn db.Database) *LookCommand {
	return &LookCommand{
		db: dbConn,
	}
}

// Register registers the handlers
func (l *LookCommand) Register(registry command.Registry, server command.Server) {
	registry.RegisterInteractionCreate("lookforplayers", l.SearchCommand)
	registry.RegisterMessageReactionAddHandler(l.handleReactionAdd)
	registry.RegisterMessageReactionRemoveHandler(l.handleReactionRemove)
}

// InstallSlashCommands registers the slash commands
// TODO: Make configurable for specific guilds
func (l *LookCommand) InstallSlashCommands(s *discordgo.Session) error {
	app := discordgo.ApplicationCommand{
		Name:        "lookforplayers",
		Description: "Send out an invitation to look for players!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "game",
				Description: "Name of the game",
				Required:    true,
			}, {
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "amount",
				Description: "Amount of people you need for the game",
				Required:    true,
			}, {
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "time",
				Description: "At what time do you want to play? Format hh:mm (example: 15:45)",
				Required:    false,
			}, {
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "notifyrole",
				Description: "Notify a role with your invitation!",
				Required:    false,
			},
		},
	}

	if err := slash.InstallSlashCommand(s, tmGaming, app); err != nil {
		return fmt.Errorf("error installing lfp in TM Gaming: %w", err)
	}

	if err := slash.InstallSlashCommand(s, itf, app); err != nil {
		return fmt.Errorf("error installing lfp in ITF: %w", err)
	}

	return nil
}

func (l *LookCommand) SearchCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	conf, isLFP, err := l.checkConfig(i.GuildID, i.ChannelID)
	if err != nil {
		l.sendInvisibleInteractionResponse(s, i, err.Error())
		return
	}
	if !isLFP {
		l.sendInvisibleInteractionResponse(s, i, "This command only works in Requests channels")
		return // not from a guild
	}

	inviteChannelID := conf.AdvertiseChannelID
	var name, selectedRoleID string
	var amount float64
	timeString := "Now!"
	var ok bool

	for _, option := range i.Data.Options {
		switch option.Name {
		case "game":
			name, ok = option.Value.(string)
			if !ok {
				l.sendInvisibleInteractionResponse(s, i, "Please enter a valid name.")
				return
			}
			if len(name) < 2 || len(name) > 25 {
				l.sendInvisibleInteractionResponse(s, i, "Your game needs to be between 2-25 characters long")
				return
			}
			if matched, _ := regexp.MatchString(`^[A-Za-z0-9 ]+$`, name); !matched {
				l.sendInvisibleInteractionResponse(s, i, "Your game cannot contain any special characters")
				return
			}

		case "amount":
			amount, ok = option.Value.(float64)
			if !ok {
				l.sendInvisibleInteractionResponse(s, i, "Please enter a valid amount.")
				return
			}
			if amount < 2 || amount > 40 {
				l.sendInvisibleInteractionResponse(s, i, "Your game needs to contain between 2-40 players")
				return
			}

		case "notifyrole":
			selectedRoleID, ok = option.Value.(string)
			if !ok {
				l.sendInvisibleInteractionResponse(s, i, "Please enter a valid role.")
				return
			}
			roles, _ := s.GuildRoles(i.GuildID)
			for _, role := range roles {
				if selectedRoleID == role.ID && role.Color != 0x9c9c9c {
					l.sendInvisibleInteractionResponse(s, i, "Please enter a valid gaming role.")
					return
				}
			}
			//Time case could be extended to handle different time zones for erasmus students
		case "time":
			timeString, ok = option.Value.(string)
			if !ok {
				l.sendInvisibleInteractionResponse(s, i, "Please enter a valid time in format 15:45.")
				return
			}
			if _, err := time.Parse("15:04", timeString); err != nil {
				l.sendInvisibleInteractionResponse(s, i, "Please enter your time in format hh:mm (For example 15:50)")
				return
			}
		}
	}

	err = l.createInviteEmbed(s, i, name, int(amount), timeString, selectedRoleID, inviteChannelID)
	content := fmt.Sprintf("Invite created in <#%v>!", inviteChannelID)
	if err != nil {
		content = err.Error()
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: content,
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func (l *LookCommand) checkConfig(guildID, channelID string) (*db.LookingForPlayersConfiguration, bool, error) {
	conf, err := l.db.ConfigForGuild(guildID)
	if err != nil {
		return nil, false, err
	}
	if conf == nil {
		// not in our DB
		return nil, false, nil
	}

	for _, lfp := range conf.LookingForPlayers {
		if channelID == lfp.AdvertiseChannelID {
			return &lfp, true, nil
		}
		for _, reqID := range lfp.RequestChannelIDs {
			if channelID == reqID {
				return &lfp, true, nil
			}
		}
	}

	// no lfp found
	return nil, false, nil
}

func (l *LookCommand) checkHiveConfig(guildID, channelID string) (*db.HiveConfiguration, bool, error) {
	conf, err := l.db.ConfigForGuild(guildID)
	if err != nil {
		return nil, false, err
	}
	if conf == nil {
		// not in our DB
		return nil, false, nil
	}

	for _, hive := range conf.Hives {
		for _, reqID := range hive.RequestChannelIDs {
			if channelID == reqID {
				return &hive, true, nil
			}
		}
	}

	// no hive found
	return nil, false, nil
}

func (l *LookCommand) createInviteEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, gameName string, amount int, timeString string, roleID string, inviteChannelID string) error {
	embed := &discordgo.MessageEmbed{
		Title: gameName,
		Color: 0x33FF33,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%v is looking for players!", i.Member.User.Username),
			IconURL: i.Member.User.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Host",
				Value:  i.Member.User.Mention(),
				Inline: true,
			}, {
				Name:   "Players needed",
				Value:  strconv.Itoa(amount),
				Inline: true,
			}, {
				Name:   "Playing at",
				Value:  timeString,
				Inline: true,
			}, {
				Name:   "Joined players",
				Value:  i.Member.User.Mention(),
				Inline: true,
			}, {
				Name:   "Backup players",
				Value:  "\u200b",
				Inline: true,
			}, {
				Name:   "\u200b",
				Value:  "\u200b",
				Inline: true,
			}, {
				Name:   "Join active",
				Value:  "👋",
				Inline: true,
			}, {
				Name:   "Join backup",
				Value:  "💾",
				Inline: true,
			}, {
				Name:   "\u200b",
				Value:  "\u200b",
				Inline: true,
			}, {
				Name:   "Delete Invite",
				Value:  "🗑️",
				Inline: true,
			}, {
				Name:   "Start game",
				Value:  "🎮",
				Inline: true,
			}, {
				Name:   "\u200b",
				Value:  "\u200b",
				Inline: true,
			},
		},
	}

	var sentMessage *discordgo.Message
	var err error

	if roleID != "" {
		message := &discordgo.MessageSend{
			Content: fmt.Sprintf("<@&%s>", roleID),
			Embed:   embed,
		}
		sentMessage, err = s.ChannelMessageSendComplex(inviteChannelID, message)
	} else {
		sentMessage, err = s.ChannelMessageSendEmbed(inviteChannelID, embed)
	}
	if err != nil {
		return errors.New(fmt.Sprintf("Error sending embed message: %v", err))
	}
	s.MessageReactionAdd(sentMessage.ChannelID, sentMessage.ID, "👋")
	s.MessageReactionAdd(sentMessage.ChannelID, sentMessage.ID, "💾")
	s.MessageReactionAdd(sentMessage.ChannelID, sentMessage.ID, "🗑️")
	s.MessageReactionAdd(sentMessage.ChannelID, sentMessage.ID, "🎮")
	return nil
}

func (l *LookCommand) sendInvisibleInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: content,
			Flags:   64,
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func (l *LookCommand) handleReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		return
	}
	if !l.checkEmbed(s, message) {
		return
	}

	if r.Emoji.MessageFormat() == "👋" {
		_, currentPlayers, backupPlayers, neededPlayers := l.getPlayers(s, message, r.UserID, true, true)
		l.handleJoinReaction(currentPlayers, backupPlayers, message, s)
		if message.Embeds[0].Fields[2].Value == "Now!" && len(currentPlayers) >= neededPlayers {
			l.startGame(s, r, currentPlayers, backupPlayers, neededPlayers, message)
		}
	}

	hostID, currentPlayers, backupPlayers, neededPlayers := l.getPlayers(s, message, r.UserID, true, false)

	if r.Emoji.MessageFormat() == "💾" {
		l.handleJoinReaction(currentPlayers, backupPlayers, message, s)
	}

	if r.UserID == hostID {
		if r.Emoji.MessageFormat() == "🗑️" {
			//Delete message first to prevent players being notified multiple times when emoji spam (Dirk proofing)
			err = s.ChannelMessageDelete(r.ChannelID, r.MessageID)
			if err != nil {
				return
			}
			//Notify players
			err = l.messagePlayers(s, currentPlayers, message.Embeds[0], fmt.Sprintf("The invite for %s has been deleted by the host.", message.Embeds[0].Title))
			if err != nil {
				return
			}
		}

		if r.Emoji.MessageFormat() == "🎮" {
			l.startGame(s, r, currentPlayers, backupPlayers, neededPlayers, message)
		}
	}
}

func (l *LookCommand) handleReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		return
	}
	if !l.checkEmbed(s, message) {
		return
	}

	if r.Emoji.MessageFormat() == "👋" {
		_, currentPlayers, backupPlayers, _ := l.getPlayers(s, message, r.UserID, false, true)
		l.handleJoinReaction(currentPlayers, backupPlayers, message, s)
	}

	if r.Emoji.MessageFormat() == "💾" {
		_, currentPlayers, backupPlayers, _ := l.getPlayers(s, message, r.UserID, false, false)
		l.handleJoinReaction(currentPlayers, backupPlayers, message, s)
	}
}

func (l *LookCommand) checkEmbed(s *discordgo.Session, message *discordgo.Message) bool {
	if message.Author.ID != s.State.User.ID {
		return false // not the bot user
	}
	if len(message.Embeds) < 1 {
		return false // not an embed
	}
	if len(message.Embeds[0].Fields) <= 5 {
		return false // not the correct embed
	}

	if message.Embeds[0].Fields[0].Name != "Host" {
		return false // not the lookforplayers message
	}

	channel, _ := s.Channel(message.ChannelID)
	if channel.Type != discordgo.ChannelTypeGuildText {
		return false // not from a guild
	}
	return true
}

func (l *LookCommand) getPlayers(s *discordgo.Session, message *discordgo.Message, reactionUser string, add bool, active bool) (hostID string, activePlayers []string, backupPlayers map[string]bool, neededplayers int) {
	//Trim out mention
	hostID = strings.TrimRight(strings.TrimLeft(message.Embeds[0].Fields[0].Value, "<@"), ">")
	neededPlayers, _ := strconv.Atoi(message.Embeds[0].Fields[1].Value)

	//Get all players from message
	var playersID []string
	var backupPlayersID []string
	//Active + Backup players (field 3 and 4)
	for i := 3; i <= 4; i++ {
		playersMention := strings.Split(message.Embeds[0].Fields[i].Value, "\n")
		for _, player := range playersMention {
			if strings.HasSuffix(player, "\u200b") && player != "\u200b" {
				//put the backup players in a different array
				ID := strings.TrimRight(strings.TrimLeft(player, "<@"), ">\u200b")
				backupPlayersID = append(backupPlayersID, ID)
			} else if player != "\u200b" {
				ID := strings.TrimRight(strings.TrimLeft(player, "<@"), ">")
				playersID = append(playersID, ID)
			}
		}
	}

	if reactionUser != hostID {
		//There's a better way to do this but i don't know how... (it works tough)
		activePlayerIndex := 999
		//Append players without the host and bot
		for index, ID := range playersID {
			if ID == reactionUser {
				//player in list
				activePlayerIndex = index
			}
		}

		backupPlayerIndex := 999
		//Append backup players without the host and bot
		for index, ID := range backupPlayersID {
			if ID == reactionUser {
				//player in list
				backupPlayerIndex = index
			}
		}

		if active == true && backupPlayerIndex == 999 {
			if add && activePlayerIndex == 999 {
				playersID = append(playersID, reactionUser)
			}
			if !add && activePlayerIndex != 999 {
				//Remove from array
				playersID = append(playersID[:activePlayerIndex], playersID[activePlayerIndex+1:]...)
			}
		}

		if active == false && activePlayerIndex == 999 {
			if add && backupPlayerIndex == 999 {
				backupPlayersID = append(backupPlayersID, reactionUser)
			}
			if !add && backupPlayerIndex != 999 {
				//Remove from backup array
				backupPlayersID = append(backupPlayersID[:backupPlayerIndex], backupPlayersID[backupPlayerIndex+1:]...)
			}
		}
	}

	backupPlayers = make(map[string]bool)

	if len(playersID) < neededPlayers {
		activePlayers = playersID
		for _, backupPlayer := range backupPlayersID {
			backupPlayers[backupPlayer] = false
		}
		message.Embeds[0].Color = 0x33FF33
	} else {
		activePlayers = playersID[:neededPlayers]
		for _, player := range playersID[neededPlayers:] {
			backupPlayers[player] = true
		}
		for _, backupPlayer := range backupPlayersID {
			backupPlayers[backupPlayer] = false
		}
		message.Embeds[0].Color = 0xFF0000
	}
	return hostID, activePlayers, backupPlayers, neededPlayers
}

func (l *LookCommand) startGame(s *discordgo.Session, r *discordgo.MessageReactionAdd, currentPlayers []string, backupPlayers map[string]bool, neededPlayers int, message *discordgo.Message) {
	//Create voice channel
	//TODO: Make configurable in config file!
	conf, isLFP, err := l.checkConfig(r.GuildID, r.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}
	if !isLFP {
		log.Println("not in lfp")
		return
	}

	hiveconf, _, err := l.checkHiveConfig(r.GuildID, conf.HiveChannelID)
	if err != nil {
		log.Println(err)
		return
	}

	//Make the voice channel
	h := hive.NewHiveCommand(l.db)
	channelSize, err := strconv.Atoi(message.Embeds[0].Fields[1].Value)
	if err != nil {
		log.Println(err)
		return
	}
	channel, err := h.CreateVoiceChannel(s, hiveconf, message.Embeds[0].Title, hiveconf.VoiceCategoryID, r.GuildID, channelSize, false)
	if err != nil {
		log.Println(err)
		return
	}

	//Delete message first to prevent players being notified multiple times when emoji spam (Dirk proofing)
	err = s.ChannelMessageDelete(r.ChannelID, r.MessageID)
	if err != nil {
		return
	}

	//Notify players, except the host
	err = l.messagePlayers(s, currentPlayers[1:], message.Embeds[0], fmt.Sprintf("%s is starting now! You can join the channel here! <#%s>\nIf this does not show up, you make one yourself with `/hive type voice name:%s size:%s` in the request channel", message.Embeds[0].Title, channel.ID, message.Embeds[0].Title, message.Embeds[0].Fields[1].Value))
	if len(currentPlayers) < neededPlayers && len(backupPlayers) != 0 {
		var backupID []string
		for player := range backupPlayers {
			backupID = append(backupID, player)
		}
		//Message first backup players
		err = l.messagePlayers(s, backupID[:neededPlayers-len(currentPlayers)], message.Embeds[0], fmt.Sprintf("%s is starting now! You can join the channel here! <#%s>\nIf this does not show up, you make one yourself with `/hive type voice name:%s size:%s` in the request channel", message.Embeds[0].Title, channel.ID, message.Embeds[0].Title, message.Embeds[0].Fields[1].Value))
		//Message host about backup players
		err = l.messagePlayers(s, currentPlayers[:1], message.Embeds[0], fmt.Sprintf("I have notified every joined player and needed backup player(s)! You can join the channel here! <#%s>\nIf this does not show up, you make one yourself with `/hive type voice name:%s size:%s` in the request channel", channel.ID, message.Embeds[0].Title, message.Embeds[0].Fields[1].Value))
	} else {
		//Message host
		err = l.messagePlayers(s, currentPlayers[:1], message.Embeds[0], fmt.Sprintf("I have notified every joined player! You can join the channel here! <#%s>\nIf this does not show up, you make one yourself with `/hive type voice name:%s size:%s` in the request channel", channel.ID, message.Embeds[0].Title, message.Embeds[0].Fields[1].Value))
	}
	if err != nil {
		return
	}
}

func (l *LookCommand) messagePlayers(s *discordgo.Session, playerList []string, embed *discordgo.MessageEmbed, message string) error {
	embed.Fields = embed.Fields[:5]
	for _, user := range playerList {
		dmChannel, _ := s.UserChannelCreate(user)
		messageSend := &discordgo.MessageSend{
			Content: message,
			Embed:   embed,
		}
		_, err := s.ChannelMessageSendComplex(dmChannel.ID, messageSend)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *LookCommand) handleJoinReaction(currentPlayers []string, backupPlayers map[string]bool, message *discordgo.Message, s *discordgo.Session) {
	activePlayersString := "\u200b"
	backupPlayersString := "\u200b"

	if len(currentPlayers) != 0 {
		activePlayersString = ""
		for _, player := range currentPlayers {
			activePlayersString += fmt.Sprintf("<@%s>\n", player)
		}
		if len(backupPlayers) != 0 {
			backupPlayersString = ""
			for player, active := range backupPlayers {
				if active {
					//If active join selected
					backupPlayersString += fmt.Sprintf("<@%s>\n", player)
				} else {
					//If backup selected
					backupPlayersString += fmt.Sprintf("<@%s>\u200b\n", player)
				}
			}
		}
	}

	message.Embeds[0].Fields[3].Value = activePlayersString
	message.Embeds[0].Fields[4].Value = backupPlayersString

	_, err := s.ChannelMessageEditEmbed(message.ChannelID, message.ID, message.Embeds[0])
	if err != nil {
		log.Println(err)
	}
}

// Info return the commands in this package
func (l *LookCommand) Info() []command.Command {
	return []command.Command{}
}
