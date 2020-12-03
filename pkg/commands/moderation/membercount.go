package moderation

import (
	"fmt"
	"strconv"

	"github.com/itfactory-tm/thomas-bot/pkg/embed"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/sudo"
)

func (m *ModerationCommands) membercount(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if !sudo.IsBotDev(msg.Author.ID) {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("%s is not in the sudoers file. This incident will be reported.", msg.Author.ID))
		return
	}

	// Get the guild status
	g, err := s.State.Guild(msg.GuildID)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error getting guild: %v", err))
		return
	}

	//put the roles into a map and count how many users have that role
	roleMap := make(map[string]int)
	for _, member := range g.Members {
		for _, role := range member.Roles {
			if _, exists := roleMap[role]; !exists {
				roleMap[role] = 0
			}
			roleMap[role]++
		}
	}

	//Create embed
	embed := embed.NewEmbed()
	embed.SetTitle("Membercount")
	embed.AddField("Totaal", strconv.Itoa(g.MemberCount))

	////Print to embed if the role has users
	for _, role := range g.Roles {
		userCount := roleMap[role.ID]
		if userCount > 1 {
			embed.AddField(role.Name, strconv.Itoa(userCount))
		}
	}

	embed.InlineAllFields()
	_, err = s.ChannelMessageSendEmbed(msg.ChannelID, embed.MessageEmbed)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error sending embed message: %v", err))
		return
	}
}
