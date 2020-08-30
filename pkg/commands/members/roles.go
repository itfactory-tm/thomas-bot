package members

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// TODO: add ability for commands to query config
const prefix = "tm"

func (m *MemberCommands) sayRole(s *discordgo.Session, msg *discordgo.MessageCreate) {
	ch, err := s.UserChannelCreate(msg.Author.ID)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, "Cannot DM user")
		return
	}
	if ch.ID != msg.ChannelID && msg.Message.Content == fmt.Sprintf("%s!role", prefix) {
		s.ChannelMessageDelete(msg.ChannelID, msg.Message.ID)
	}

	m.SendRoleDM(s, msg.Author.ID)
}

// SendRoleDM sends a role selection DM to the user
func (m *MemberCommands) SendRoleDM(s *discordgo.Session, userID string) {
	ch, err := s.UserChannelCreate(userID)
	if err != nil {
		return
	}

	msg, err := s.ChannelMessageSend(ch.ID, roleMessage)
	if err != nil {
		log.Println("Role DM error", err)
		return
	}

	for _, emoji := range roleEmoji.Keys() {
		err := s.MessageReactionAdd(ch.ID, msg.ID, emoji.(string))
		if err != nil {
			log.Printf("Error adding help emoji: %q\n", err)
		}
	}
}

func (m *MemberCommands) handleRoleReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		return
	}

	if message.Author.ID != s.State.User.ID {
		return // not the bot user
	}

	if message.Content != roleMessage {
		return // not the role message
	}

	wantedRole, roleExists := roleEmoji.Get(r.Emoji.MessageFormat())
	if !roleExists {
		log.Printf("Role emoji %s not found", r.Emoji.MessageFormat())
	}

	ch, err := s.UserChannelCreate(r.UserID)
	if err != nil {
		log.Printf("Cannot DM user", err)
		return
	}

	member, err := s.GuildMember(itfDiscord, r.UserID)
	if err == nil {
		for _, role := range member.Roles {
			if role == wantedRole {
				s.ChannelMessageSend(ch.ID, "Oopsie! You already have the role you requested!")
				return
			}
		}
	}

	s.ChannelMessageSend(ch.ID, "Thank you! I have asked our moderators for permissions to assign the role you asked.")
	if r.Emoji.MessageFormat() == "👩‍🏫" {
		s.ChannelMessageSend(ch.ID, "Not already working at Thomas More? We're hiring! http://werkenbij.thomasmore.be/")
	}

	msg, err := s.ChannelMessageSend(roleChannelID, fmt.Sprintf("<@%s> wants role <@&%s>\n Allow/Deny or Remove all others and assign requested role?", r.UserID, wantedRole))
	if err != nil {
		return // let's handle this later
	}
	s.MessageReactionAdd(roleChannelID, msg.ID, "✅")
	s.MessageReactionAdd(roleChannelID, msg.ID, "❌")
	s.MessageReactionAdd(roleChannelID, msg.ID, "☝️")
}

func (m *MemberCommands) handleRolePermissionReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.ChannelID != roleChannelID {
		return
	}

	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		return
	}

	if r.Emoji.MessageFormat() != "✅" && r.Emoji.MessageFormat() != "☝️" {
		return
	}
	matches := userIDRoleIDRegex.FindAllStringSubmatch(message.Content, -1)
	if len(matches) != 1 {
		return /// invalid message
	}
	if len(matches[0]) != 3 {
		return /// invalid message
	}

	userID := matches[0][1]
	roleID := matches[0][2]

	if r.Emoji.MessageFormat() == "☝️" {
		member, err := s.GuildMember(itfDiscord, userID)
		if err != nil {
			s.ChannelMessageSend(roleChannelID, fmt.Sprintf("Error getting roles of <@%s>, aborting operation: %q\n", userID, err))
			return
		}
		for _, role := range member.Roles {
			s.GuildMemberRoleRemove(itfDiscord, userID, role)
		}
	}

	err = s.GuildMemberRoleAdd(itfDiscord, userID, roleID)
	if err != nil {
		s.ChannelMessageSend(roleChannelID, fmt.Sprintf("Error assigning role %q\n", err))
		log.Printf("Error assigning role %q\n", err)
	}
}
