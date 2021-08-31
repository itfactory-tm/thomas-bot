package schedule

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/itfactory-tm/thomas-bot/pkg/db"
	"github.com/itfactory-tm/thomas-bot/pkg/embed"
	"github.com/itfactory-tm/thomas-bot/pkg/util/slash"

	ical "github.com/arran4/golang-ical"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

// ScheduleCommand contains the schedule command
type ScheduleCommand struct {
	db db.Database
}

// NewScheduleCommand gives a new ScheduleCommand
func NewScheduleCommand(db db.Database) *ScheduleCommand {
	return &ScheduleCommand{
		db: db,
	}
}

// Register registers the handlers
func (s *ScheduleCommand) Register(registry command.Registry, server command.Server) {
	registry.RegisterInteractionCreate("schedule", s.SaySchedule)
}

// InstallSlashCommands registers the slash commands
func (s *ScheduleCommand) InstallSlashCommands(session *discordgo.Session) error {
	guilds, err := s.db.GetAllConfigurations()
	if err != nil {
		return err
	}

	for _, config := range guilds {
		if len(config.Schedules) < 1 {
			log.Println(config.GuildID, "has no schedules")
			continue
		}
		classes := []*discordgo.ApplicationCommandOptionChoice{}
		for _, class := range config.Schedules {
			classes = append(classes, &discordgo.ApplicationCommandOptionChoice{
				Name:  class.ClassName,
				Value: class.ClassName,
			})
		}

		err := slash.InstallSlashCommand(session, config.GuildID, discordgo.ApplicationCommand{
			Name:        "schedule",
			Description: "get a class schedule",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "class",
					Type:        discordgo.ApplicationCommandOptionString,
					Description: "the class name",
					Choices:     classes,
					Required:    true,
				},
				{
					Name:        "publish",
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Description: "post the reply in channel",
					Required:    false,
				},
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

type classSchedule struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Room      string
	Teachers  string
}

// SaySchedule handles a schedule interaction
func (s *ScheduleCommand) SaySchedule(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	conf, err := s.db.ConfigForGuild(i.GuildID)
	if err != nil {
		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error fetching DB",
				Flags:   64, // hidden
			},
		})
		return
	}

	publish := false
	var name string
	var ok bool
	if len(i.ApplicationCommandData().Options) < 1 {
		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No options sent",
				Flags:   64, // hidden
			},
		})
		return
	}

	if name, ok = i.ApplicationCommandData().Options[0].Value.(string); !ok {
		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No valid option",
				Flags:   64, // hidden
			},
		})

		return
	}

	if len(i.ApplicationCommandData().Options) >= 2 {
		publish = i.ApplicationCommandData().Options[1].Value.(bool)
	}

	url := ""
	for _, sch := range conf.Schedules {
		if sch.ClassName == name {
			url = sch.URL
			break
		}
	}

	if url == "" {
		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Could not get schedule: unknown class",
				Flags:   64, // hidden
			},
		})
	}

	events, err := parseSchedule(url)
	if err != nil {
		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Could not get schedule: " + err.Error(),
				Flags:   64, // hidden
			},
		})
		return
	}

	embeds := []*discordgo.MessageEmbed{}
	for _, event := range events {
		e := embed.NewEmbed()
		e.SetTitle(event.Name)
		e.SetAuthor(event.Room)
		e.SetDescription(event.StartTime.Format("Mon Jan 2 15:04") + " - " + event.EndTime.Format("15:04") + "\n" + event.Teachers)

		embeds = append(embeds, e.MessageEmbed)
	}

	if len(embeds) > 10 {
		embeds = embeds[:10]
	}

	flags := 64
	if publish {
		flags = 0
	}
	if len(embeds) == 0 {
		sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No classes found in the next week, enjoy!",
				Flags:   uint64(flags),
			},
		})
		return
	}

	err = sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Here is your schedule:",
			Embeds:  embeds,
			Flags:   uint64(flags),
		},
	})
	if err != nil {
		log.Println("ScheduleCommand.SaySchedule:", err)
	}
}

// function that parses a given ical URL and returns a schedule
func parseSchedule(icalURL string) ([]classSchedule, error) {
	resp, err := http.Get(icalURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	cal, err := ical.ParseCalendar(resp.Body)
	if err != nil {
		return nil, err
	}

	out := make([]classSchedule, 0)
	for _, icalEvent := range cal.Events() {
		fixICalTime(icalEvent)
		start, err := icalEvent.GetStartAt()
		if err != nil {
			log.Println(err)
			continue
		}
		end, err := icalEvent.GetEndAt()
		if err != nil {
			log.Println(err)
			continue
		}

		var teachers string
		if p := strings.Split(icalEvent.GetProperty(ical.ComponentPropertyDescription).Value, "Staff member(s):"); len(p) > 1 {
			teachers = strings.Split(strings.TrimSpace(p[1]), "\\n")[0]
		}

		if end.After(time.Now()) && start.Before(time.Now().Add(time.Hour*24*7)) { // only show classes in the next 7 days
			out = append(out, classSchedule{
				Name:      icalEvent.GetProperty(ical.ComponentPropertySummary).Value,
				StartTime: start,
				EndTime:   end,
				Room:      icalEvent.GetProperty(ical.ComponentPropertyLocation).Value,
				Teachers:  teachers,
			})
		}
	}

	return out, nil
}

func fixICalTime(icalEvent *ical.VEvent) {
	icalStart := icalEvent.GetProperty(ical.ComponentPropertyDtStart)
	if !strings.Contains(icalStart.Value, "Z") {
		icalStart.Value = icalStart.Value + "Z"
	}

	icalEvent.SetProperty(ical.ComponentPropertyDtStart, icalStart.Value)

	icalEnd := icalEvent.GetProperty(ical.ComponentPropertyDtEnd)
	if !strings.Contains(icalEnd.Value, "Z") {
		icalEnd.Value = icalEnd.Value + "Z"
	}
	icalEvent.SetProperty(ical.ComponentPropertyDtEnd, icalEnd.Value)
}

// Info return the commands in this package
func (s *ScheduleCommand) Info() []command.Command {
	return []command.Command{}
}
