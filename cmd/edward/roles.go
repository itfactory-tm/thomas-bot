package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

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

	members, err := dg.GuildMembers("689786403596533841", "", 1000)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	for _, member := range members {
		dg.GuildMemberRoleAdd("689786403596533841", member.User.ID, "689819133206200590")
		log.Printf("Gave %s the guest role", member.User.Username)
	}

	/*
		m := members.NewMemberCommand()

		members, err := dg.GuildMembers(itfDiscord, "", 1000)
		if err != nil {
			return fmt.Errorf("error getting members: %w", err)
		}



	*/

	return nil
}
