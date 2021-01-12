package hive

import (
	"fmt"
	"log"
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/embed"
	"github.com/itfactory-tm/thomas-bot/pkg/sudo"
)

var verifyRegex = regexp.MustCompile(`!verify ([0-9]*) (.*)$`)

const infoDeskID = "794973874634752040" //TODO: make me configurable

// SayVerify handles the tm!verify command
func (h *HiveCommand) SayVerify(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !sudo.IsAdmin(m.Author.ID) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is not in the sudoers file. This incident will be reported.", m.Author.ID))
		return
	}

	matched := verifyRegex.FindStringSubmatch(m.Message.Content)

	if len(matched) < 3 {
		s.ChannelMessageSend(m.ChannelID, "invalid syntax, needs to be ID + description")
		return
	}

	channel, err := s.Channel(matched[1])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	e := embed.NewEmbed()
	e.SetTitle("Hive Channel")
	e.AddField("name", channel.Name)
	e.AddField("description", matched[2])
	e.AddField("id", channel.ID)

	msg, err := s.ChannelMessageSendEmbed(infoDeskID, e.MessageEmbed)
	if err != nil {
		log.Println(err)
	}
	s.MessageReactionAdd(infoDeskID, msg.ID, "👋")
}
