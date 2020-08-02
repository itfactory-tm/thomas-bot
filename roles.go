package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/elliotchance/orderedmap"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

const roleChannelID = "739512119338467449" // one day we need to stop doing these...

const roleMessage = `We need to assign you a role inside our Discord which will help you  gain access to the class specific channels.
Select the following emoji(s) for roles you want to request, note that our moderation team has to approve these first.
1️⃣: 1ITF Student
2️⃣: 2ITF Student
3️⃣: 3ITF Student
👩‍🎓: Alumni
👩‍💻: OHO Student
👩‍🏫: Teacher`

var userIDRoleIDRegex = *regexp.MustCompile(`<@(.*)> wants role <@&(.*)>.*`)

var roleEmoji = orderedmap.NewOrderedMap()

func init() {
	registerCommand(command.Command{
		Name:        "role",
		Category:    command.CategoryAlgemeen,
		Description: "Modify your ITFactory Discord role",
		Hidden:      false,
		Handler:     sayRole,
	})

	// very upset Discord does not support non-binary emoji
	roleEmoji.Set("1️⃣", "687567949795557386") // 1ITF
	roleEmoji.Set("2️⃣", "687568334379679771") // 2ITF
	roleEmoji.Set("3️⃣", "687568470820388864") // 3ITF
	roleEmoji.Set("👩‍🎓", "688368287255494702") // Alumni
	roleEmoji.Set("👩‍🏫", "687567374198767617") // Teacher
	roleEmoji.Set("👩‍💻", "689844328528478262") // OHO
}

func sayRole(s *discordgo.Session, m *discordgo.MessageCreate) {
	ch, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Cannot DM user")
		return
	}

	if ch.ID != m.ChannelID && m.Message.Content == fmt.Sprintf("%s!role", c.Prefix) {
		s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
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

func handleRoleReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		return
	}

	if message.Author.ID != s.State.User.ID {
		return // not the bot user
	}

	if r.UserID == s.State.User.ID {
		return // the bot itself reacted
	}

	if r.ChannelID == roleChannelID {
		handleRolePermissionReaction(s, r, message)
		return
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

func handleRolePermissionReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd, message *discordgo.Message) {
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

	err := s.GuildMemberRoleAdd(itfDiscord, userID, roleID)
	if err != nil {
		s.ChannelMessageSend(roleChannelID, fmt.Sprintf("Error assigning role %q\n", err))
		log.Printf("Error assigning role %q\n", err)
	}
}
