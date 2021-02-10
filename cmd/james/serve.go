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

	"github.com/itfactory-tm/thomas-bot/pkg/db"

	"github.com/itfactory-tm/thomas-bot/pkg/commands/help"
	"github.com/itfactory-tm/thomas-bot/pkg/commands/hive"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/commands/game"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"

	"github.com/itfactory-tm/thomas-bot/pkg/command"

	discordha "github.com/meyskens/discord-ha"
)

func init() {
	rootCmd.AddCommand(NewServeCmd())
}

type serveCmdOptions struct {
	Token  string
	Prefix string `default:"bob"`

	EtcdEndpoints []string `envconfig:"ETCD_ENDPOINTS"`
	MongoDBURL    string   `envconfig:"MONGODB_URL"`
	MongoDBDB     string   `envconfig:"MONGODB_DB"`
	ConfigPath    string   `default:"./config.json" envconfig:"CONFIG"`

	commandRegex *regexp.Regexp
	dg           *discordgo.Session
	ha           discordha.HA
	handlers     []command.Interface
	db           db.Database

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
		Long:    `This connects to Discord and handle all game related events`,
		RunE:    s.RunE,
		PreRunE: s.Validate,
	}

	// TODO: switch to viper
	err := envconfig.Process("thomasbob", &s)
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

	dg, err := discordgo.New("Bot " + s.Token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	dg.UpdateStreamingStatus(0, fmt.Sprintf("Thomas Bob rev. %s", revision), "")
	haLogger := log.New(os.Stdout, "discordha: ", log.Ldate|log.Ltime)
	s.ha, err = discordha.New(&discordha.Config{
		Session:       dg,
		HA:            len(s.EtcdEndpoints) > 0,
		EtcdEndpoints: s.EtcdEndpoints,
		Context:       ctx,
		Log:           *haLogger,
	})
	if err != nil {
		return fmt.Errorf("error creating Discord HA: %w", err)
	}

	err = dg.Open()
	if err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}
	defer dg.Close()

	if s.MongoDBDB != "" {
		s.db, err = db.NewMongoDB(s.MongoDBURL, s.MongoDBDB)
		if err != nil {
			return err
		}
	} else {
		// local fallback
		s.db, err = db.NewLocalDB(s.ConfigPath)
		if err != nil {
			return err
		}
	}

	s.RegisterHandlers()

	s.ha.AddHandler(s.onMessage)
	s.ha.AddHandler(s.onMessageReactionAdd)

	log.Println("Thomas Bob is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	s.ha.Stop()

	return nil
}

func (s *serveCmdOptions) RegisterHandlers() {
	s.handlers = []command.Interface{
		game.NewUserCommand(),
		game.NewMuteCommand(),
		help.NewHelpCommand(),
		hive.NewHiveCommandForBob(s.db),
	}

	for _, handler := range s.handlers {
		handler.Register(s, s)
	}
}

func (s *serveCmdOptions) onMessage(sess *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == sess.State.User.ID {
		return
	}

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
func (s *serveCmdOptions) onMessageReactionAdd(sess *discordgo.Session, m *discordgo.MessageReactionAdd) {
	// Ignore all reactions created by the bot itself
	if m.UserID == sess.State.User.ID {
		return
	}

	for _, handler := range s.onMessageReactionAddHandler {
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

func (s *serveCmdOptions) GetDiscordHA() discordha.HA {
	return s.ha
}

func (s *serveCmdOptions) GetAllCommandInfos() []command.Command {
	out := []command.Command(nil)
	for _, handler := range s.handlers {
		out = append(out, handler.Info()...)
	}

	return out
}
