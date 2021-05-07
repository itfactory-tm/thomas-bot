package images

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"

	"github.com/itfactory-tm/thomas-bot/pkg/embed"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

type imageFunction func(s *discordgo.Session, in *discordgo.InteractionCreate) *discordgo.MessageEmbed

// ImagesCommands contains the tm!hello command
type ImagesCommands struct {
	images map[string]imageFunction
}

// NewImagesCommands gives a new ImagesCommands
func NewImagesCommands() *ImagesCommands {
	i := &ImagesCommands{}

	i.images = map[string]imageFunction{
		"erasmus":  i.sayErasmus,
		"partners": i.sayPartners,
		"loesje":   i.sayLoesje,
		"geit":     i.sayGeit,
		"paard":    i.sayPaard,
		"schaap":   i.saySchaap,
		"steun":    i.saySteun,
		"love":     i.sayLove,
	}
	return i
}

// Register registers the handlers
func (i *ImagesCommands) Register(registry command.Registry, server command.Server) {
	registry.RegisterInteractionCreate("image", i.slashCommand)
}

// InstallSlashCommands registers the slash commands
func (i *ImagesCommands) InstallSlashCommands(session *discordgo.Session) error {
	var choices []*discordgo.ApplicationCommandOptionChoice
	for name := range i.images {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  name,
			Value: name,
		})
	}

	app := discordgo.ApplicationCommand{
		Name:        "image",
		Description: "Gives an image",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "name of the picture",
				Required:    true,
				Choices:     choices,
			},
		},
	}

	cmds, err := session.ApplicationCommands(session.State.User.ID, "")
	if err != nil {
		return err
	}
	exists := false
	for _, cmd := range cmds {
		if cmd.Name == "image" {
			exists = reflect.DeepEqual(app.Options, cmd.Options)
		}
	}

	if !exists {
		_, err = session.ApplicationCommandCreate(session.State.User.ID, "", &app)
	}

	return err
}

// Info return the commands in this package
func (i *ImagesCommands) Info() []command.Command {
	return []command.Command{}
}

func (i *ImagesCommands) slashCommand(s *discordgo.Session, in *discordgo.InteractionCreate) {
	if len(in.Data.Options) > 0 {
		if key, ok := in.Data.Options[0].Value.(string); ok {
			if fn, ok := i.images[key]; ok {
				err := s.InteractionRespond(in.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionApplicationCommandResponseData{
						Embeds: []*discordgo.MessageEmbed{
							fn(s, in),
						},
					},
				})

				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	err := s.InteractionRespond(in.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "sorry I didn't find that image",
			Flags:   64,
		},
	})

	if err != nil {
		log.Println(err)
	}

}

func (i *ImagesCommands) sayErasmus(s *discordgo.Session, in *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	e := embed.NewEmbed()
	e.SetTitle("Erasmus @ ITfactory")
	e.SetImage("https://static.eyskens.me/thomas-bot/sem_2_2020.gif")
	e.SetURL("https://thomasmore365.sharepoint.com/sites/james/NL/international?tmbaseCampus=Geel")
	return e.MessageEmbed
}

func (i *ImagesCommands) sayPartners(s *discordgo.Session, in *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	e := embed.NewEmbed()
	e.SetTitle("Partners in education")
	e.SetImage("https://static.eyskens.me/thomas-bot/voorstelling_partners_in_education.png")
	return e.MessageEmbed
}

func (i *ImagesCommands) sayLove(s *discordgo.Session, in *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	e := embed.NewEmbed()
	guild, err := s.Guild(in.GuildID)
	if err != nil {
		log.Println(guild)
		return nil
	}
	e.SetTitle(fmt.Sprintf("<3 %s <3", guild.Name))
	e.SetImage("https://static.eyskens.me/thomas-bot/love.gif")
	return e.MessageEmbed
}

func (i *ImagesCommands) sayLoesje(s *discordgo.Session, in *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	j := rand.Intn(7)
	e := embed.NewEmbed()
	e.SetTitle("Loesje")
	e.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/loesje%d.png", j+1))
	return e.MessageEmbed
}

func (i *ImagesCommands) sayGeit(s *discordgo.Session, in *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	j := rand.Intn(4)
	e := embed.NewEmbed()
	e.SetTitle("Geit")
	e.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/geit%d.png", j+1))
	return e.MessageEmbed
}

func (i *ImagesCommands) sayPaard(s *discordgo.Session, in *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	j := rand.Intn(2)
	e := embed.NewEmbed()
	e.SetTitle("Paard")
	e.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/paard%d.png", j+1))
	return e.MessageEmbed
}

func (i *ImagesCommands) saySchaap(s *discordgo.Session, in *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	j := rand.Intn(9)
	e := embed.NewEmbed()
	e.SetTitle("Schaap")
	e.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/schaap%d.png", j+1))
	return e.MessageEmbed
}

func (i *ImagesCommands) saySteun(s *discordgo.Session, in *discordgo.InteractionCreate) *discordgo.MessageEmbed {
	j := rand.Intn(40)
	e := embed.NewEmbed()
	e.SetTitle("Steun")
	e.SetImage(fmt.Sprintf("https://static.eyskens.me/thomas-bot/examensteun/%02d.png", j+1))
	return e.MessageEmbed
}
