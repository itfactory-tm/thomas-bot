package voice

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func FindVoiceUser(dg *discordgo.Session, guildID, userID string) (string, error) {
	g, err := dg.Guild(guildID)
	if err != nil {
		return "", err
	}

	for _, user := range g.VoiceStates {
		if user.UserID == userID {
			return user.ChannelID, nil
		}
	}

	return "", errors.New("user not in voice")
}
