package members

import (
	"fmt"
	"log"

	"github.com/itfactory-tm/thomas-bot/pkg/db"
	"github.com/itfactory-tm/thomas-bot/pkg/embed"

	"github.com/bwmarrin/discordgo"
)

func (m *MemberCommands) roleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ch, err := s.UserChannelCreate(i.Member.User.ID)
	if err != nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: "error sending a DM to you",
				Flags:   64, // ephemeral
			},
		})

		if err != nil {
			log.Println(err)
		}
		return
	}

	// if not in DM delete command
	if ch.ID == i.ChannelID {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: "I'm sorry I have no idea which server you are in, please use tm!role in a channel in the Discord server I need to help you with.",
				Flags:   64, // ephemeral
			},
		})

		if err != nil {
			log.Println(err)
		}
		return
	}

	m.SendRoleDM(s, i.GuildID, i.Member.User.ID)

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "I sent you a DM!",
			Flags:   64, // ephemeral
		},
	})

	if err != nil {
		log.Println(err)
	}

}

// SendRoleDM sends a role selection DM to the user
func (m *MemberCommands) SendRoleDM(s *discordgo.Session, guildID, userID string) {
	conf, err := m.db.ConfigForGuild(guildID)
	if err != nil || conf == nil {
		return
	}

	ch, err := s.UserChannelCreate(userID)
	if err != nil {
		return
	}

	if len(conf.RoleManagement.Roles) <= 0 {
		s.ChannelMessageSend(ch.ID, "I'm sorry this server hasn't told me any roles I am allowed to give you :(")
		return
	}

	if conf.RoleManagement.Message != "" {
		_, err := s.ChannelMessageSend(ch.ID, conf.RoleManagement.Message)
		if err != nil {
			log.Println("Role DM error", err)
			return
		}
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		log.Println("Guild error", err)
		return
	}
	e := embed.NewEmbed()
	e.SetTitle("Role Request")
	e.SetAuthor(guildID)

	roles := ""
	for _, crole := range conf.RoleManagement.Roles {
		role := findRole(guild.Roles, crole.ID)
		if role != nil {
			roles += fmt.Sprintf("%s: %s\n", crole.Emoji, role.Name)
		}
	}

	e.AddField("Roles", roles)

	msg, err := s.ChannelMessageSendEmbed(ch.ID, e.MessageEmbed)
	if err != nil {
		log.Println("Role DM error", err)
		return
	}

	for _, crole := range conf.RoleManagement.Roles {
		err := s.MessageReactionAdd(ch.ID, msg.ID, crole.Emoji)
		if err != nil {
			log.Printf("Error adding help emoji: %q\n", err)
		}
	}
}

func findRole(all []*discordgo.Role, want string) *discordgo.Role {
	for _, role := range all {
		if role.ID == want {
			return role
		}
	}
	return nil
}

func (m *MemberCommands) handleRoleReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		return
	}

	if message.Author.ID != s.State.User.ID {
		return // not the bot user
	}

	if len(message.Embeds) <= 0 {
		return // not the role message
	}

	if message.Embeds[0].Title != "Role Request" {
		return // not the role message
	}

	guildID := message.Embeds[0].Author.Name
	conf, err := m.db.ConfigForGuild(guildID)
	if err != nil || conf == nil {
		return
	}

	wantedRole := findRoleWithEmoji(conf.RoleManagement.Roles, r.Emoji.MessageFormat())
	if wantedRole == nil {
		log.Printf("Role emoji %s not found", r.Emoji.MessageFormat())
	}

	ch, err := s.UserChannelCreate(r.UserID)
	if err != nil {
		log.Printf("Cannot DM user", err)
		return
	}

	member, err := s.GuildMember(guildID, r.UserID)
	if err == nil {
		for _, role := range member.Roles {
			if role == wantedRole.ID {
				s.ChannelMessageSend(ch.ID, "Oopsie! You already have the role you requested!")
				return
			}
		}
	}

	s.ChannelMessageSend(ch.ID, "Thank you! I have asked our moderators for permissions to assign the role you asked.")
	if r.Emoji.MessageFormat() == "👩‍🏫" {
		s.ChannelMessageSend(ch.ID, "Not already working at Thomas More? We're hiring! http://werkenbij.thomasmore.be/")
	}

	msg, err := s.ChannelMessageSend(conf.RoleManagement.RoleAdminChannelID, fmt.Sprintf("<@%s> wants role <@&%s>\n Allow/Deny or Remove all others and assign requested role?", r.UserID, wantedRole.ID))
	if err != nil {
		return // let's handle this later
	}
	s.MessageReactionAdd(conf.RoleManagement.RoleAdminChannelID, msg.ID, "✅")
	s.MessageReactionAdd(conf.RoleManagement.RoleAdminChannelID, msg.ID, "❌")
	s.MessageReactionAdd(conf.RoleManagement.RoleAdminChannelID, msg.ID, "☝️")
}

func (m *MemberCommands) handleRolePermissionReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	conf, err := m.db.ConfigForGuild(r.GuildID)
	if err != nil || conf == nil {
		return
	}

	if r.ChannelID != conf.RoleManagement.RoleAdminChannelID {
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
		member, err := s.GuildMember(r.GuildID, userID)
		if err != nil {
			s.ChannelMessageSend(conf.RoleManagement.RoleAdminChannelID, fmt.Sprintf("Error getting roles of <@%s>, aborting operation: %q\n", userID, err))
			return
		}
		for _, role := range member.Roles {
			s.GuildMemberRoleRemove(r.GuildID, userID, role)
		}
	}

	err = s.GuildMemberRoleAdd(r.GuildID, userID, roleID)
	if err != nil {
		s.ChannelMessageSend(conf.RoleManagement.RoleAdminChannelID, fmt.Sprintf("Error assigning role %q\n", err))
		log.Printf("Error assigning role %q\n", err)
		return
	}

	s.ChannelMessageSend(conf.RoleManagement.RoleAdminChannelID, fmt.Sprintf("Assigned <@&%s> role for <@%s>", roleID, userID))
}

func findRoleWithEmoji(roles []db.Role, wantedEmjoi string) *db.Role {
	for _, role := range roles {
		if role.Emoji == wantedEmjoi {
			return &role
		}
	}

	return nil
}
