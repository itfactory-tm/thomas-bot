package honeypot

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
	"github.com/itfactory-tm/thomas-bot/pkg/db"
	"github.com/itfactory-tm/thomas-bot/pkg/sudo"
)

var dieMessages = []string{
	"\"Dat is een buis\" - Freddi Fish",
	"ga toch eventjes je cyber security les heralen",
	"de instructies waren duidelijk volgende keer beter!",
	"het KULoket is geinformeerd...",
	"Gunther is trots op je... sorry ik heb trots fout begrepen...",
	"To the shadow realm you go",
	"Niet op vreemde linkjes klikken he",
}

// HoneypotCommands contains the vegan honeypot module
type HoneypotCommands struct {
	server command.Server
	db     db.Database
}

// NewHoneypotCommands gives a new HoneypotCommands
func NewHoneypotCommands(conn db.Database) *HoneypotCommands {
	return &HoneypotCommands{
		db: conn,
	}
}

// Info return the commands in this package
func (m *HoneypotCommands) Info() []command.Command {
	return []command.Command{}
}

// Register registers the handlers
func (m *HoneypotCommands) Register(registry command.Registry, server command.Server) {
	registry.RegisterMessageCreateHandler("", m.checkMessageCreateAsync)

	m.server = server
}

// InstallSlashCommands registers the slash commands
func (m *HoneypotCommands) InstallSlashCommands(session *discordgo.Session) error {
	return nil
}

func (m *HoneypotCommands) checkMessageCreateAsync(s *discordgo.Session, msg *discordgo.MessageCreate) {
	conf, err := m.db.ConfigForGuild(msg.GuildID)
	if err != nil {
		return
	}

	if conf.HoneypotChannelID == "" || msg.ChannelID != conf.HoneypotChannelID {
		return
	}

	if sudo.IsAdmin(msg.Author.ID) || msg.Author.Bot {
		return
	}

	// somebody had bitten the dust
	err = s.GuildBanCreateWithReason(msg.GuildID, msg.Author.ID, "Has posted in honeypot", 1)
	if err != nil {
		log.Panicln(err)
		return
	}

	randomInt := rand.Intn(len(dieMessages))

	s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Sorry %s, %s", msg.Author.Username, dieMessages[randomInt]))

	time.Sleep(10 * time.Second)

	// soft ban you can come back once more
	s.GuildBanDelete(msg.GuildID, msg.Author.ID)
}
