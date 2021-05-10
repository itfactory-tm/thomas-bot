package slash

import (
	"reflect"

	"github.com/bwmarrin/discordgo"
)

func InstallSlashCommand(session *discordgo.Session, guildID string, app discordgo.ApplicationCommand) error {
	cmds, err := session.ApplicationCommands(session.State.User.ID, guildID)
	if err != nil {
		return err
	}

	exists := true
	same := false
	var slashcmd *discordgo.ApplicationCommand
	for _, cmd := range cmds {
		if cmd.Name == "hive" {
			exists = true
			same = reflect.DeepEqual(app.Options, cmd.Options)
			slashcmd = cmd
		}
	}

	if !same && exists && slashcmd != nil {
		_, err = session.ApplicationCommandEdit(slashcmd.ID, session.State.User.ID, "", &app)
	} else if !same {
		_, err = session.ApplicationCommandCreate(session.State.User.ID, "", &app)
	}
	return err
}
