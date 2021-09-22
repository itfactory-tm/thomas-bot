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

const apiString = "https://tmmenumanagement.azurewebsites.net/api/WeekMenu/"

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
					{
						Name:  "Turnhout",
						Value: "Turnhout",
					},
					{
						Name:  "Vorselaar",
						Value: "Vorselaar",
					},
					/*{
						Name:  "De Nayer",
						Value: "De Nayer", //De Nayer gebruikt "undefined"
					},
					{
						Name:  "Antwerpen",
						Value: "Antwerpen",
					},
					{
						Name:  "Mechelen",
						Value: "Mechelen",
					},*/
					/*
						----------------
							Bijkomend
						----------------
						Antwerpen campus Sanderus gebruikt "undefined"
						Antwerpen campus Sint-Andries gebruikt "undefined"
						De Nayer gebruikt "undefined"
						Campus De Ham gebruikt "Mechelen"
						Campus De Vest gebruikt "Mechelen"
						Campus Kruidtuin gebruikt "Mechelen"
						Campus Lucas Faydherbe gebruikt "Mechelen"
					*/
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
	if len(data) == 0 {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "That campus does not have a menu for this week yet!",
			},
		})

		if err != nil {
			log.Println(err)
		}
		return
	}

	pdata := data[0].(map[string]interface{})

	var items map[string]interface{}
	var startDate time.Time

	// retrieve list of categories and the startdate
	for k, v := range pdata {
		switch k {
		case "items":
			items = v.(map[string]interface{})
		case "startdate":
			startDate, _ = time.Parse(time.RFC3339, v.(string))
		}
	}

	var categoryWeeks []map[string]interface{}

	// extrapolate categories from list
	for _, v := range items {
		categoryWeeks = append(categoryWeeks, v.(map[string]interface{}))
	}

	var finalMenu WeekMenu

	// initialize finalMenu dates for later use
	for a := range finalMenu.Days {
		finalMenu.Days[a].Date = startDate.Add(time.Duration(a * 86400000000000)) //24h * 3600s/h * 1000 000 000ns/s
	}

	// Pull the actual menu data from the categories
	// and group by week
	for _, categoryweek := range categoryWeeks {
		for k, v := range categoryweek {
			var dayJ, _ = json.Marshal(v)
			var day CategoryDay
			err := json.Unmarshal(dayJ, &day) // easiest way to get our CategoryDay struct out is by converting to and from JSON
			if err != nil {
				log.Fatalf(err.Error())
			}
			switch k {
			case "Monday":
				finalMenu.Days[0].MenuItems = append(finalMenu.Days[0].MenuItems, day)
			case "Tuesday":
				finalMenu.Days[1].MenuItems = append(finalMenu.Days[1].MenuItems, day)
			case "Wednesday":
				finalMenu.Days[2].MenuItems = append(finalMenu.Days[2].MenuItems, day)
			case "Thursday":
				finalMenu.Days[3].MenuItems = append(finalMenu.Days[3].MenuItems, day)
			case "Friday":
				finalMenu.Days[4].MenuItems = append(finalMenu.Days[4].MenuItems, day)
			}
		}
	}

	embeds := []*discordgo.MessageEmbed{}
	for _, day := range finalMenu.Days {
		if day.Date.After(time.Now()) || day.Date.Day() == time.Now().Day() {
			e := embed.NewEmbed()

			e.Title = day.Date.Format("Monday")

			// Check if the fields contain data
			for _, item := range day.MenuItems {
				if item.ShortDescriptionEN == "" {
					if item.ShortDescriptionNL != "" {
						e.AddField(item.Category.NameNL, item.ShortDescriptionNL)
					} else if item.Category.NameEN != "" {
						e.AddField(item.Category.NameEN, "There is no "+item.Category.NameEN+" available today")
					}
				} else {
					e.AddField(item.Category.NameEN, item.ShortDescriptionEN)
				}
			}
			if len(e.Fields) == 0 {
				e.AddField("​", "There is no menu available this day")
			}

			e.InlineAllFields()

			embeds = append(embeds, e.MessageEmbed)
		}
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
func GetSiteContent(campus string) []interface{} {
	res, err := http.Get(apiString + campus)
	if err != nil {
		log.Fatalf(err.Error())
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = res.Body.Close()
	if err != nil {
		log.Fatalf(err.Error())
	}

	dataStr := ""
	err = json.Unmarshal(content, &dataStr) // yes the data is sent inside a string
	if err != nil {
		log.Fatalf(err.Error())
	}

	/*
		Omzetten naar een slice van een interface?!??

		Wat betekent dit en waarom wordt ik hier ziek van?
		Er is een groot probleem; er staan ID's in de JSON als namen,
		waardoor we geen structuur met vaste namen kunnen declareren
		(dit gaat, maar encoding/json zet structs met onbekende structuur
		dan om naar lege structs)

		door om te zetten naar dit formaat wordt alles behouden, maar moeten
		we alles manueel uit de data halen.
		Bovendien zal alles kapot gaan als er iets zou veranderen aan de
		structuur van de data
	*/
	var data []interface{}
	err = json.Unmarshal([]byte(dataStr), &data)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return data
}
