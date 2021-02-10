package hive

import (
	"fmt"
	"log"
	"strings"

	"github.com/itfactory-tm/thomas-bot/pkg/embed"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/sudo"
)

// SayAttendance handles the tm!attendance command
func (h *HiveCommand) SayAttendance(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !sudo.IsAdmin(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is not in the sudoers file. This incident will be reported.", m.Author.ID))
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	UIDs := []string{}
	for _, perm := range channel.PermissionOverwrites {
		if perm.Type == discordgo.PermissionOverwriteTypeMember {
			UIDs = append(UIDs, perm.ID)
		}
	}

	names := []string{}
	for _, uid := range UIDs {
		u, err := s.GuildMember(m.GuildID, uid)
		if err != nil {
			continue // we don't care for now
		}
		name := ""
		if u.Nick != "" {
			name = u.Nick
		} else {
			name = u.User.Username
		}
		names = append(names, name)
	}

	hasMessages := true
	activeUIDs := []string{}
	lastID := ""
	c := 0
	for hasMessages {
		c++
		if c > 20 { // max 20000 messages
			break
		}
		messages, err := s.ChannelMessages(m.ChannelID, 100, lastID, "", "")
		if err != nil {
			break
		}

		for _, message := range messages {
			activeUIDs = append(activeUIDs, message.Author.ID)
		}

		if len(messages) == 0 {
			hasMessages = false
			break
		}

		lastID = messages[len(messages)-1].ID
	}

	activeUIDs = removeDuplicateValues(activeUIDs)
	activeNames := []string{}
	for _, uid := range activeUIDs {
		if uid == s.State.User.ID { // rm the bot
			continue
		}
		u, err := s.GuildMember(m.GuildID, uid)
		if err != nil {
			continue // we don't care for now
		}
		name := ""
		if u.Nick != "" {
			name = u.Nick
		} else {
			name = u.User.Username
		}
		activeNames = append(activeNames, name)
	}

	e := embed.NewEmbed()
	e.SetTitle("Attendance List")
	if len(names) > 0 {
		e.AddField("people in channel", strings.Join(names, "\n"))
		e.AddField("number of people in channel", fmt.Sprintf("%d", len(names)))
	}
	if len(activeNames) > 0 {
		e.AddField("active people in channel", strings.Join(activeNames, "\n"))
		e.AddField("number of active people in channel", fmt.Sprintf("%d", len(activeNames)))
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
	if err != nil {
		log.Println(err)
	}
}

func removeDuplicateValues(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
