package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/commands/members"
	"github.com/itfactory-tm/thomas-bot/pkg/db"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(NewServeCmd())
}

var studentRoles = map[string]bool{
	"687567949795557386": true, // 1ITF
	"687568334379679771": true, // 2ITF
	"687568470820388864": true, // 3ITF
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

	dbConn, err := db.NewLocalDB("./config.json")
	if err != nil {
		return fmt.Errorf("error creating database connection: %w", err)
	}

	roleCMD := members.NewMemberCommand(dbConn)

	// get all members
	var members []*discordgo.Member
	needsMore := true
	after := ""
	for needsMore {
		var newMembers []*discordgo.Member
		newMembers, err = dg.GuildMembers("687565213943332875", after, 1000)
		if err != nil {
			return fmt.Errorf("error getting members: %w", err)
		}
		members = append(members, newMembers...)
		if len(newMembers) < 1000 {
			needsMore = false
		} else {
			after = newMembers[len(newMembers)-1].User.ID
		}
	}

	for _, member := range members {
		isStudent := false
		// check if member is in the student role
		for _, role := range member.Roles {
			if studentRoles[role] {
				isStudent = true
				break
			}
		}

		if !isStudent {
			continue
		}

		// create DM
		dm, err := dg.UserChannelCreate(member.User.ID)
		if err != nil {
			log.Printf("error creating DM channel for %q: %s\n", member.User.ID, err)
			continue
		}

		// send message
		dg.ChannelMessageSend(dm.ID, "Hi there! The new academic year is almost here :) I am checking up on all students to give them new roles for next semester. Can you let me know which ones you want below?")
		roleCMD.SendRoleDM(dg, "687565213943332875", member.User.ID)
		log.Printf("Sent roles request to %q\n", member.User.ID)

		time.Sleep(time.Second * 5)
	}

	return nil
}
