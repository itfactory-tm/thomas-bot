package game

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
	"github.com/itfactory-tm/thomas-bot/pkg/db"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TODO: make configurable in config file
//LFPDesk channel id
const lfpDeskID = "828205635249635369"

//LFP Request channel id
const lfpReqID = "828204894187421696"

// LookCommand contains the /lookforplayers command
type LookCommand struct {
	db db.Database
}

// NewSearchCommand gives a new SearchCommand
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

// TODO: Make configurable for specific guilds
// InstallSlashCommands registers the slash commands
func (l *LookCommand) InstallSlashCommands(s *discordgo.Session) error {
	app := discordgo.ApplicationCommand{
		Name:        "lookForPlayers",
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
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "notifyrole",
				Description: "Notify a role with your invitation!",
				Required:    false,
			}, {
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "time",
				Description: "At what time do you want to play? Format hh:mm (example: 15:45)",
				Required:    false,
			},
		},
	}

	cmds, err := s.ApplicationCommands(s.State.User.ID, "773847927910432789")
	if err != nil {
		return err
	}
	exists := false
	for _, cmd := range cmds {
		if cmd.Name == "lookForPlayers" {
			exists = reflect.DeepEqual(app.Options, cmd.Options)
		}
	}

	if !exists {
		_, err = s.ApplicationCommandCreate(s.State.User.ID, "773847927910432789", &app)
	}

	return err
}

func (l *LookCommand) SearchCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//conf, isHive, err := l.checkConfig(i.GuildID, i.ChannelID)
	//if err != nil {
	//	l.sendInteractionResponse(s,i,err.Error())
	//	return
	//}
	//if !isHive {
	//	l.sendInteractionResponse(s,i,"This command only works in Requests channels")
	//	return// not from a guild
	//}

	if i.ChannelID != lfpReqID {
		l.sendInteractionResponse(s, i, "This command only works in Requests channels")
		return
	}

	inviteChannelID := lfpDeskID
	var name, selectedRoleID string
	var amount float64
	timeString := "Now!"
	var ok bool

	for _, option := range i.Data.Options {
		switch option.Name {
		case "game":
			name, ok = option.Value.(string)
			if !ok {
				l.sendInteractionResponse(s, i, "Please enter a valid name.")
				return
			}
			if len(name) < 2 || len(name) > 25 {
				l.sendInteractionResponse(s, i, "Your game needs to be between 2-25 characters long")
				return
			}
			if matched, _ := regexp.MatchString(`^[A-Za-z0-9 ]+$`, name); !matched {
				l.sendInteractionResponse(s, i, "Your game cannot contain any special characters")
				return
			}

		case "amount":
			amount, ok = option.Value.(float64)
			if !ok {
				l.sendInteractionResponse(s, i, "Please enter a valid amount.")
				return
			}
			if amount < 2 || amount > 40 {
				l.sendInteractionResponse(s, i, "Your game needs to contain between 2-40 players")
				return
			}
		case "notifyrole":
			selectedRoleID, ok = option.Value.(string)
			if !ok {
				l.sendInteractionResponse(s, i, "Please enter a valid role.")
				return
			}
			roles, _ := s.GuildRoles(i.GuildID)
			for _, role := range roles {
				if selectedRoleID == role.ID && role.Color != 0x9c9c9c {
					l.sendInteractionResponse(s, i, "Please enter a valid gaming role.")
					return
				}
			}
		case "time":
			timeString, ok = option.Value.(string)
			if !ok {
				l.sendInteractionResponse(s, i, "Please enter a valid time in format 15:45.")
				return
			}
			if _, err := time.Parse("15:04", timeString); err != nil {
				l.sendInteractionResponse(s, i, "Please enter your time in format hh:mm (For example 15:50)")
				return
			}
		}
	}

	err := l.createInviteEmbed(s, i, name, int(amount), timeString, selectedRoleID, inviteChannelID)
	content := fmt.Sprintf("Invite created in <#%v>!", inviteChannelID)
	if err != nil {
		content = err.Error()
	}
	l.sendInteractionResponse(s, i, content)
}

