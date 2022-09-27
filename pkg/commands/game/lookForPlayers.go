package game

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/itfactory-tm/thomas-bot/pkg/sudo"

	"github.com/itfactory-tm/thomas-bot/pkg/util/slash"

	"github.com/itfactory-tm/thomas-bot/pkg/commands/hive"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
	"github.com/itfactory-tm/thomas-bot/pkg/db"
)

var reg = regexp.MustCompile(`^[A-Za-z0-9 ]+$`)

var buttons = []discordgo.MessageComponent{
	discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Join",
				Style:    discordgo.SuccessButton,
				CustomID: "lfp_join",
				Emoji: discordgo.ComponentEmoji{
					Name: "👋",
				},
			},
			discordgo.Button{
				Label:    "Join as Backup",
				Style:    discordgo.SecondaryButton,
				CustomID: "lfp_backup",
				Emoji: discordgo.ComponentEmoji{
					Name: "💾",
				},
			},
		},
	},
	discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Delete",
				Style:    discordgo.DangerButton,
				CustomID: "lfp_delete",
				Emoji: discordgo.ComponentEmoji{
					Name: "🗑",
				},
			},
			discordgo.Button{
				Label:    "Start",
				Style:    discordgo.SuccessButton,
				CustomID: "lfp_start",
				Emoji: discordgo.ComponentEmoji{
					Name: "🎮",
				},
			},
		},
	},
}

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

	registry.RegisterInteractionCreate("lfp_join", l.handleBtnClick)
	registry.RegisterInteractionCreate("lfp_backup", l.handleBtnClick)
	registry.RegisterInteractionCreate("lfp_delete", l.handleBtnClick)
	registry.RegisterInteractionCreate("lfp_start", l.handleBtnClick)
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
	conf, _ := l.db.GetAllConfigurations()
	for _, c := range conf {
		if len(c.LookingForPlayers) > 0 {
			if err := slash.InstallSlashCommand(s, c.GuildID, app); err != nil {
				return fmt.Errorf("error installing lfp in %s: %w", c.GuildID, err)
			}
		}
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

	for _, option := range i.ApplicationCommandData().Options {
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
			if matched := reg.MatchString(name); !matched {
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
			if len(timeString) > 25 {
				l.sendInvisibleInteractionResponse(s, i, "Please enter a valid time, >25 characters is a very weird time")
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
		Data: &discordgo.InteractionResponseData{
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
	return nil, false, errors.New("no LFP configured")
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
				Name:   "Players joined",
				Value:  fmt.Sprintf("1/%d", amount),
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
			},
		},
	}

	message := &discordgo.MessageSend{
		Embed:      embed,
		Components: buttons,
	}

	if roleID != "" {
		message.Content = fmt.Sprintf("<@&%s>", roleID)
	}
	_, err := s.ChannelMessageSendComplex(inviteChannelID, message)

	return err
}

func (l *LookCommand) sendInvisibleInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   64,
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func (l *LookCommand) handleBtnClick(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}
	message := i.Message
	uid := i.Member.User.ID

	defer s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	if i.MessageComponentData().CustomID == "lfp_join" {
		// check if we need to remove
		_, _, playersIDs, backupPlayersIDs := l.getPlayers(message)
		for _, p := range playersIDs {
			if p == uid {
				// handle remove
				_, activePlayers, activeBackupPlayers, backupPlayers, neededPlayers := l.removePlayer(message, uid)
				l.handleJoinReaction(activePlayers, activeBackupPlayers, backupPlayers, neededPlayers, message, s)
				return
			}
		}

		_, activePlayers, activeBackupPlayers, backupPlayers, neededPlayers := l.addPlayer(message, uid)
		l.handleJoinReaction(activePlayers, activeBackupPlayers, backupPlayers, neededPlayers, message, s)
		if message.Embeds[0].Fields[2].Value == "Now!" && len(activePlayers) >= neededPlayers {
			l.startGame(s, i, activePlayers, backupPlayersIDs, neededPlayers, message)
		}

		return
	}

	if i.MessageComponentData().CustomID == "lfp_backup" {
		_, _, _, backupPlayersIDs := l.getPlayers(message)
		for _, p := range backupPlayersIDs {
			if p == uid {
				// handle remove
				_, activePlayers, activeBackupPlayers, backupPlayers, neededPlayers := l.removeBackup(message, uid)
				l.handleJoinReaction(activePlayers, activeBackupPlayers, backupPlayers, neededPlayers, message, s)
				return
			}
		}

		_, activePlayers, activeBackupPlayers, backupPlayers, neededPlayers := l.addBackup(message, uid)
		l.handleJoinReaction(activePlayers, activeBackupPlayers, backupPlayers, neededPlayers, message, s)
		return
	}

	hostID, neededPlayers, playersIDs, backupPlayersIDs := l.getPlayers(message)
	_, activePlayers := l.buildBackup(message, playersIDs, neededPlayers)
	//If host of LFP or gamer admin
	if uid == hostID || sudo.IsItfGameAdmin(uid) {
		if i.MessageComponentData().CustomID == "lfp_delete" {
			//Delete message first to prevent players being notified multiple times when emoji spam (Dirk proofing)
			err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
			if err != nil {
				return
			}
			//Notify players
			err = l.messagePlayers(s, activePlayers, message.Embeds[0], fmt.Sprintf("The invite for %s has been deleted by the host.", message.Embeds[0].Title))
			if err != nil {
				return
			}
		}

		if i.MessageComponentData().CustomID == "lfp_start" {
			l.startGame(s, i, activePlayers, backupPlayersIDs, neededPlayers, message)
		}
	} else if i.MessageComponentData().CustomID == "lfp_delete" { //If delete was pressed as normal player
		_, activePlayers, activeBackupPlayers, backupPlayers, _ := l.removePlayer(message, uid)
		l.handleJoinReaction(activePlayers, activeBackupPlayers, backupPlayers, neededPlayers, message, s)
		return
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

	return channel.Type == discordgo.ChannelTypeGuildText
}

func (l *LookCommand) getPlayerIndexes(playersID, backupPlayersID []string, reactionUser string) (int, int) {
	//There's a better way to do this but i don't know how... (it works tough)
	activePlayerIndex := -1
	//Check if the user is in the list
	for index, ID := range playersID {
		if ID == reactionUser {
			//player in list
			activePlayerIndex = index
		}
	}

	backupPlayerIndex := -1
	for index, ID := range backupPlayersID {
		if ID == reactionUser {
			backupPlayerIndex = index
		}
	}
	return activePlayerIndex, backupPlayerIndex
}

func (l *LookCommand) addPlayer(message *discordgo.Message, reactionUser string) (hostID string, activePlayers, activeBackupPlayers, backupPlayers []string, neededplayers int) {
	hostID, neededPlayers, playersIDs, backupPlayersIDs := l.getPlayers(message)

	if reactionUser != hostID {
		activePlayerIndex, backupPlayerIndex := l.getPlayerIndexes(playersIDs, backupPlayersIDs, reactionUser)

		// if in neither of the lists add them to players
		if backupPlayerIndex == -1 && activePlayerIndex == -1 {
			playersIDs = append(playersIDs, reactionUser)
		}

		// if in backup move them to active
		if backupPlayerIndex != -1 && activePlayerIndex == -1 {
			backupPlayersIDs = append(backupPlayersIDs[:backupPlayerIndex], backupPlayersIDs[backupPlayerIndex+1:]...)
			playersIDs = append(playersIDs, reactionUser)
		}
	}

	activeBackupPlayers, activePlayers = l.buildBackup(message, playersIDs, neededPlayers)
	return hostID, activePlayers, activeBackupPlayers, backupPlayersIDs, neededPlayers
}

func (l *LookCommand) removePlayer(message *discordgo.Message, reactionUser string) (hostID string, activePlayers, activeBackupPlayers, backupPlayers []string, neededplayers int) {
	hostID, neededPlayers, playersIDs, backupPlayersIDs := l.getPlayers(message)

	if reactionUser != hostID {
		activePlayerIndex, backupPlayerIndex := l.getPlayerIndexes(playersIDs, backupPlayersIDs, reactionUser)

		// remove from both lists! We don't want to see leavers at all
		if activePlayerIndex != -1 {
			playersIDs = append(playersIDs[:activePlayerIndex], playersIDs[activePlayerIndex+1:]...)
		}
		if backupPlayerIndex != -1 {
			backupPlayersIDs = append(backupPlayersIDs[:backupPlayerIndex], backupPlayersIDs[backupPlayerIndex+1:]...)
		}
	}

	activeBackupPlayers, activePlayers = l.buildBackup(message, playersIDs, neededPlayers)
	return hostID, activePlayers, activeBackupPlayers, backupPlayersIDs, neededPlayers
}

func (l *LookCommand) addBackup(message *discordgo.Message, reactionUser string) (hostID string, activePlayers, activeBackupPlayers, backupPlayers []string, neededplayers int) {
	hostID, neededPlayers, playersIDs, backupPlayersIDs := l.getPlayers(message)

	if reactionUser != hostID {
		activePlayerIndex, backupPlayerIndex := l.getPlayerIndexes(playersIDs, backupPlayersIDs, reactionUser)

		// no longer be an active player
		if activePlayerIndex != -1 {
			playersIDs = append(playersIDs[:activePlayerIndex], playersIDs[activePlayerIndex+1:]...)
		}

		// if not a backup today, become one
		if backupPlayerIndex == -1 {
			backupPlayersIDs = append(backupPlayersIDs, reactionUser)
		}

	}

	activeBackupPlayers, activePlayers = l.buildBackup(message, playersIDs, neededPlayers)
	return hostID, activePlayers, activeBackupPlayers, backupPlayersIDs, neededPlayers
}

func (l *LookCommand) removeBackup(message *discordgo.Message, reactionUser string) (hostID string, activePlayers, activeBackupPlayers, backupPlayers []string, neededplayers int) {
	hostID, neededPlayers, playersIDs, backupPlayersIDs := l.getPlayers(message)

	if reactionUser != hostID {
		_, backupPlayerIndex := l.getPlayerIndexes(playersIDs, backupPlayersIDs, reactionUser)

		if backupPlayerIndex != -1 {
			backupPlayersIDs = append(backupPlayersIDs[:backupPlayerIndex], backupPlayersIDs[backupPlayerIndex+1:]...)
		}
	}

	activeBackupPlayers, activePlayers = l.buildBackup(message, playersIDs, neededPlayers)
	return hostID, activePlayers, activeBackupPlayers, backupPlayersIDs, neededPlayers
}

func (l *LookCommand) buildBackup(message *discordgo.Message, playersIDs []string, neededPlayers int) ([]string, []string) {
	var activeBackupPlayers []string
	var activePlayers []string

	if len(playersIDs) < neededPlayers {
		activePlayers = playersIDs
		message.Embeds[0].Color = 0x33FF33
	} else {
		activePlayers = playersIDs[:neededPlayers]
		activeBackupPlayers = playersIDs[neededPlayers:]
		message.Embeds[0].Color = 0xFF0000
	}
	return activeBackupPlayers, activePlayers
}

func (l *LookCommand) getPlayers(message *discordgo.Message) (string, int, []string, []string) {
	//Trim out mention
	hostID := strings.TrimRight(strings.TrimLeft(message.Embeds[0].Fields[0].Value, "<@"), ">")
	//neededPlayers is the number y in x/y
	var neededPlayers int
	neededPlayersSplit := strings.Split(message.Embeds[0].Fields[1].Value, "/")
	//TODO: Remove this code in the future
	//This if statement is for older lookingforplayers embeds (players needed = 8 instead of Players joined = 2/8)
	if len(neededPlayersSplit) > 1 {
		//New - Players joined = 2/8
		neededPlayers, _ = strconv.Atoi(neededPlayersSplit[1])
	} else {
		//Old - players needed = 8
		neededPlayers, _ = strconv.Atoi(message.Embeds[0].Fields[1].Value)
	}

	//Get all players from message
	var playersIDs []string
	var backupPlayersIDs []string
	//Active + Backup players (field 3 and 4)
	for i := 3; i <= 4; i++ {
		playersMention := strings.Split(message.Embeds[0].Fields[i].Value, "\n")
		for _, player := range playersMention {
			if strings.HasSuffix(player, "\u200b") && player != "\u200b" {
				//put the backup players in a different array
				ID := strings.TrimRight(strings.TrimLeft(player, "<@"), ">\u200b")
				backupPlayersIDs = append(backupPlayersIDs, ID)
			} else if player != "\u200b" {
				ID := strings.TrimRight(strings.TrimLeft(player, "<@"), ">")
				playersIDs = append(playersIDs, ID)
			}
		}
	}
	return hostID, neededPlayers, playersIDs, backupPlayersIDs
}

func (l *LookCommand) startGame(s *discordgo.Session, i *discordgo.InteractionCreate, activePlayers, backupPlayers []string, neededPlayers int, message *discordgo.Message) {
	//Create voice channel
	//TODO: Make configurable in config file!
	conf, isLFP, err := l.checkConfig(i.GuildID, i.ChannelID)
	if err != nil {
		log.Println(err)
		return
	}
	if !isLFP {
		log.Println("not in lfp")
		return
	}

	hiveconf, _, err := l.checkHiveConfig(i.GuildID, conf.HiveChannelID)
	if err != nil {
		log.Println(err)
		return
	}

	//Make the voice channel
	h := hive.NewHiveCommand(l.db)
	channel, err := h.CreateVoiceChannel(s, hiveconf, message.Embeds[0].Title, hiveconf.VoiceCategoryID, i.GuildID, neededPlayers, false)
	if err != nil {
		log.Println(err)
		return
	}

	//Delete message first to prevent players being notified multiple times when button spam (Dirk proofing)
	err = s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
	if err != nil {
		return
	}

	//Notify players, except the host
	err = l.messagePlayers(s, activePlayers[1:], message.Embeds[0], fmt.Sprintf("%s is starting now! You can join the channel here! <#%s>\nIf this does not show up, you make one yourself with `/hive type voice name:%s size:%d` in the request channel", message.Embeds[0].Title, channel.ID, message.Embeds[0].Title, neededPlayers))
	if len(activePlayers) < neededPlayers && len(backupPlayers) != 0 {
		//Find out how many backup players need to be invited
		backupsToAdd := neededPlayers - len(activePlayers)
		if backupsToAdd > len(backupPlayers) {
			backupsToAdd = len(backupPlayers)
		}
		//Message needed backup players
		err = l.messagePlayers(s, backupPlayers[:backupsToAdd], message.Embeds[0], fmt.Sprintf("%s is starting now! You can join the channel here! <#%s>\nIf this does not show up, you make one yourself with `/hive type voice name:%s size:%d` in the request channel", message.Embeds[0].Title, channel.ID, message.Embeds[0].Title, neededPlayers))
		//Get needed backup players
		var backupPlayersString string
		for _, backupPlayer := range backupPlayers[:backupsToAdd] {
			backupPlayersString += fmt.Sprintf("\n<@%s>", backupPlayer)
		}
		//Message host about backup players
		err = l.messagePlayers(s, activePlayers[:1], message.Embeds[0], fmt.Sprintf("I have notified every joined player and needed backup player(s)! You can join the channel here! <#%s>\nIf this does not show up, you make one yourself with `/hive type voice name:%s size:%d` in the request channel\n**Notified backup players:**%s", channel.ID, message.Embeds[0].Title, neededPlayers, backupPlayersString))
	} else {
		//Message host
		err = l.messagePlayers(s, activePlayers[:1], message.Embeds[0], fmt.Sprintf("I have notified every joined player! You can join the channel here! <#%s>\nIf this does not show up, you make one yourself with `/hive type voice name:%s size:%d` in the request channel", channel.ID, message.Embeds[0].Title, neededPlayers))
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

func (l *LookCommand) handleJoinReaction(activePlayers, activeBackupPlayers, backupPlayers []string, neededPlayers int, message *discordgo.Message, s *discordgo.Session) {
	activePlayersString := "\u200b"
	backupPlayersString := "\u200b"

	if len(activePlayers) != 0 {
		activePlayersString = ""
		for _, player := range activePlayers {
			activePlayersString += fmt.Sprintf("<@%s>\n", player)
		}
		if len(backupPlayers) != 0 || len(activeBackupPlayers) != 0 {
			backupPlayersString = ""
			for _, player := range activeBackupPlayers {
				//join active selected
				backupPlayersString += fmt.Sprintf("<@%s>\n", player)
			}
			for _, player := range backupPlayers {
				//If backup selected, we put a non breaking space to know in the future that the user selected Join as Backup
				backupPlayersString += fmt.Sprintf("<@%s>\u200b\n", player)
			}
		}
	}

	//Players needed value
	message.Embeds[0].Fields[1].Value = fmt.Sprintf("%d/%d", len(activePlayers), neededPlayers)
	//Joined Players list
	message.Embeds[0].Fields[3].Value = activePlayersString
	//Backup Players list
	message.Embeds[0].Fields[4].Value = backupPlayersString

	_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Components: buttons,
		Embed:      message.Embeds[0],
		ID:         message.ID,
		Channel:    message.ChannelID,
	})
	if err != nil {
		log.Println(err)
	}
}

// Info return the commands in this package
func (l *LookCommand) Info() []command.Command {
	return []command.Command{}
}
