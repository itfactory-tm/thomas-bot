package main

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/davecgh/go-spew/spew"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Token string
}

var c config
var handlers = map[string]func(*discordgo.Session, *discordgo.MessageCreate){}
var commandRegex = regexp.MustCompile(`tm!(\w*)\b`)

func main() {
	err := envconfig.Process("thomasbot", &c)
	if c.Token == "" {
		log.Fatal("No token specified")
	}

	dg, err := discordgo.New("Bot " + c.Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	// Register handlers
	dg.AddHandler(onMessage)

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	// TODO: add connection error handlers

	log.Println("Thomas Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	spew.Dump(m.Content, commandRegex.MatchString(m.Content))

	if commandRegex.MatchString(m.Content) {
		if fn, exists := handlers[commandRegex.FindStringSubmatch(m.Content)[1]]; exists {
			fn(s, m)
		}
	}
}

func registerCommand(name string, fn func(*discordgo.Session, *discordgo.MessageCreate)) {
	handlers[name] = fn
}
