package images

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/itfactory-tm/thomas-bot/pkg/embed"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

// ImagesCommands contains the tm!hello command
type ImagesCommands struct{}

// NewImagesCommands gives a new ImagesCommands
func NewImagesCommands() *ImagesCommands {
	return &ImagesCommands{}
}

func (i *ImagesCommands) Register(registry command.Registry, server command.Server) {
	registry.RegisterMessageCreateHandler("erasmus", i.sayErasmus)
	registry.RegisterMessageCreateHandler("partners", i.sayPartners)
	registry.RegisterMessageCreateHandler("love", i.sayLove)
	registry.RegisterMessageCreateHandler("loesje", i.sayGeit)
	registry.RegisterMessageCreateHandler("geit", i.sayErasmus)
	registry.RegisterMessageCreateHandler("paard", i.sayPaard)
	registry.RegisterMessageCreateHandler("schaap", i.saySchaap)
	registry.RegisterMessageCreateHandler("command", i.saySteun)
}

func (i *ImagesCommands) SayHello(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Beep bop boop! I am Thomas Bot, fork me on GitHub!")
}

func (i *ImagesCommands) Info() []command.Command {
	return []command.Command{
		command.Command{
			Name:        "geit",
			Category:    command.CategoryFun,
			Description: "De E-F-blok geiten nu ook online",
			Hidden:      false,
		},
		command.Command{
			Name:        "paard",
			Category:    command.CategoryFun,
			Description: "De E-F-blok paardjes nu ook online",
			Hidden:      false,
		},
		command.Command{
			Name:        "schaap",
			Category:    command.CategoryFun,
			Description: "De E-F-blok schapen nu ook online",
			Hidden:      false,
		},
		command.Command{
			Name:        "steun",
			Category:    command.CategoryFun,
			Description: "Examens? We komen er samen door!",
			Hidden:      false,
		},
		command.Command{
			Name:        "erasmus",
			Category:    command.CategoryLinks,
			Description: "Link naar erasmus",
			Hidden:      false,
		},
		command.Command{
			Name:        "partners",
			Category:    command.CategoryLinks,
			Description: "Link naar partners",
			Hidden:      false,
		},
		command.Command{
			Name:        "love",
			Category:    command.CategoryFun,
			Description: "Toon wat liefde aan elkaar <3",
			Hidden:      false,
		},
		command.Command{
			Name:        "loesje",
			Category:    command.CategoryFun,
			Description: "'Een fan van loesje' heeft wijze spreuken",
			Hidden:      false,
		},
	}
}

func (i *ImagesCommands) sayErasmus(s *discordgo.Session, m *discordgo.MessageCreate) {
	e := embed.NewEmbed()
	e.SetTitle("Erasmus @ ITfactory")
	e.SetImage("https://static.eyskens.me/thomas-bot/sem_2_2020.gif")
	e.SetURL("https://thomasmore365.sharepoint.com/sites/james/NL/international?tmbaseCampus=Geel")
	s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
}

func (i *ImagesCommands) sayPartners(s *discordgo.Session, m *discordgo.MessageCreate) {
	e := embed.NewEmbed()
	e.SetTitle("Partners in education")
	e.SetImage("https://static.eyskens.me/thomas-bot/voorstelling_partners_in_education.png")
	s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
}

func (i *ImagesCommands) sayLove(s *discordgo.Session, m *discordgo.MessageCreate) {
	e := embed.NewEmbed()
	e.SetTitle("<3 IT-Factory <3")
	e.SetImage("https://static.eyskens.me/thomas-bot/love.gif")
	s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
}

func (i *ImagesCommands) sayLoesje(s *discordgo.Session, m *discordgo.MessageCreate) {
	j := rand.Intn(7)
	e := embed.NewEmbed()
	e.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/loesje%d.png", j+1))
	s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
}

func (i *ImagesCommands) sayGeit(s *discordgo.Session, m *discordgo.MessageCreate) {
	j := rand.Intn(9)
	e := embed.NewEmbed()
	e.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/geit%d.png", j+1))
	s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
}

func (i *ImagesCommands) sayPaard(s *discordgo.Session, m *discordgo.MessageCreate) {
	j := rand.Intn(3)
	e := embed.NewEmbed()
	e.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/paard%d.png", j+1))
	s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
}

func (i *ImagesCommands) saySchaap(s *discordgo.Session, m *discordgo.MessageCreate) {
	j := rand.Intn(9)
	e := embed.NewEmbed()
	e.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/schaap%d.png", j+1))
	s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
}

func (i *ImagesCommands) saySteun(s *discordgo.Session, m *discordgo.MessageCreate) {
	j := rand.Intn(40)
	e := embed.NewEmbed()
	e.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/examensteun/%02d.png", j+1))
	log.Println(fmt.Sprintf("https://static.eyskens.me/thomas-bot/examensteun/%02d.png", j+1))
	s.ChannelMessageSendEmbed(m.ChannelID, e.MessageEmbed)
}
