package slash

import (
	"fmt"
	"reflect"

	"github.com/bwmarrin/discordgo"
)

func InstallSlashCommand(session *discordgo.Session, guildID string, app discordgo.ApplicationCommand) error {
	cmds, err := session.ApplicationCommands(session.State.User.ID, guildID)
	if err != nil {
		return err
	}

	exists := false
	same := false
	var slashcmd *discordgo.ApplicationCommand
	for _, cmd := range cmds {
		if cmd.Name == app.Name {
			exists = true
			same = reflect.DeepEqual(app.Options, cmd.Options)
			slashcmd = cmd
		}
	}

	if !same && exists && slashcmd != nil {
		_, err = session.ApplicationCommandEdit(slashcmd.ID, session.State.User.ID, "", &app)
		return fmt.Errorf("error in ApplicationCommandEdit: %w", err)
	} else {
		_, err = session.ApplicationCommandCreate(session.State.User.ID, "", &app)
		return fmt.Errorf("error in ApplicationCommandCreate: %w", err)
	}
	return nil
}
