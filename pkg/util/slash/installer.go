package slash

import (
	"fmt"
	"reflect"

	"github.com/bwmarrin/discordgo"
)

func InstallSlashCommand(session *discordgo.Session, guildID string, app discordgo.ApplicationCommand) error {
	cmds, err := session.ApplicationCommands(session.State.User.ID, guildID)
	if err != nil {
		return fmt.Errorf("error in ApplicationCommands get: %w", err)
	}

	exists := false
	same := false
	var slashcmd *discordgo.ApplicationCommand
	for _, cmd := range cmds {
		if cmd.Name == app.Name {
			exists = true
			same = reflect.DeepEqual(app.Options, cmd.Options)
			if len(cmd.Options) == 0 && len(app.Options) == 0 {
				same = true // hack as the api retuns a nil and we don't
			}
			slashcmd = cmd
		}
	}

	if !same && exists && slashcmd != nil {
		_, err = session.ApplicationCommandEdit(session.State.User.ID, guildID, slashcmd.ID, &app)
		if err != nil {
			return fmt.Errorf("error in ApplicationCommandEdit %s: %w", slashcmd.Name, err)
		}
	} else if !same {
		_, err = session.ApplicationCommandCreate(session.State.User.ID, guildID, &app)
		if err != nil {
			return fmt.Errorf("error in ApplicationCommandCreate: %w", err)
		}
	}

	return nil
}
