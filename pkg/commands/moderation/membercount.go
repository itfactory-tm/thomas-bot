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

	//Get updated memberlist
	memberList, err := s.GuildMembers(msg.GuildID, "", 1000)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error getting member list: %v", err))
		return
	}
	for len(memberList) < g.MemberCount {
		members, err := s.GuildMembers(msg.GuildID, memberList[len(memberList)-1].GuildID, 1000)
		if err != nil {
			s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error getting member list: %v", err))
			return
		}
		memberList = append(memberList, members...)
	}

	//Put the roles into a map and count how many users have that role
	roleMap := make(map[string]int)
	for _, member := range memberList {
		for _, role := range member.Roles {
			if _, exists := roleMap[role]; !exists {
				roleMap[role] = 0
			}
			roleMap[role]++
		}
	}

	//Create embed
	embedmessage := embed.NewEmbed()
	embedmessage.SetTitle("Membercount")
	embedmessage.SetThumbnail(g.IconURL())
	embedmessage.SetFooter(fmt.Sprintf("Guild total %v; Members counted: %v", g.MemberCount, len(memberList)))
	embedmessage.AddField("Total", strconv.Itoa(g.MemberCount))

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
}

func sendEmbed(s *discordgo.Session, msg *discordgo.MessageCreate, embedmessage *embed.Embed) {
	embedmessage.InlineAllFields()
	_, err := s.ChannelMessageSendEmbed(msg.ChannelID, embedmessage.MessageEmbed)
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error sending embed message: %v", err))
		return
	}
}
