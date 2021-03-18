package moderation

import (
	"errors"
	"strings"
	"time"

	discordha "github.com/meyskens/discord-ha"

	"github.com/bwmarrin/discordgo"
	"mvdan.cc/xurls/v2"
)

func removeLink(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	parts := strings.Split(m.Message.Content, " ")
	for _, part := range parts {
		rxStrict := xurls.Strict()
		if len(rxStrict.FindAllString(part, -1)) > 0 {
			return true
		}
	}

	return false
}

func (m *ModerationCommands) getUser(dg *discordgo.Session, gid, uid string) (*discordgo.Member, error) {
	obj, err := m.server.GetDiscordHA().CacheRead("user", gid+uid, &discordgo.Member{})
	if err != nil {
		user, err := dg.GuildMember(gid, uid)
		if err != nil {
			return nil, err
		}

		m.server.GetDiscordHA().CacheWrite("user", gid+uid, user, 2*time.Minute)
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

func (m *ModerationCommands) notifyUser(dg *discordgo.Session, id string) {
	_, err := m.server.GetDiscordHA().CacheRead("notify", id, "")
	if !errors.Is(err, discordha.ErrorCacheKeyNotExist) {
		// limit self spam
		return
	}
	c, err := dg.UserChannelCreate(id)
	if err != nil {
		return
	}

	dg.ChannelMessageSend(c.ID, "Hallo! Ik heb een bericht van je verwijderd omdat het inging tegen de Thomas More ITFactory Discord regels.")
	m.server.GetDiscordHA().CacheWrite("notify", id, "", 3*time.Minute)
}

func (m *ModerationCommands) notifyUserReaction(dg *discordgo.Session, id string) {
	_, err := m.server.GetDiscordHA().CacheRead("reactionnotify", id, "")
	if !errors.Is(err, discordha.ErrorCacheKeyNotExist) {
		// limit self spam
		return
	}
	c, err := dg.UserChannelCreate(id)
	if err != nil {
		return
	}

	dg.ChannelMessageSend(c.ID, "Hallo! Ik heb je reactie van je verwijderd omdat het inging tegen de Thomas More ITFactory Discord regels.")
	m.server.GetDiscordHA().CacheWrite("reactionnotify", id, "", 3*time.Minute)
}
