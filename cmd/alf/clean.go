package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/discordha"

	"github.com/bwmarrin/discordgo"

	"github.com/kelseyhightower/envconfig"

	"github.com/spf13/cobra"
)

// TODO: automate these
const itfDiscord = "687565213943332875"

// default channel
const hiveCategoryID = "775436992136871957"

func init() {
	rootCmd.AddCommand(NewCleanCmd())
}

type cleanCmdOptions struct {
	Token string

	ha *discordha.HA
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

	// small in memory structure to keep candidates to remove
	shouldRemove := map[string]bool{}

	for {
		time.Sleep(60 * time.Second)

		state, err := dg.State.Guild(itfDiscord)
		if err != nil {
			log.Println(err)
			continue
		}

		channels, err := dg.GuildChannels(itfDiscord)
		if err != nil {
			return err
		}

		for _, channel := range channels {
			if channel.ParentID == hiveCategoryID && channel.Type == discordgo.ChannelTypeGuildVoice {
				inUse := false
				for _, vs := range state.VoiceStates {
					if vs.ChannelID == channel.ID {
						inUse = true
					}
				}

				// on first occurance: mark to remove, on second occurance remove
				if _, wasMarkedAsRemove := shouldRemove[channel.ID]; wasMarkedAsRemove && !inUse {
					log.Println("Deleting", channel.ID)
					dg.ChannelDelete(channel.ID)
				}

				if !inUse {
					shouldRemove[channel.ID] = true
				} else {
					delete(shouldRemove, channel.ID)
				}
			}
		}

	}
	return nil
}
