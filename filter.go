package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/discordha"

	"github.com/bwmarrin/discordgo"
)

type checkFn func(s *discordgo.Session, m *discordgo.MessageCreate) bool

var checks = []checkFn{
	removeLink,
}

const itfGuestRole = "687568536356257890"

func checkMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author == nil {
		// reactions are also edit events
		return
	}
	if _, err := ha.CacheRead("check", fmt.Sprintf("%s%s%s", m.ChannelID, m.Author.ID, m.Content), ""); err != nil {
		return
	}
	ha.CacheWrite("check", fmt.Sprintf("%s%s%s", m.ChannelID, m.Author.ID, m.Content), "true", time.Minute)
	user, err := getUser(m.GuildID, m.Author.ID)
	if err != nil {
		return
	}

	if isUserSafe(user) {
		return
	}

	toRemove := false
	for _, check := range checks {
		if check(s, m) {
			// remove if the check sends true
			toRemove = true
			break
		}
	}

	if toRemove {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		log.Printf("Removed message from %s aka %s: %s\n", m.Author.ID, m.Author.Username, m.Message.Content)
		notifyUser(m.Author.ID)
	}
}

func checkReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	user, err := getUser(r.GuildID, r.UserID)
	if err != nil {
		log.Println(err)
		return
	}

	if isUserSafe(user) {
		return
	}

	var num int
	obj, err := ha.CacheRead("reaction", r.GuildID+r.UserID, num)
	if errors.Is(err, discordha.ErrorCacheKeyNotExist) {
		obj = 0
	} else if err != nil {
		log.Println(err)
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
		notifyUserReaction(r.UserID)
	}

	ha.CacheWrite("reaction", r.GuildID+r.UserID, i, 2*time.Minute)
}

func getUser(gid, uid string) (*discordgo.Member, error) {
	obj, err := ha.CacheRead("user", gid+uid, &discordgo.Member{})
	if err != nil {
		user, err := dg.GuildMember(gid, uid)
		if err != nil {
			return nil, err
		}

		ha.CacheWrite("user", gid+uid, user, 2*time.Minute)
		return user, nil
	}

	return obj.(*discordgo.Member), nil
}

// checks if the user is somebody we should trust
func isUserSafe(m *discordgo.Member) bool {
	safe := true // i trust people on first sight
	for _, role := range m.Roles {
		if role == itfGuestRole {
			safe = false
		}
	}

	return safe
}

func notifyUser(id string) {
	_, err := ha.CacheRead("notify", id, "")
	if !errors.Is(err, discordha.ErrorCacheKeyNotExist) {
		// limit self spam
		return
	}
	c, err := dg.UserChannelCreate(id)
	if err != nil {
		return
	}

	dg.ChannelMessageSend(c.ID, "Hallo! Ik heb een bericht van je verwijderd omdat het inging tegen de Thomas More ITFactory Discord regels.")
	ha.CacheWrite("notify", id, "", 3*time.Minute)
}

func notifyUserReaction(id string) {
	_, err := ha.CacheRead("reactionnotify", id, "")
	if !errors.Is(err, discordha.ErrorCacheKeyNotExist) {
		// limit self spam
		return
	}
	c, err := dg.UserChannelCreate(id)
	if err != nil {
		return
	}

	dg.ChannelMessageSend(c.ID, "Hallo! Ik heb je reactie van je verwijderd omdat het inging tegen de Thomas More ITFactory Discord regels.")
	ha.CacheWrite("reactionnotify", id, "", 3*time.Minute)
}
