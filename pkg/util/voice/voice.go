package voice

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

const itfDiscord = "687565213943332875"

func FindVoiceUser(dg *discordgo.Session, guildID, userID string) (string, error) {
	// may regret this
	if guildID == "" {
		guildID = itfDiscord
	}
	g, err := dg.State.Guild(guildID)
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
