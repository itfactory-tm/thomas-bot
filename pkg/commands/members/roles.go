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

		_, err = s.ChannelMessageSendComplex(ch.ID, &discordgo.MessageSend{
			Content: rs.Message,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							MinValues:   1,
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
			return
		}

		for _, mr := range member.Roles {
			if mr == val {
				s.ChannelMessageSend(ch.ID, fmt.Sprintf("Oopsie! You already have the role %q, no worries I will not re-request it!", role.Name))
				return
			}
		}

		s.ChannelMessageSend(ch.ID, fmt.Sprintf("Thank you! I have asked our moderators for permissions to assign the role %q", role.Name))
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
		s.InteractionRespond(i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("<@%s> has denied to give <@%s> the role <@&%s>", i.Member.User.ID, userID, roleID),
				},
			})
		s.ChannelMessageSend(dm.ID, fmt.Sprintf("I'm sorry, your request for role %q has been denied.", role.Name))
		return
	}

	// remove default role
	if conf.RoleManagement.DefaultRole != "" {
		s.GuildMemberRoleRemove(i.GuildID, userID, conf.RoleManagement.DefaultRole)
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
			s.GuildMemberRoleRemove(i.GuildID, userID, role.ID)
		}
	}

	err = s.GuildMemberRoleAdd(i.GuildID, userID, roleID)
	if err != nil {
		s.InteractionRespond(i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Error assigning role %q\n", err),
				},
			})
		log.Printf("Error assigning role %q\n", err)
		return
	}

	s.InteractionRespond(i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("<@%s> assigned <@&%s> role for <@%s>", i.Member.User.ID, roleID, userID),
			},
		})

	s.ChannelMessageSend(dm.ID, fmt.Sprintf("Good news! your request for role %q has been approved!", role.Name))
}
