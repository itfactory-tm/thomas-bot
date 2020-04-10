package main

import (
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/itfactory-tm/thomas-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Token                    string
	Prefix                   string `default:"tm"`
	GiphyKey                 string
	TwitterEnabled           bool   `envconfig:"TWITTER_ENABLED"`
	TwitterConsumerKey       string `envconfig:"TWITTER_CONSUMER_KEY"`
	TwitterConsumerSecret    string `envconfig:"TWITTER_CONSUMER_SECRET"`
	TwitterAccessToken       string `envconfig:"TWITTER_ACCESS_TOKEN"`
	TwitterAccessTokenSecret string `envconfig:"TWITTER_ACCESS_TOKEN_SECRET"`
}

var c config
var handlers = map[string]command.Command{}
var commandRegex *regexp.Regexp

func main() {
	err := envconfig.Process("thomasbot", &c)
	if err != nil {
		log.Fatal(err)
	}
	if c.Token == "" {
		log.Fatal("No token specified")
	}

	commandRegex = regexp.MustCompile(c.Prefix + `!(\w*)\b`)

	dg, err := discordgo.New("Bot " + c.Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	// Register handlers
	dg.AddHandler(onMessage)
	dg.AddHandler(onReactionAdd)

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	// TODO: add connection error handlers

	dg.UpdateStreamingStatus(0, "Thomas Bot", "https://github.com/itfactory-tm/thomas-bot")

	go postHashtagTweets(dg)

	log.Println("Thomas Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if commandRegex.MatchString(m.Content) {
		if c, exists := handlers[commandRegex.FindStringSubmatch(m.Content)[1]]; exists {
			c.Handler(s, m)
		}
	}
}

func onReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	handleHelpReaction(s, r)
}

func registerCommand(c command.Command) {
	handlers[c.Name] = c
	if _, exists := helpData[c.Category]; !exists {
		if !c.Hidden {
			helpData[c.Category] = map[string]command.Command{}
		}
	}
	helpData[c.Category][c.Name] = c
}

func registerCommandDEPRECATED(name, category, helpText string, fn func(*discordgo.Session, *discordgo.MessageCreate)) {
	registerCommand(command.Command{
		Name:        name,
		Category:    command.StringToCategory(category),
		Description: helpText,
		Hidden:      false,
		Handler:     fn,
	})
}
