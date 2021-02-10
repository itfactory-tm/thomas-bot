package hive

import (
	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/db"
)

func (h *HiveCommand) getConfigForRequestChannel(m *discordgo.MessageCreate) (*db.HiveConfiguration, bool, error) {
	conf, err := h.db.ConfigForGuild(m.GuildID)
	if err != nil {
		return nil, false, err
	}
	if conf == nil {
		// not in our DB
		return nil, false, nil
	}

	for _, hive := range conf.Hives {
		for _, reqID := range hive.RequestChannelIDs {
			if m.ChannelID == reqID {
				return &hive, true, nil
			}
		}
	}

	// no hive found
	return nil, false, nil
}
