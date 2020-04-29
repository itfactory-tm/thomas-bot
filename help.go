package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/itfactory-tm/thomas-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
)

// map[category][]map[command]description
var helpData = map[command.Category]map[string]command.Command{}

func init() {
	registerCommand(command.Command{
		Name:        "help",
		Category:    command.CategoryAlgemeen,
		Description: "Lijst van alle commandos (u bent hier)",
		Hidden:      false,
		Handler:     sayHelp,
	})
}

func sayHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := helpMenu()

	ch, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Cannot DM user")
		return
	}

	if ch.ID != m.ChannelID && m.Message.Content == fmt.Sprintf("%s!help", c.Prefix) {
		s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
	}

	ints := []int{}
	for category := range helpData {
		ints = append(ints, int(category))
	}

	em, err := s.ChannelMessageSendEmbed(ch.ID, embed.MessageEmbed)
	if err != nil {
		log.Printf("Cannot send embed to %q \n", ch.ID)
		return
	}

	sort.Ints(ints)
	for _, i := range ints {
		err := s.MessageReactionAdd(ch.ID, em.ID, intToEmoji(i))
		if err != nil {
			log.Println(err)
		}
	}
}

func handleHelpReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		s.ChannelMessageSend(r.ChannelID, "Cannot get message of reaction")
		return
	}

	if message.Author.ID != s.State.User.ID {
		return // not the bot user
	}

	if r.UserID == s.State.User.ID {
		return // the bot itself reacted
	}

	if len(message.Embeds) < 1 {
		return // not the help message
	}

	if message.Embeds[0].Title != "Help" {
		return // not the help message
	}

	i := emojiToInt(r.Emoji.MessageFormat())
	if i < 0 {
		return // not a valid emoji
	}

	category := command.Category(i)

	messages := []string{}
	if data, ok := helpData[category]; ok {
		keys := []string{}
		for key := range data {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			messages = append(messages, fmt.Sprintf("* `%s`: %s", data[key].Name, data[key].Description))
		}

		embed := helpMenu()
		embed.AddField(command.CategoryToString(category), strings.Join(messages, "\n"))

		s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, embed.MessageEmbed)
	}
}

func helpMenu() *Embed {
	embed := NewEmbed()
	embed.SetTitle("Help")

	categories := []string{}
	for category := range helpData {
		categories = append(categories, fmt.Sprintf("* `%d`:\t%s", category, command.CategoryToString(category)))
	}
	sort.Strings(categories)
	embed.AddField("kies een categorie", strings.Join(categories, "\n"))

	return embed
}

func intToEmoji(i int) string {
	switch i {
	case 0:
		return "0️⃣"
	case 1:
		return "1️⃣"
	case 2:
		return "2️⃣"
	case 3:
		return "3️⃣"
	case 4:
		return "4️⃣"
	case 5:
		return "5️⃣"
	case 6:
		return "6️⃣"
	case 7:
		return "7️⃣"
	case 8:
		return "8️⃣"
	case 9:
		return "9️⃣"
	}

	return ""
}

func emojiToInt(i string) int {
	switch i {
	case "0️⃣":
		return 0
	case "1️⃣":
		return 1
	case "2️⃣":
		return 2
	case "3️⃣":
		return 3
	case "4️⃣":
		return 4
	case "5️⃣":
		return 5
	case "6️⃣":
		return 6
	case "7️⃣":
		return 7
	case "8️⃣":
		return 8
	case "9️⃣":
		return 9
	}

	return -1
}
