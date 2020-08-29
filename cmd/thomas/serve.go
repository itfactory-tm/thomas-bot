package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/itfactory-tm/thomas-bot/pkg/commands/members"

	"github.com/itfactory-tm/thomas-bot/pkg/command"

	"github.com/itfactory-tm/thomas-bot/pkg/commands/hello"

	"github.com/kelseyhightower/envconfig"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/discordha"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewServeCmd())
}

type serveCmdOptions struct {
	Token                    string
	Prefix                   string `default:"tm"`
	GiphyKey                 string
	TwitterEnabled           bool     `envconfig:"TWITTER_ENABLED"`
	TwitterConsumerKey       string   `envconfig:"TWITTER_CONSUMER_KEY"`
	TwitterConsumerSecret    string   `envconfig:"TWITTER_CONSUMER_SECRET"`
	TwitterAccessToken       string   `envconfig:"TWITTER_ACCESS_TOKEN"`
	TwitterAccessTokenSecret string   `envconfig:"TWITTER_ACCESS_TOKEN_SECRET"`
	HCaptchaSiteKey          string   `envconfig:"HCAPTCHA_SITE_KEY"`
	HCaptchaSiteSecret       string   `envconfig:"HCAPTCHA_SITE_SECRET"`
	BindAddr                 string   `default:":8080" envconfig:"BIND_ADDR"`
	EtcdEndpoints            []string `envconfig:"ETCD_ENDPOINTS"`

	commandRegex *regexp.Regexp
	dg           *discordgo.Session
	ha           *discordha.HA

	onMessageCreateHandlers     map[string][]func(*discordgo.Session, *discordgo.MessageCreate)
	onMessageEditHandlers       map[string][]func(*discordgo.Session, *discordgo.MessageUpdate)
	onMessageReactionAddHandler []func(*discordgo.Session, *discordgo.MessageReactionAdd)
	onGuildMemberAddHandler     []func(*discordgo.Session, *discordgo.GuildMemberAdd)
}

// NewServeCmd generates the `serve` command
func NewServeCmd() *cobra.Command {
	s := serveCmdOptions{}
	c := &cobra.Command{
		Use:     "serve",
		Short:   "Run the server",
		Long:    `This connects to Discord and handle all events, it will also serve several HTTP pages`,
		RunE:    s.RunE,
		PreRunE: s.Validate,
	}

	// TODO: switch to viper
	err := envconfig.Process("thomasbot", &s)
	if err != nil {
		log.Fatalf("Error processing envvars: %q\n", err)
	}

	return c
}

func (s *serveCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	if s.Token == "" {
		return errors.New("No token specified")
	}

	s.RegisterHandlers()

	return nil
}