func (l *LookCommand) checkConfig(guildID, challelID string) (*db.HiveConfiguration, bool, error) {
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
			if challelID == reqID {
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
				Name:   "Join",
				Value:  "👋",
				Inline: true,
			}, {
				Name:   "Delete Invite",
				Value:  "🗑️",
				Inline: true,
			}, {
				Name:   "Start game",
				Value:  "🎮",
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
	s.MessageReactionAdd(sentMessage.ChannelID, sentMessage.ID, "🗑️")
	s.MessageReactionAdd(sentMessage.ChannelID, sentMessage.ID, "🎮")
	return nil
}

func (l *LookCommand) sendInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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

func (l *LookCommand) handleReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		return
	}

	if !l.checkEmbed(s, message) {
		return
	}

	hostID, currentPlayers, backupPlayers, neededPlayers := l.getPlayers(s, message)

	if r.Emoji.MessageFormat() == "👋" {
		l.handleJoinReaction(currentPlayers, backupPlayers, message, s)
		if message.Embeds[0].Fields[2].Value == "Now!" && len(currentPlayers) >= neededPlayers {
			l.startGame(s, r, currentPlayers, message, hostID, err)
		}
	}

	if r.Emoji.MessageFormat() == "🗑️" {
		if r.UserID == hostID {
			//Notify players
			l.messagePlayers(s, r, currentPlayers, fmt.Sprintf("The invite for %s has been cancelled.", message.Embeds[0].Title))
		}
	}

	if r.Emoji.MessageFormat() == "🎮" {
		if r.UserID == hostID {
			l.startGame(s, r, currentPlayers, message, hostID, err)
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

	_, currentPlayers, backupPlayers, _ := l.getPlayers(s, message)

	if r.Emoji.MessageFormat() == "👋" {
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

func (l *LookCommand) getPlayers(s *discordgo.Session, message *discordgo.Message) (hostID string, activePlayers []*discordgo.User, backupPlayers []*discordgo.User, neededplayers int) {
	//Trim out mention
	hostID = strings.TrimRight(strings.TrimLeft(message.Embeds[0].Fields[0].Value, "<@"), ">")
	neededPlayers, _ := strconv.Atoi(message.Embeds[0].Fields[1].Value)

	//Get all players from reaction
	reactionPlayers, err := s.MessageReactions(message.ChannelID, message.ID, "👋", 100, "", "")
	if err != nil {
		log.Println(err)
		return
	}

	//Make new array with the hostUser
	hostUser, _ := s.User(hostID)
	joinedPlayers := []*discordgo.User{hostUser}

	//Append ratedPlayers without the host and bot
	for _, player := range reactionPlayers {
		if player.ID != hostID && player.ID != s.State.User.ID {
			joinedPlayers = append(joinedPlayers, player)
		}
	}

	if len(joinedPlayers) < neededPlayers {
		activePlayers = joinedPlayers
		message.Embeds[0].Color = 0x33FF33
	} else {
		activePlayers = joinedPlayers[:neededPlayers]
		backupPlayers = joinedPlayers[neededPlayers:]
		message.Embeds[0].Color = 0xFF0000
	}
	return hostID, activePlayers, backupPlayers, neededPlayers
}

func (l *LookCommand) startGame(s *discordgo.Session, r *discordgo.MessageReactionAdd, currentPlayers []*discordgo.User, message *discordgo.Message, hostID string, err error) {
	//Notify players, except the host
	messagePlayerSuccessful := l.messagePlayers(s, r, currentPlayers[1:], fmt.Sprintf("The game %s is starting now! Your host should make a voice channel soon!\nOr make one yourself with `/hive type voice name:%s size:%s` in the request channel", message.Embeds[0].Title,message.Embeds[0].Title,message.Embeds[0].Fields[1].Value))
	if !messagePlayerSuccessful {
		return
	}
	message.Embeds[0].Fields = message.Embeds[0].Fields[:5]
	messageSend := &discordgo.MessageSend{
		Content: fmt.Sprintf("I have notified every joined player! Here is your invite to notify backup players if needed.\nDon't forget to make a voice channel with `/hive type voice name:%s size:%s` in the request channel",message.Embeds[0].Title,message.Embeds[0].Fields[1].Value),
		Embed:   message.Embeds[0],
	}
	//Dm invite to host
	dmChannel, err := s.UserChannelCreate(hostID)
	_, err = s.ChannelMessageSendComplex(dmChannel.ID, messageSend)
	if err != nil {
		log.Println(err)
	}
}

func (l *LookCommand) messagePlayers(s *discordgo.Session, r *discordgo.MessageReactionAdd, currentPlayers []*discordgo.User, message string) bool {
	//Delete message first to prevent players being notified multiple times when emoji spam (Dirk proofing)
	err := s.ChannelMessageDelete(r.ChannelID, r.MessageID)
	if err != nil {
		return false
	}
	for _, user := range currentPlayers {
		dmChannel, _ := s.UserChannelCreate(user.ID)
		_, err := s.ChannelMessageSend(dmChannel.ID, message)
		if err != nil {
			log.Println(err)
		}
	}
	return true
}

func (l *LookCommand) handleJoinReaction(activePlayers []*discordgo.User, backupPlayers []*discordgo.User, message *discordgo.Message, s *discordgo.Session) {
	activePlayersString := "\u200b"
	backupPlayersString := "\u200b"

	if len(activePlayers) != 0 {
		activePlayersString = ""
		for _, player := range activePlayers {
			activePlayersString += fmt.Sprintf("%s\n", player.Mention())
		}
		if len(backupPlayers) != 0 {
			backupPlayersString = ""
			for _, player := range backupPlayers {
				backupPlayersString += fmt.Sprintf("%s\n", player.Mention())
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
