package hive

import (
	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/db"
)

func (h *HiveCommand) getConfigForRequestChannel(guildID, challelID string) (*db.HiveConfiguration, bool, error) {
	conf, err := h.db.ConfigForGuild(guildID)
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

func (h *HiveCommand) getConfigForRequestCategory(s *discordgo.Session, guildID, channelID string) (*db.HiveConfiguration, bool, error) {
	conf, err := h.db.ConfigForGuild(guildID)
	if err != nil {
		return nil, false, err
	}
	if conf == nil {
		// not in our DB
		return nil, false, nil
	}

	channel, err := s.Channel(channelID)
	if err != nil {
		return nil, false, err
	}

	for _, hive := range conf.Hives {
		if channel.ParentID == hive.VoiceCategoryID || channel.ParentID == hive.TextCategoryID {
			return &hive, true, nil
		}
	}

	// no hive found
	return nil, false, nil
}

func (h *HiveCommand) isPrivilegedChannel(channelID string, conf *db.HiveConfiguration) bool {
	switch channelID {
	case conf.VoiceCategoryID:
	case conf.TextCategoryID:
		return true

	}

	for _, requestID := range conf.RequestChannelIDs {
		if channelID == requestID {
			return true
		}
	}

	return false
}