func (s *serveCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.commandRegex = regexp.MustCompile(s.Prefix + `!(\w*)\b`)

	http := NewServeHTTPCmd()
	err := http.PreRunE(cmd, args)
	if err != nil {
		return err
	}

	go http.RunE(cmd, args)

	dg, err := discordgo.New("Bot " + s.Token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	s.ha, err = discordha.New(discordha.Config{
		Session:       dg,
		HA:            len(s.EtcdEndpoints) > 0,
		EtcdEndpoints: s.EtcdEndpoints,
		Context:       ctx,
	})
	if err != nil {
		return fmt.Errorf("error creating Discord HA: %w", err)
	}

	// TODO: Register handlers
	dg.AddHandler(s.onMessage)
	dg.AddHandler(s.onMessageUpdate)
	dg.AddHandler(s.onMessageReactionAdd)
	dg.AddHandler(s.onGuildMemberAdd)

	err = dg.Open()
	if err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}
	defer dg.Close()

	dg.UpdateStreamingStatus(0, fmt.Sprintf("Thomas Bot rev. %s", revision), "")

	// TODO: go postHashtagTweets(ctx, dg)
	// TODO: go serve()

	log.Println("Thomas Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	return nil
}

func (s *serveCmdOptions) RegisterHandlers() {
	handlers := []command.Interface{
		hello.NewHelloCommand(),
		members.NewMemberCommand(),
	}

	for _, handler := range handlers {
		handler.Register(s)
	}
}

func (s *serveCmdOptions) onMessage(sess *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == sess.State.User.ID {
		return
	}

	if ok, err := s.ha.Lock(m); !ok {
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer s.ha.Unlock(m)
	matchedHandlers := []func(*discordgo.Session, *discordgo.MessageCreate){}
	if handlers, hasHandlers := s.onMessageCreateHandlers[""]; hasHandlers {
		matchedHandlers = append(matchedHandlers, handlers...)
	}

	if s.commandRegex.MatchString(m.Content) {
		if handlers, hasHandlers := s.onMessageCreateHandlers[s.commandRegex.FindStringSubmatch(m.Content)[1]]; hasHandlers {
			matchedHandlers = append(matchedHandlers, handlers...)
		}
	}

	for _, handler := range matchedHandlers {
		handler(sess, m)
	}
}

func (s *serveCmdOptions) onMessageUpdate(sess *discordgo.Session, m *discordgo.MessageUpdate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == sess.State.User.ID {
		return
	}

	if ok, err := s.ha.Lock(m); !ok {
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer s.ha.Unlock(m)
	matchedHandlers := []func(*discordgo.Session, *discordgo.MessageUpdate){}
	if handlers, hasHandlers := s.onMessageEditHandlers[""]; hasHandlers {
		matchedHandlers = append(matchedHandlers, handlers...)
	}

	if s.commandRegex.MatchString(m.Content) {
		if handlers, hasHandlers := s.onMessageEditHandlers[s.commandRegex.FindStringSubmatch(m.Content)[1]]; hasHandlers {
			matchedHandlers = append(matchedHandlers, handlers...)
		}
	}

	for _, handler := range matchedHandlers {
		handler(sess, m)
	}
}

func (s *serveCmdOptions) onMessageReactionAdd(sess *discordgo.Session, m *discordgo.MessageReactionAdd) {
	// Ignore all reactions created by the bot itself
	if m.UserID == sess.State.User.ID {
		return
	}

	if ok, err := s.ha.Lock(m); !ok {
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer s.ha.Unlock(m)

	for _, handler := range s.onMessageReactionAddHandler {
		handler(sess, m)
	}
}

func (s *serveCmdOptions) onGuildMemberAdd(sess *discordgo.Session, m *discordgo.GuildMemberAdd) {
	if ok, err := s.ha.Lock(m); !ok {
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer s.ha.Unlock(m)

	for _, handler := range s.onGuildMemberAddHandler {
		handler(sess, m)
	}
}

func (s *serveCmdOptions) RegisterMessageCreateHandler(command string, fn func(*discordgo.Session, *discordgo.MessageCreate)) {
	if s.onMessageCreateHandlers == nil {
		s.onMessageCreateHandlers = map[string][]func(*discordgo.Session, *discordgo.MessageCreate){}
	}

	if _, exists := s.onMessageCreateHandlers[command]; !exists {
		s.onMessageCreateHandlers[command] = []func(*discordgo.Session, *discordgo.MessageCreate){}
	}

	s.onMessageCreateHandlers[command] = append(s.onMessageCreateHandlers[command], fn)
}
func (s *serveCmdOptions) RegisterMessageEditHandler(command string, fn func(*discordgo.Session, *discordgo.MessageUpdate)) {
	if s.onMessageEditHandlers == nil {
		s.onMessageEditHandlers = map[string][]func(*discordgo.Session, *discordgo.MessageUpdate){}
	}

	if _, exists := s.onMessageEditHandlers[command]; !exists {
		s.onMessageEditHandlers[command] = []func(*discordgo.Session, *discordgo.MessageUpdate){}
	}

	s.onMessageEditHandlers[command] = append(s.onMessageEditHandlers[command], fn)
}
func (s *serveCmdOptions) RegisterMessageReactionAddHandler(fn func(*discordgo.Session, *discordgo.MessageReactionAdd)) {
	if s.onMessageReactionAddHandler == nil {
		s.onMessageReactionAddHandler = []func(*discordgo.Session, *discordgo.MessageReactionAdd){}
	}

	s.onMessageReactionAddHandler = append(s.onMessageReactionAddHandler, fn)
}
func (s *serveCmdOptions) RegisterGuildMemberAddHandler(fn func(*discordgo.Session, *discordgo.GuildMemberAdd)) {
	if s.onGuildMemberAddHandler == nil {
		s.onGuildMemberAddHandler = []func(*discordgo.Session, *discordgo.GuildMemberAdd){}
	}

	s.onGuildMemberAddHandler = append(s.onGuildMemberAddHandler, fn)
}