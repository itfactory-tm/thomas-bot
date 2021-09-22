package menu

import (
	"log"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/embed"
	"github.com/itfactory-tm/thomas-bot/pkg/util/slash"

	"io/ioutil"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"

	"encoding/json"
)

const apiString = "https://tmmenumanagement.azurewebsites.net/api/Menu/"

// only works for Geel...
// not anymore!
type WeekMenu struct {
	Days [5]struct {
		MenuItems []CategoryDay
		Date      time.Time
	}
}

type CategoryDay struct {
	ShortDescriptionNL string
	ShortDescriptionEN string
	Category           struct {
		ID     string
		NameNL string
		NameEN string
	}
	ChoiceGroups []interface{}
}

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
		Name:        "menu",
		Description: "Loads the cafetaria menu",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "campus",
				Description: "The campus to get the menu from",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "Geel",
						Value: "Geel",
					},
					{
						Name:  "Lier",
						Value: "Lier",
					},
					/*{
						Name:  "De Nayer",
						Value: "De Nayer",
					},
					{
						Name:  "Lier",
						Value: "Lier",
					},
					{
						Name:  "Antwerpen",
						Value: "Antwerpen",
					},
					{
						Name:  "Mechelen",
						Value: "Mechelen",
					},
					{
						Name:  "Turnhout",
						Value: "Turnhout",
					},
					{
						Name:  "Vorselaar",
						Value: "Vorselaar",
					},*/
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
	var selectedCampus = i.ApplicationCommandData().Options[0].Value.(string)

	data := GetSiteContent(selectedCampus)

	currentMenu := []MenuData{}
	for _, item := range data {
		// if is today or after today
		if item.Curdate.After(time.Now()) || item.Curdate.Day() == time.Now().Day() {
			currentMenu = append(currentMenu, item)
		}
	}

	embeds := []*discordgo.MessageEmbed{}
	for _, item := range currentMenu {
		e := embed.NewEmbed()
		e.Title = item.Curdate.Format("Monday")
		for _, item := range item.Items {
			e.AddField(item.Category.NameEN, item.ShortDescriptionEN)
		}

		embeds = append(embeds, e.MessageEmbed)
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Here is the menu: ",
			Embeds:  embeds,
		},
	})

	if err != nil {
		log.Println(err)
	}

}

// Info return the commands in this package
func (h *MenuCommand) Info() []command.Command {
	return []command.Command{}
}

//	GetSiteContent returns the json from the api
func GetSiteContent(campus string) []MenuData {
	res, err := http.Get(apiString + campus)
	if err != nil {
		log.Fatalf(err.Error())
	}

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatalf(err.Error())
	}

	dataStr := ""
	json.Unmarshal(content, &dataStr) // yes the data is sent inside a string

	data := []MenuData{}
	json.Unmarshal([]byte(dataStr), &data)

	return data
}
