package main

import (
	"fmt"
	"log"
	"os"

	"github.com/itfactory-tm/thomas-bot/pkg/command"

	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("You are using the old entrypoint please update your tooling")
	os.Exit(1)
}

func onNewMember(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	if ok, err := ha.Lock(g); !ok {
		if err != nil {
			log.Printf("Error lockin on new memner: %q\n", err)
		}
		return
	}
	if g.GuildID != itfDiscord {
		ha.Unlock(g)
		return
	}

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
