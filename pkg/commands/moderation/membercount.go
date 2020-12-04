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

	//Put the roles into a map and count how many users have that role
	roleMap := make(map[string]int)
	//Map for debugging => array for later? doesn't need values to check doubles
	memberMap := make(map[string]int)
	for _, member := range g.Members {
		if _, exists := memberMap[member.User.ID]; !exists {
			memberMap[member.User.ID] = 0
			for _, role := range member.Roles {
				if _, exists := roleMap[role]; !exists {
					roleMap[role] = 0
				}
				roleMap[role]++
			}
		}
		memberMap[member.User.ID]++
	}

	//Create embed
	embedmessage := embed.NewEmbed()
	embedmessage.SetTitle("Membercount")
	embedmessage.AddField("Totaal", strconv.Itoa(g.MemberCount))
	embedmessage.AddField("Totaal len", strconv.Itoa(len(g.Members)))

	//Print to embed if the role has more than 1 user (filters bot roles)
	for _, role := range g.Roles {
		userCount := roleMap[role.ID]
		if userCount > 1 {
			embedmessage.AddField(role.Name, strconv.Itoa(userCount))
			//Discord embeds only allow 25 fields => make new embed
			if len(embedmessage.Fields) >= 25 {
				sendEmbed(s, msg, embedmessage)
				embedmessage = embed.NewEmbed()
				embedmessage.SetTitle("Membercount")
			}
		}
	}

	//Prevent sending empty embed
	if len(embedmessage.Fields) != 0 {
		sendEmbed(s, msg, embedmessage)
	}

	endString := "Double users:\n"
	for key, value := range memberMap {
		if memberMap[key] > 1 {
			endString = endString + fmt.Sprintf("User: <@%v> => %v \n", key, value)
		}
	}
	s.ChannelMessageSend(msg.GuildID, endString)
}

func sendEmbed(s *discordgo.Session, msg *discordgo.MessageCreate, embedmessage *embed.Embed) {
	embedmessage.InlineAllFields()
	_, err := s.ChannelMessageSendEmbed(msg.ChannelID, embedmessage.MessageEmbed)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error sending embed message: %v", err))
		return
	}
}
