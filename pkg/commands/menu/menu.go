package menu

import (
	"github.com/itfactory-tm/thomas-bot/pkg/util/slash"
	"log"
	"time"

	"io/ioutil"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

const apiString = "https://tmmenumanagement.azurewebsites.net/api/Menu/"

type MenuCommand struct{}

func NewMenuCommand() *MenuCommand {
	return &MenuCommand{}
}

//	Register registers the handlers
func (h *MenuCommand) Register(registry command.Registry, server command.Server) {
	registry.RegisterInteractionCreate("menu", h.SayMenu)
}

//	InstallSlashCommands registers the slash commands
func (h *MenuCommand) InstallSlashCommands(session *discordgo.Session) error {
	app := discordgo.ApplicationCommand{
		Name: "menu",
		Description: "Loads the cafetaria menu",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type: discordgo.ApplicationCommandOptionString,
				Name: "campus",
				Description: "The campus to get the menu from",
				Required: true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name: "Geel",
						Value: "Geel",
					},
					{
						Name: "De Nayer",
						Value: "De Nayer",
					},
					{
						Name: "Lier",
						Value: "Lier",
					},
					{
						Name: "Antwerpen",
						Value: "Antwerpen",
					},
					{
						Name: "Mechelen",
						Value: "Mechelen",
					},
					{
						Name: "Turnhout",
						Value: "Turnhout",
					},
					{
						Name: "Vorselaar",
						Value: "Vorselaar",
					},
				},
			},
		},
	}

	if err := slash.InstallSlashCommand(session, "", app); err != nil {
		return err
	}

	return nil
}

//	SayMenu relays the menu
//	TODO: pull the different meals from the api
func (h *MenuCommand) SayMenu(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var selectedCampus = i.Data.Options[0].Value.(string)

	embed := &discordgo.MessageEmbed{
		Title: "Menu campus "+selectedCampus+" | "+time.Now().Format("2 Jan"),
		Color: 0x33FF33,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Meal",
				Value: "🥪 Sandwich\n"+
					"🍲 Main course\n"+
					"🥣 Soup\n"+
					"🥗 Vegetarian\n",
					Inline: true,
			},
			{
				Name: "Item",
				Value: "Cheese Sandwich\n"+
					"Spaghetti\n"+
					"Pea soup\n"+
					"Quinoa Salad\n",
				Inline: true,
			},
		},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "Here is the menu: "+GetSiteContent(selectedCampus),
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
		},
	})

	if(err!=nil){
		log.Println(err)
	}
}

// Info return the commands in this package
func (h *MenuCommand) Info() []command.Command {
	return []command.Command{}
}

//	GetSiteContent returns the json from the api
func GetSiteContent(campus string) string {
	res, err := http.Get(apiString+campus)
	if err != nil {
		log.Fatalf(err.Error())
	}

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return string(content)
}
