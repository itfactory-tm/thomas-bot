package pronostiek

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/itfactory-tm/thomas-bot/pkg/util/slash"

	"github.com/itfactory-tm/thomas-bot/pkg/embed"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

// PronostiekCommand contains the /pronostiek command
type PronostiekCommand struct {
}

// PronostiekCommand gives a new PronostiekCommand
func NewPronostiekCommand() *PronostiekCommand {
	return &PronostiekCommand{}
}

// Register registers the handlers
func (p *PronostiekCommand) Register(registry command.Registry, server command.Server) {
	registry.RegisterInteractionCreate("pronostiek", p.slashCommand)
}

// InstallSlashCommands registers the slash commands
func (p *PronostiekCommand) InstallSlashCommands(session *discordgo.Session) error {
	return slash.InstallSlashCommand(session, "687565213943332875", discordgo.ApplicationCommand{
		Name:        "pronostiek",
		Description: "Post the Current EK Pronostiek",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "rank",
				Description: "name of rank",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "studenten",
						Value: "Studenten",
					},
					{
						Name:  "docenten",
						Value: "Docenten",
					},
				},
			},
		},
	})
}

// Info return the commands in this package
func (p *PronostiekCommand) Info() []command.Command {
	return []command.Command{}
}

func (p *PronostiekCommand) slashCommand(s *discordgo.Session, in *discordgo.InteractionCreate) {
	rank := ""
	if len(in.ApplicationCommandData().Options) > 0 {
		if key, ok := in.ApplicationCommandData().Options[0].Value.(string); ok {
			rank = key
		}
	}

	resp, err := http.Get(fmt.Sprintf("https://prono.inmijneendje.be/api/rank%s", rank))
	if err != nil {
		s.InteractionRespond(in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("An error occured: %q", err),
				Flags:   64,
			},
		})
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.InteractionRespond(in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("An error occured: %q", err),
				Flags:   64,
			},
		})
		return
	}

	var data []Rank
	err = json.Unmarshal(body, &data)
	if err != nil {
		s.InteractionRespond(in.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("An error occured: %q", err),
				Flags:   64,
			},
		})
		return
	}

	e := embed.NewEmbed()
	e.SetTitle(fmt.Sprintf("Rank %s", rank))

	out := []*discordgo.MessageEmbed{}

	fieldCount := 0
	for i, r := range data {
		fieldCount += 3
		if fieldCount > 15 { // if over 15 (5 rows) reset the embed for a 2nd one
			e.InlineAllFields()
			out = append(out, e.MessageEmbed)

			e = embed.NewEmbed()
			e.SetTitle(fmt.Sprintf("Rank %s", rank))

			fieldCount = 3 // because we will send 3 already
		}
		e.AddField("Name", fmt.Sprintf("%s %s", getMedal(i), r.Name))
		e.AddField("Score", r.Totalscore)
		e.AddField("Correct", fmt.Sprintf("%d", r.AllCorrect))

	}

	e.InlineAllFields()
	out = append(out, e.MessageEmbed)

	err = s.InteractionRespond(in.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ook een gokje wagen? https://prono.inmijneendje.be/",
			Embeds:  out,
		},
	})

	if err != nil {
		log.Println(err)
	}

}

func getMedal(i int) string {
	switch i {
	case 0:
		return "🥇"
	case 1:
		return "🥈"
	case 2:
		return "🥉"
	default:
		return fmt.Sprintf("%d", i+1)
	}
}
