package members

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/db"
)

func (m *MemberCommands) roleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "I cannot do this in DM, sorry",
			},
		})
		return
	}
	ch, err := s.UserChannelCreate(i.Member.User.ID)
	if err != nil {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "error sending a DM to you",
				Flags:   64, // ephemeral
			},
		})

		if err != nil {
			log.Println(err)
		}
		return
	}

	if ch.ID == i.ChannelID {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "I'm sorry I have no idea which server you are in, please use tm!role in a channel in the Discord server I need to help you with.",
				Flags:   64, // ephemeral
			},
		})

		if err != nil {
			log.Println(err)
		}
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "I sent you a DM!",
			Flags:   64, // ephemeral
		},
	})

	if err != nil {
		log.Println(err)
	}

	m.SendRoleDM(s, i.GuildID, i.Member.User.ID)

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

	if len(conf.RoleManagement.RoleSets) <= 0 {
		s.ChannelMessageSend(ch.ID, "I'm sorry this server hasn't told me any roles I am allowed to give you :(")
		return
	}

	guild, err := s.Guild(guildID)
	if err != nil {
		log.Println("Guild error", err)
		return
	}

	for _, rs := range conf.RoleManagement.RoleSets {
		roles := []discordgo.SelectMenuOption{}
		for _, crole := range rs.Roles {
			role := findRole(guild.Roles, crole.ID)
			if role != nil {
				roles = append(roles, discordgo.SelectMenuOption{
					Label:       role.Name,
					Value:       role.ID,
					Description: role.Name,
					Emoji: discordgo.ComponentEmoji{
						Name: crole.Emoji,
					},
					Default: false,
				})
			}
		}

		// discord requires the maximum options to be as long as the list but not more than 25
		maxValues := len(roles)
		if maxValues > 25 {
			maxValues = 25
		}
		minValues := 1

		_, err = s.ChannelMessageSendComplex(ch.ID, &discordgo.MessageSend{
			Content: rs.Message,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							MinValues:   &minValues,
							MaxValues:   maxValues,
							CustomID:    "rolereq--" + guildID,
							Placeholder: "Select the roles you want to request",
							Options:     roles,
						},
					},
				},
			},
		})

		if err != nil {
			log.Println(err)
		}

		time.Sleep(3 * time.Second)
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

func (m *MemberCommands) handleRoleRequest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := strings.Split(i.MessageComponentData().CustomID, "--")
	if len(data) < 2 {
		return // not valid ID
	}
	guildID := data[1]
	conf, err := m.db.ConfigForGuild(guildID)
	if err != nil || conf == nil {
		return // no guild data
	}

	ch, err := s.UserChannelCreate(i.User.ID)
	if err != nil {
		log.Printf("Cannot DM user", err)
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	member, err := s.GuildMember(guildID, i.User.ID)
	if err != nil {
		log.Println("error looking up memeber", err)
		return
	}

L:
	for _, val := range i.MessageComponentData().Values {

		var role *discordgo.Role
		guildRoles, err := s.GuildRoles(guildID)
		if err != nil {
			log.Println("error getting guild roles", err)
			return
		}
		for _, gr := range guildRoles {
			if gr.ID == val {
				role = gr
				break
			}
		}

		if role == nil {
			s.ChannelMessageSend(ch.ID, "Oh no! I cannot find that role any longer...")
			continue L
		}

		for _, mr := range member.Roles {
			if mr == val {
				// ignoring this for the current run
				//s.ChannelMessageSend(ch.ID, fmt.Sprintf("Oopsie! You already have the role %q, no worries I will not re-request it!", role.Name))
				//continue L

				log.Printf("Role %q from %q will be re-requested", role.Name, member.User.Username)
			}
		}

		s.ChannelMessageSend(ch.ID, fmt.Sprintf("Thank you! I have asked our moderators for permissions to assign the role %q", role.Name))
		var configRole *db.Role
		for _, rs := range conf.RoleManagement.RoleSets {
			for _, crole := range rs.Roles {
				if crole.ID == val {
					configRole = &crole
					break
				}
			}
		}
		if configRole != nil && configRole.AutoApprove {
			time.Sleep(time.Second)
			s.ChannelMessageSend(ch.ID, fmt.Sprintf("I have assigned the role %q to you", role.Name))
			err = s.GuildMemberRoleAdd(guildID, i.User.ID, configRole.ID)
			if err != nil {
				s.ChannelMessage(ch.ID, fmt.Sprintf("I was unable to assign the role %q to you, please contact a moderator", role.Name))
				log.Printf("Error assigning role %q\n", err)
				continue
			}

			_, err = s.ChannelMessageSend(conf.RoleManagement.RoleAdminChannelID, fmt.Sprintf("The role <@&%s> was automatically assigned to <@%s>", configRole.ID, i.User.ID))
			if err != nil {
				log.Printf("Error sending role request message to admin channel %q\n", err)
			}
			continue
		}

		if role.Name == "Docent" {
			s.ChannelMessageSend(ch.ID, "Not already working at Thomas More? We're hiring! http://werkenbij.thomasmore.be/")
		}

		s.ChannelMessageSendComplex(conf.RoleManagement.RoleAdminChannelID, &discordgo.MessageSend{
			Content: fmt.Sprintf("<@%s> wants role <@&%s>", i.User.ID, role.ID),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "Add Role",
							Style: discordgo.SuccessButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "➕",
							},
							CustomID: fmt.Sprintf("roleresponse--add--%s--%s", role.ID, i.User.ID),
						},
						discordgo.Button{
							Label: "Replace role of type",
							Style: discordgo.SecondaryButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "🔄",
							},
							CustomID: fmt.Sprintf("roleresponse--replace--%s--%s", role.ID, i.User.ID),
						},
						discordgo.Button{
							Label: "Deny",
							Style: discordgo.DangerButton,
							Emoji: discordgo.ComponentEmoji{
								Name: "❌",
							},
							CustomID: fmt.Sprintf("roleresponse--deny--%s--%s", role.ID, i.User.ID),
						},
					},
				},
			},
		})

	}
}

