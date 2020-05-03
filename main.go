package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/discordha"

	"github.com/itfactory-tm/thomas-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
)

// to be overwritten in build
var revision = "dev"

type config struct {
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
}

var c config
var handlers = map[string]command.Command{}
var commandRegex *regexp.Regexp
var dg *discordgo.Session
var ha *discordha.HA

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := envconfig.Process("thomasbot", &c)
	if err != nil {
		log.Fatalf("Error processing envvars: %q\n", err)
	}
	if c.Token == "" {
		log.Fatal("No token specified")
	}

	commandRegex = regexp.MustCompile(c.Prefix + `!(\w*)\b`)

	dg, err = discordgo.New("Bot " + c.Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	ha, err = discordha.New(discordha.Config{
		Session:       dg,
		HA:            len(c.EtcdEndpoints) > 0,
		EtcdEndpoints: c.EtcdEndpoints,
		Context:       ctx,
	})
	if err != nil {
		log.Fatal("error creating Discord HA,", err)
	}

	// Register handlers
	dg.AddHandler(onMessage)
	dg.AddHandler(onMessageEdit)
	dg.AddHandler(onReactionAdd)
	dg.AddHandler(onNewMember)

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	defer dg.Close()

	dg.UpdateStreamingStatus(0, fmt.Sprintf("Thomas Bot rev. %s", revision), "")

	go postHashtagTweets(dg, ctx)
	go serve()

	log.Println("Thomas Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if ok, err := ha.Lock(m); !ok {
		if err != nil {
			log.Println(err)
		}
		return
	}
	defer ha.Unlock(m)
	go checkMessage(s, m)

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

func onMessageEdit(s *discordgo.Session, u *discordgo.MessageUpdate) {
	if ok, err := ha.Lock(u); !ok {
		if err != nil {
			log.Printf("Error locking message edit: %q\n", err)
		}
		return
	}
	defer ha.Unlock(u)
	m := &discordgo.MessageCreate{
		u.Message,
	}

	go checkMessage(s, m)
}

func onReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}
	if ok, err := ha.Lock(r); !ok {
		if err != nil {
			log.Printf("Error locking reaction add: %q\n", err)
		}
		return
	}
	defer ha.Unlock(r)
	go checkReaction(s, r)
	handleHelpReaction(s, r)
}

func onNewMember(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	if ok, err := ha.Lock(g); !ok {
		if err != nil {
			log.Printf("Error lockin on new memner: %q\n", err)
		}
		return
	}
	defer ha.Unlock(g)
	if g.GuildID != itfDiscord {
		return
	}
	err := s.GuildMemberRoleAdd(g.GuildID, g.Member.User.ID, "687568536356257890") // gast role
	if err != nil {
		log.Printf("Cannot set role for user %s: %q\n", g.Member.User.ID, err)
	}

	s.ChannelMessageSend(itfWelcome, fmt.Sprintf("Welkom <@%s> op de **IT Factory Official** Discord server. Je wordt automatisch toegevoegd als **gast**. Indien je student of alumnus bent en  toegang wil tot de studenten- of alumnikanalen, gelieve dan een van de moderatoren te contacteren om de juiste rol te krijgen. Indien je graag informatie hebt over onze opleiding, neem dan gerust een kijkje op ons <#693046715665874944>.", g.User.ID))

	c, err := s.UserChannelCreate(g.Member.User.ID)
	if err != nil {
		log.Printf("Cannot DM user %s\n", g.Member.User.ID)
		return
	}

	s.ChannelMessageSend(c.ID, fmt.Sprintf("Hallo %s", g.User.Username))
	time.Sleep(time.Second)
	s.ChannelMessageSend(c.ID, "Welkom op de ITFactory Discord!")
	time.Sleep(time.Second)
	s.ChannelMessageSend(c.ID, "Mijn naam is Thomas Bot, ik ben een bot die jou kan helpen!")
	time.Sleep(time.Second)
	s.ChannelMessageSend(c.ID, "Nieuw op Discord? Geen probleem hier is een handleiding: https://itf.to/discord-help")
	embed := NewEmbed()
	embed.SetImage("https://static.eyskens.me/thomas-bot/opendeurdag-1.png")
	embed.SetURL("https://itf.to/discord-help")
	s.ChannelMessageSendEmbed(c.ID, embed.MessageEmbed)

	time.Sleep(time.Second)
	s.ChannelMessageSend(c.ID, "Heb je hulp nodig zeg dan tm!help")
	time.Sleep(time.Second)
	s.ChannelMessageSend(c.ID, "Let op, ik kan enkel antwoorden op commandos die starten met `tm!` niet op gewone berichten.")
	time.Sleep(5 * time.Second)
	s.ChannelMessageSend(c.ID, "Klaar voor de virtuele opendeurdag? Je kan best starten in <#693046715665874944>")
}

func registerCommand(c command.Command) {
	handlers[c.Name] = c
	if _, exists := helpData[c.Category]; !exists {
		helpData[c.Category] = map[string]command.Command{}
	}
	if !c.Hidden {
		helpData[c.Category][c.Name] = c
	}
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
