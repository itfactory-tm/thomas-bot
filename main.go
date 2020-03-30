package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Token string
}

var c config

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

	if strings.Contains(m.Content, "tm!hello") {
		s.ChannelMessageSend(m.ChannelID, "Beep bop boop! Ik ben Thomas Bot, fork me on GitHub!")
	}
}
