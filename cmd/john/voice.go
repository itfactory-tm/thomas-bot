package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	discordha "github.com/meyskens/discord-ha"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/mixer"

	"github.com/kelseyhightower/envconfig"

	"github.com/spf13/cobra"
)

// audioConnected per Guild ID
var audioConnected = map[string]bool{}

func init() {
	rootCmd.AddCommand(NewVoiceCmd())
}

type voiceCmdOptions struct {
	Token         string
	EtcdEndpoints []string `envconfig:"ETCD_ENDPOINTS"`

	ha discordha.HA
}

// NewVoiceCmd generates the `serve` command
func NewVoiceCmd() *cobra.Command {
	s := voiceCmdOptions{}
	c := &cobra.Command{
		Use:     "voice",
		Short:   "Run the voice server",
		Long:    `This is a separate instance for voice services. This can only run in a single replica`,
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

func (v *voiceCmdOptions) Validate(cmd *cobra.Command, args []string) error {
	if v.Token == "" {
		return errors.New("No token specified")
	}

	return nil
}

func (v *voiceCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	log.Println("Starting John...")

	ctx := context.TODO()

	dg, err := discordgo.New("Bot " + v.Token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	haLogger := log.New(os.Stderr, "discordha: ", log.Ldate|log.Ltime)
	v.ha, err = discordha.New(&discordha.Config{
		Session:                          dg,
		HA:                               len(v.EtcdEndpoints) > 0,
		EtcdEndpoints:                    v.EtcdEndpoints,
		Context:                          ctx,
		LockTTL:                          1 * time.Second,
		LockUpdateInterval:               500 * time.Millisecond,
		Log:                              *haLogger,
		DoNotParticipateInLeaderElection: true,
	})
	if err != nil {
		return fmt.Errorf("error creating Discord HA: %w", err)
	}

	log.Println("Watching etcd...")
	voiceQueueChan := v.ha.WatchVoiceCommands(ctx, "thomasbot")

	err = dg.Open()
	if err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}
	defer dg.Close()

	for {
		q := <-voiceQueueChan
		if audioConnected[q.GuildID] {
			continue
		}
		connected := make(chan struct{})
		fmt.Printf("Connecting to %s\n", q.ChannelID)
		go v.connectVoice(dg, connected, q.GuildID, q.ChannelID, q.UserID)
		<-connected
		// send again for voice to pick up
		v.ha.SendVoiceCommand(q)
	}

	return nil
}

func (v *voiceCmdOptions) connectVoice(dg *discordgo.Session, connected chan struct{}, guildID, channelID, userID string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if audioConnected[guildID] {
		connected <- struct{}{}
		return
	}

	audioConnected[guildID] = true
	voiceQueueChan := v.ha.WatchVoiceCommands(ctx, "thomasbot")

	dgv, err := dg.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		log.Fatal(err)
		return
	}

	connected <- struct{}{}

	encoder := mixer.NewEncoder()
	encoder.VC = dgv
	go encoder.Run()

	doneChan := make(chan struct{})
	go func() {
		var i uint64
		for {
			select {
			case f := <-voiceQueueChan:
				log.Println(f)
				go encoder.Queue(uint64(i), path.Join("./sounds/", f.File))
				i++
			case <-doneChan:
				return
			}
		}
	}()

	time.Sleep(5 * time.Second) // waiting for first audio
	for !encoder.HasFinishedAll() {
		time.Sleep(5 * time.Second)
	}

	// Close connections once all are played
	dgv.Disconnect()
	dgv.Close()
	encoder.Stop()
	audioConnected[guildID] = false
	doneChan <- struct{}{}
}