func (m *MemberCommands) handleRolePermissionResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	data := strings.Split(i.MessageComponentData().CustomID, "--")
	if len(data) < 4 {
		return // not valid ID
	}
	permType := data[1]
	roleID := data[2]
	userID := data[3]

	conf, err := m.db.ConfigForGuild(i.GuildID)
	if err != nil || conf == nil {
		return // no guild data
	}

	if i.ChannelID != conf.RoleManagement.RoleAdminChannelID {
		return
	}

	dm, err := s.UserChannelCreate(userID)
	if err != nil {
		return
	}

	var role *discordgo.Role
	guildRoles, err := s.GuildRoles(i.GuildID)
	if err != nil {
		log.Println("error getting guild roles", err)
		return
	}
	for _, gr := range guildRoles {
		if gr.ID == roleID {
			role = gr
			break
		}
	}

	if role == nil {
		return
	}

	if permType == "deny" {
		msg := fmt.Sprintf("<@%s> denied the request from <@%s> for role <@&%s>", i.Member.User.ID, userID, role.ID)
		s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel:    i.Message.ChannelID,
			ID:         i.Message.ID,
			Content:    &msg,
			Components: []discordgo.MessageComponent{},
		})
		s.ChannelMessageSend(dm.ID, fmt.Sprintf("I'm sorry, your request for role %q has been denied.", role.Name))
		return
	}

	// remove default role
	user, err := s.GuildMember(i.GuildID, userID)
	if err != nil {
		log.Println("error getting user", err)
		return
	}
	if conf.RoleManagement.DefaultRole != "" {
		if hasRole(user, conf.RoleManagement.DefaultRole) {
			s.GuildMemberRoleRemove(i.GuildID, userID, conf.RoleManagement.DefaultRole)
		}
	}

	if permType == "replace" {
		var currentRoleSet db.RoleSet

		for _, rs := range conf.RoleManagement.RoleSets {
			for _, r := range rs.Roles {
				if r.ID == roleID {
					currentRoleSet = rs
					break
				}
			}
		}

		for _, role := range currentRoleSet.Roles {
			if hasRole(user, role.ID) {
				s.GuildMemberRoleRemove(i.GuildID, userID, role.ID)
			}
		}
	}

	err = s.GuildMemberRoleAdd(i.GuildID, userID, roleID)
	if err != nil {
		msg := fmt.Sprintf("Error assigning role %q\n", err)
		s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			Channel: i.Message.ChannelID,
			ID:      i.Message.ID,
			Content: &msg,
		})
		log.Printf("Error assigning role %q\n", err)
		return
	}

	msg := fmt.Sprintf("<@%s> assigned <@&%s> role for <@%s>", i.Member.User.ID, roleID, userID)
	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    i.Message.ChannelID,
		ID:         i.Message.ID,
		Content:    &msg,
		Components: []discordgo.MessageComponent{},
	})
	if err != nil {
		log.Println("error responding to interaction", err)
		s.ChannelMessageSend(i.ChannelID, fmt.Sprintf("<@%s> assigned <@&%s> role for <@%s> (and interaction response failed, sad)", i.Member.User.ID, roleID, userID))
		return
	}

	s.ChannelMessageSend(dm.ID, fmt.Sprintf("Good news! Your request for role %q has been approved!", role.Name))
}

func hasRole(user *discordgo.Member, roleID string) bool {
	for _, r := range user.Roles {
		if r == roleID {
			return true
		}
	}
	return false
}
