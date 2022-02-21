package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/db"

	"github.com/bwmarrin/discordgo"

	"github.com/kelseyhightower/envconfig"

	"github.com/spf13/cobra"
)

const agora = "687565214555570195"

func init() {
	rootCmd.AddCommand(NewCleanCmd())
}

type cleanCmdOptions struct {
	Token string

	MongoDBURL string `envconfig:"MONGODB_URL"`
	MongoDBDB  string `envconfig:"MONGODB_DB"`
	ConfigPath string `default:"./config.json" envconfig:"CONFIG"`

	dg           *discordgo.Session
	shouldRemove map[string]bool
	db           db.Database
}

// NewCleanCmd generates the `clean` command
func NewCleanCmd() *cobra.Command {
	s := cleanCmdOptions{}
	c := &cobra.Command{
		Use:     "clean",
		Short:   "Run the voice channel cleaner",
		Long:    `This is a separate instance cleaning unused channels in The Hive`,
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

func (v *cleanCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	if v.Token == "" {
		return errors.New("No token specified")
	}

	return nil
}

func (v *cleanCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	log.Println("Starting Alf...")

	var err error
	if v.MongoDBDB != "" {
		v.db, err = db.NewMongoDB(v.MongoDBURL, v.MongoDBDB)
		if err != nil {
			return err
		}
	} else {
		// local fallback
		v.db, err = db.NewLocalDB(v.ConfigPath)
		if err != nil {
			return err
		}
	}

	dg, err := discordgo.New("Bot " + v.Token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	dg.StateEnabled = true
	dg.State.TrackVoice = true

	err = dg.Open()
	if err != nil {
		return fmt.Errorf("error opening Discord session: %w", err)
	}

	v.dg = dg

	go func() {
		for {
			time.Sleep(time.Second)
			now := time.Now().UTC()
			_, m, d := now.Date()
			if m == time.March && d == 30 && now.Hour() == 17 && now.Minute() == 15 {
				dg.ChannelMessageSend(agora, "Happy birthday to me")
				time.Sleep(time.Second)
				dg.ChannelMessageSend(agora, "Happy birthday to me")
				time.Sleep(time.Second)
				dg.ChannelMessageSend(agora, "Happy birthday to Thomas Bot")
				time.Sleep(time.Second)
				dg.ChannelMessageSend(agora, "Happy birthday to me")
				time.Sleep(2 * time.Second)
				dg.ChannelMessageSend(agora, "Happy birthday Edward, James, John, Thomas!")
				time.Sleep(time.Second)
				dg.ChannelMessageSend(agora, "Bedankt <@687715371255463972> en <@252083102992695296> en alle contributors om mij te laten draaien :)")
				time.Sleep(time.Second)
				dg.ChannelMessageSend(agora, "Groetjes, Alf de enige (maar stille) Thomas Bot microservice die een besef van tijd heeft")
				return
			}
		}
	}()

	go func() {
		for {
			tz, err := time.LoadLocation("Europe/Brussels")
			if err != nil {
				panic(err)
			}

			dirk, _ := time.Parse("2006-01-02", "2022-06-18")
			time.Sleep(time.Second)
			now := time.Now().In(tz)
			// calculate days till dirk
			days := dirk.Sub(now).Hours() / 24

			if days <= 100 && days >= 0 && now.Hour() == 0 && now.Minute() == 0 && now.Year() == 2022 {
				dg.ChannelMessageSend(agora, fmt.Sprintf("@<177531421152247809> you have %d days left of being 39 years old.", int(days)))
				time.Sleep(2 * time.Minute)
			}

			if days == 0 {
				dg.ChannelMessageSend(agora, "Wait a minute...")
				time.Sleep(time.Minute)
				dg.ChannelMessageSend(agora, "Happy birthday to you")
				time.Sleep(time.Second)
				dg.ChannelMessageSend(agora, "Happy birthday to you")
				time.Sleep(time.Second)
				dg.ChannelMessageSend(agora, "Happy birthday to @<177531421152247809>")
				time.Sleep(time.Second)
				dg.ChannelMessageSend(agora, "Happy birthday to you")
				time.Sleep(time.Second)
				dg.ChannelMessageSend(agora, "Groetjes, Alf de enige (maar stille) Thomas Bot microservice die een besef van tijd heeft. In een complot van Maartje, Vic en Sofie!")
			}
		}
	}()

	// small in memory structure to keep candidates to remove
	v.shouldRemove = map[string]bool{}

	go func() {
		for {
			time.Sleep(300 * time.Second)

			guilds, err := dg.UserGuilds(100, "", "")
			for len(guilds) > 0 {
				for _, guild := range guilds {
					log.Println("Checking", guild.Name)
					v.checkGuild(guild.ID)
				}
				guilds, err = dg.UserGuilds(100, "", guilds[len(guilds)-1].ID)
				if err != nil {
					break
				}
			}
		}
	}()

	log.Println("Thomas Bot Alf (cleanup) is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	return nil
}

func (v *cleanCmdOptions) checkGuild(guildID string) {
	state, err := v.dg.State.Guild(guildID)
	if err != nil {
		log.Println(err)
		return
	}

	channels, err := v.dg.GuildChannels(guildID)
	if err != nil {
		log.Println(err)
		return
	}

	conf, err := v.db.ConfigForGuild(guildID)
	if err != nil {
		log.Println(err)
		return
	}

	if conf == nil || len(conf.Hives) == 0 {
		return
	}

	for _, channel := range channels {
		if conf, isHive, _ := v.getConfigForRequestCategory(conf, channel); isHive && channel.Type == discordgo.ChannelTypeGuildVoice {
			if conf.Prefix != "" {
				if !strings.HasPrefix(channel.Name, conf.Prefix) {
					continue
				}
			}
			log.Println("looking at", channel.Name)
			inUse := false
			for _, vs := range state.VoiceStates {
				if vs.ChannelID == channel.ID {
					inUse = true
					log.Println(channel.Name, "is in use")
					break
				}
			}

			if !inUse {
				log.Println(channel.Name, "looks sus")
			}

			// on first occurance: mark to remove, on second occurance remove
			if wasMarkedAsRemove := v.shouldRemove[channel.ID]; wasMarkedAsRemove && !inUse {
				log.Println("Deleting", channel.ID, channel.Name)

				if conf.JunkyardCategoryID == "" {
					// no junkyard we need to delete
					_, err := v.dg.ChannelDelete(channel.ID)
					if err != nil {
						log.Println(err)
					}
				} else {
					// junkyard is set we need to move it there
					j, err := v.dg.Channel(conf.JunkyardCategoryID)
					if err != nil {
						log.Println(err)
						continue
					}
					_, err = v.dg.ChannelEditComplex(channel.ID, &discordgo.ChannelEdit{
						ParentID:             conf.JunkyardCategoryID,
						PermissionOverwrites: j.PermissionOverwrites,
					})
					if err != nil {
						log.Println(err)
					}
				}

				delete(v.shouldRemove, channel.ID)
			}

			v.shouldRemove[channel.ID] = !inUse
		}
	}
}

func (v *cleanCmdOptions) getConfigForRequestCategory(conf *db.Configuration, channel *discordgo.Channel) (*db.HiveConfiguration, bool, error) {
	for _, hive := range conf.Hives {
		if channel.ParentID == hive.VoiceCategoryID {
			return &hive, true, nil
		}
	}

	// no hive found
	return nil, false, nil
}
