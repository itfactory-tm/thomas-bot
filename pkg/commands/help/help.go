package help

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
	"github.com/itfactory-tm/thomas-bot/pkg/embed"
)

// TODO: make commands able to read config
const prefix = "tm"

// HelpCommand contains the tm!hello command
type HelpCommand struct {
	// map[category][]map[command]description
	helpData map[command.Category]map[string]command.Command

	server command.Server
}

// NewHelpCommand gives a new HelpCommand
func NewHelpCommand() *HelpCommand {
	return &HelpCommand{}
}

// Register registers the handlers
func (h *HelpCommand) Register(registry command.Registry, server command.Server) {
	registry.RegisterMessageCreateHandler("help", h.sayHelp)
	registry.RegisterMessageReactionAddHandler(h.handleHelpReaction)

	h.server = server
}

// PopulateHelpData populates the internal help data in memory
func (h *HelpCommand) PopulateHelpData() {
	commands := h.server.GetAllCommandInfos()
	h.helpData = map[command.Category]map[string]command.Command{}

	for _, c := range commands {
		if _, exists := h.helpData[c.Category]; !exists {
			h.helpData[c.Category] = map[string]command.Command{}
		}
		if !c.Hidden {
			h.helpData[c.Category][c.Name] = c
		}
	}
}

func (h *HelpCommand) sayHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	if h.helpData == nil {
		h.PopulateHelpData()
	}
	embed := h.helpMenu()

	ch, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Cannot DM user")
		return
	}

	if ch.ID != m.ChannelID && m.Message.Content == fmt.Sprintf("%s!help", prefix) {
		s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
	}

	ints := []int{}
	for category := range h.helpData {
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
			log.Printf("Error adding help emoji: %q\n", err)
		}
	}
}

// Info return the commands in this package
func (h *HelpCommand) Info() []command.Command {
	return []command.Command{
		command.Command{
			Name:        "help",
			Category:    command.CategoryAlgemeen,
			Description: "List of all commands (you are here)",
			Hidden:      false,
		},
	}
}

func (h *HelpCommand) helpMenu() *embed.Embed {
	embed := embed.NewEmbed()
	embed.SetTitle("Help")

	categories := []string{}
	for category := range h.helpData {
		categories = append(categories, fmt.Sprintf("* `%d`:\t%s", category, command.CategoryToString(category)))
	}
	sort.Strings(categories)
	embed.AddField("choose a category", strings.Join(categories, "\n"))

	return embed
}

func (h *HelpCommand) handleHelpReaction(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Println("Cannot get message of reaction", r.ChannelID)
		return
	}

	if message.Author.ID != s.State.User.ID {
		return // not the bot user
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
	if data, ok := h.helpData[category]; ok {
		keys := []string{}
		for key := range data {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			messages = append(messages, fmt.Sprintf("* `%s`: %s", data[key].Name, data[key].Description))
		}

		embed := h.helpMenu()
		embed.AddField(command.CategoryToString(category), strings.Join(messages, "\n"))

		_, err := s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, embed.MessageEmbed)
		if err != nil {
			log.Printf("Error editing help: %v", err)
		}
	}
}
