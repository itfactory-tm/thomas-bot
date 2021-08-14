package hive

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func (h *HiveCommand) handleJoin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent || i.MessageComponentData().CustomID != "hive_join" {
		return
	}

	h.join(s, i.GuildID, i.Member.User.ID, i.ChannelID, i.Message)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
}

func (h *HiveCommand) handleReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Println("Cannot get message of reaction", r.ChannelID)
		return
	}

	h.join(s, r.GuildID, r.UserID, r.ChannelID, message)
}

func (h *HiveCommand) join(s *discordgo.Session, guildID, userID, channelID string, message *discordgo.Message) {
	if message.Author.ID != s.State.User.ID {
		return // not the bot user
	}

	if len(message.Embeds) < 1 {
		return // not an embed
	}

	if len(message.Embeds[0].Fields) < 2 {
		return // not the correct embed
	}

	if message.Embeds[0].Title != "Hive Channel" {
		return // not the hive message
	}

	channel, err := s.Channel(message.Embeds[0].Fields[len(message.Embeds[0].Fields)-1].Value)
	if err != nil {
		log.Println(err)
		return
	}

	conf, isHive, err := h.getConfigForRequestChannel(guildID, channelID)
	if err != nil {
		log.Println(err)
		return
	}

	if !isHive {
		//s.ChannelMessageSend(channelID, "Sorry category not allowed, try privilege escalating otherwise!")
		return
	}

	if channel.ParentID != conf.VoiceCategoryID && channel.ParentID != conf.TextCategoryID {
		// channel no longer in hive
		return
	}

	err = s.ChannelPermissionSet(channel.ID, userID, discordgo.PermissionOverwriteTypeMember, defaultAllows, 0)
	if err != nil {
		log.Println("Cannot set permissions", err)
		return
	}

	// send message if user was not in channel before, we do keep changing the permissions to handle bugs in old permissions
	inChannel := false
	for _, ow := range channel.PermissionOverwrites {
		if ow.Type == discordgo.PermissionOverwriteTypeMember && ow.ID == userID {
			inChannel = true
			break
		}
	}

	if !inChannel {
		s.ChannelMessageSend(channel.ID, fmt.Sprintf("Welcome <@%s>, you can leave any time by saying `/leave`", userID))
	}
}
