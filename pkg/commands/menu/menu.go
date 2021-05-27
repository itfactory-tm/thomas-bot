package menu

import (
	"github.com/itfactory-tm/thomas-bot/pkg/util/slash"
	"log"

	"io/ioutil"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

const test = "828356586426597377"

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
		Options: []*discordgo.ApplicationCommandOption{},
	}

	if err := slash.InstallSlashCommand(session, "", app); err != nil {
		return err
	}

	return nil
}

//	SayMenu relays the menu
func (h *MenuCommand) SayMenu(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: GetSiteContent(),
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

func GetSiteContent() string {
	res, err := http.Get("https://tmmenumanagement.azurewebsites.net/api/Menu/Geel")
	if err != nil {
		log.Fatalf(err.Error())
	}

	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatalf(err.Error())
	}

	return "Success!, response: "+string(content)
}
