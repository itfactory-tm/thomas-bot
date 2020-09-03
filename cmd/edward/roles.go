package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/commands/members"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

const itfDiscord = "687565213943332875"

func init() {
	rootCmd.AddCommand(NewServeCmd())
}

type serveCmdOptions struct {
	Token string
}

// NewServeCmd generates the `serve` command
func NewServeCmd() *cobra.Command {
	s := serveCmdOptions{}
	c := &cobra.Command{
		Use:     "roles",
		Short:   "Run send everybody a roles request",
		Long:    `Run send everybody a roles request`,
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

	return nil
}

func (s *serveCmdOptions) RunE(cmd *cobra.Command, args []string) error {
	dg, err := discordgo.New("Bot " + s.Token)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	m := members.NewMemberCommand()

	members, err := dg.GuildMembers(itfDiscord, "", 1000)
	if err != nil {
		return fmt.Errorf("error getting members: %w", err)
	}

	for _, member := range members {
		send := false
		for _, role := range member.Roles {
			if role == "687567949795557386" || role == "687568334379679771" || role == "687568470820388864" || role == "689844328528478262" {
				send = true
			}
		}
		if send {
			fmt.Println(member.User.Username)
			m.SendRoleDM(dg, member.User.ID)
			time.Sleep(30 * time.Second)
		}
	}

	return nil
}
