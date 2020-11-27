package moderation

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/discordha"

	"github.com/bwmarrin/discordgo"
)

//TODO: make me configurable
const itfGuestRole = "687568536356257890"

type checkFn func(s *discordgo.Session, m *discordgo.MessageCreate) bool

var checks = []checkFn{
	removeLink,
}

func (m *ModerationCommands) checkMessage(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if msg.Author == nil {
		// reactions are also edit events
		return
	}
	if _, err := m.server.GetDiscordHA().CacheRead("check", fmt.Sprintf("%s%s%s", msg.ChannelID, msg.Author.ID, msg.Content), ""); err != nil {
		return
	}
	m.server.GetDiscordHA().CacheWrite("check", fmt.Sprintf("%s%s%s", msg.ChannelID, msg.Author.ID, msg.Content), "true", time.Minute)
	user, err := m.getUser(s, msg.GuildID, msg.Author.ID)
	if err != nil {
		return
	}

	if isUserSafe(user) {
		return
	}

	toRemove := false
	for _, check := range checks {
		if check(s, msg) {
			// remove if the check sends true
			toRemove = true
			break
		}
	}

	if toRemove {
		s.ChannelMessageDelete(msg.ChannelID, msg.ID)
		log.Printf("Removed message from %s aka %s: %s\n", msg.Author.ID, msg.Author.Username, msg.Message.Content)
		m.notifyUser(s, msg.Author.ID)
	}
}

func (m *ModerationCommands) checkReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	c, err := s.UserChannelCreate(r.UserID)
	if err != nil {
		log.Println(err)
		return
	}

	// is DM, no moderation needed
	if c.ID == r.ChannelID {
		return
	}
	user, err := m.getUser(s, r.GuildID, r.UserID)
	if err != nil {
		log.Printf("Error getting user %s, %q\n", r.UserID, err)
		return
	}

	if isUserSafe(user) {
		return
	}

	var num int
	obj, err := m.server.GetDiscordHA().CacheRead("reaction", r.GuildID+r.UserID, num)
	if errors.Is(err, discordha.ErrorCacheKeyNotExist) {
		obj = 0
	} else if err != nil {
		log.Printf("Error reading reaction cache for user %s, %q\n", r.UserID, err)
		return // ignoring here
	}

	var i int
	switch obj.(type) {
	case int:
		i = obj.(int)
	case float64:
		i = int(obj.(float64))
	}

	i++
	if i > 3 {
		s.MessageReactionRemove(r.ChannelID, r.MessageID, r.Emoji.APIName(), r.UserID)
		m.notifyUserReaction(s, r.UserID)
	}

	m.server.GetDiscordHA().CacheWrite("reaction", r.GuildID+r.UserID, i, 2*time.Minute)
}
