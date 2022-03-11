package links

import (
	"log"

	"github.com/itfactory-tm/thomas-bot/pkg/util/slash"

	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

// LinkCommands contains the /link command
type LinkCommands struct {
	infos    []command.Command
	iceinfos []string
	output   map[string]string
	registry command.Registry
}

// NewLinkCommands gives a new LinkCommands
func NewLinkCommands() *LinkCommands {
	return &LinkCommands{
		infos:  []command.Command{},
		output: map[string]string{},
	}
}

// Register registers the handlers
func (l *LinkCommands) Register(registry command.Registry, server command.Server) {
	registry.RegisterInteractionCreate("link", l.slashCommand)
	l.registry = registry

	l.buildLinks()
}

// InstallSlashCommands registers the slash commands
func (l *LinkCommands) InstallSlashCommands(session *discordgo.Session) error {
	if l.registry == nil {
		return nil
	}

	var icechoices []*discordgo.ApplicationCommandOptionChoice
	for _, info := range l.iceinfos {
		icechoices = append(icechoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  info,
			Value: info,
		})
	}

	slash.InstallSlashCommand(session, "808273924600365058", discordgo.ApplicationCommand{
		Name:        "link",
		Description: "Gives a useful link",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "name of the link",
				Required:    true,
				Choices:     icechoices,
			},
		},
	})

	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, info := range l.infos {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  info.Name,
			Value: info.Name,
		})
	}

	return slash.InstallSlashCommand(session, "687565213943332875", discordgo.ApplicationCommand{
		Name:        "link",
		Description: "Gives a useful link",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "name of the link",
				Required:    true,
				Choices:     choices,
			},
		},
	})
}

func (l *LinkCommands) buildLinks() {
	// Warning: these are exactly 25 links the maximum options of a Discord slash command

	l.registerLinkCommand("bot", "Link naar de git repo van deze bot", "Biep Boep, bekijk zeker mijn git repo https://github.com/itfactory-tm/thomas-bot")
	//l.registerLinkCommand("campusshop", "Link naar campusshop", "Een kijkje nemen in de campusshop? https://www.campiniamedia.be/mvc/index.jsp")
	l.registerLinkCommand("canvas", "Link naar Canvas", "Bekijk hier je leerplatform (Canvas): https://thomasmore.instructure.com/")
	l.registerInfoDagCommand("centen", "Link naar financiële informatie", "Wil je het financiële aspect van verder studeren bekijken? https://centenvoorstudenten.be/")
	//l.registerLinkCommand("coderood", "Meer info over code rood", "Heb je meer info over code rood nodig? https://www.thomasmore.be/update-code-rood")
	//l.registerLinkCommand("corona", "Link naar Corona informatie", "Zit je met vragen hoe thomasmore omgaat met corona? Bekijk dan zeker deze pagina: https://thomasmore365.sharepoint.com/sites/s-Studentenvoorzieningen/SitePages/Corona.aspx")
	l.registerInfoDagCommand("discord", "Link naar Discord documentatie", "Nog een beetje in de war over hoe Discord werkt?: https://support.discordapp.com/hc/nl")
	l.registerInfoDagCommand("ects", "Link naar ECTS fiches", "Bekijk hier de ECTS fiches van ELO-ICT: http://onderwijsaanbodkempen.thomasmore.be/2021/opleidingen/n/CQ_51236204.htm \nToegepaste Informatica: http://onderwijsaanbodkempen.thomasmore.be/2021/opleidingen/n/CQ_51236221.htm \nAlle ECTS fiches http://ects.thomasmore.be/")
	l.registerInfoDagCommand("emt", "Link naar EMT", "Heeft de IT-Factory een eigen studentenvereniging? Jazeker: https://www.facebook.com/StudentenverenigingEMT")
	l.registerLinkCommand("examen", "Link naar info over examens", "Alles over de examens vind je hier: https://thomasmore365.sharepoint.com/sites/s.itfactory/SitePages/Examens.aspx")
	l.registerInfoDagCommand("fb", "Link naar Facebook paginas", "Bekijk hier onze Facebook pagina van Toegepaste informatica: https://www.facebook.com/ToegepasteInformatica.ThomasMoreBE & ELO-ICT: https://www.facebook.com/ElektronicaICT.ThomasMoreBE & ACS: https://www.facebook.com/ACS.ThomasMoreBE")
	l.registerLinkCommand("icecube", "Link naar ice-cube", "Ice-cube, wat is dat? https://www.thomasmore.be/ice-cube Kom bij de Discord server! https://discord.thomasmore.be/")
	l.registerInfoDagCommand("inschrijven", "Link naar inschrijven", "Wil je je inschrijven? Dat kan hier! https://www.thomasmore.be/inschrijven")
	//l.registerInfoDagCommand("junior", "Link naar Junior College", "Benieuwd wat Junior College is? Bekijk het hier: https://www.thomasmore.be/site/junior-university-college")
	l.registerLinkCommand("kot", "Link naar kot informatie", "Informatie nodig rond op kot gaan? https://www.thomasmore.be/studenten/op-kot")
	l.registerLinkCommand("kuloket", "Link naar KUloket", "Kuloket raadplegen? https://kuloket.be")
	//l.registerInfoDagCommand("laptop", "Link naar info over laptops", "Welk materiaal heb ik nodig om in de IT-Factory te kunnen starten? https://www.thomasmore.be/sites/www.thomasmore.be/files/Laptopspecificaties%20voor%20IT%20Factory-studenten%202019-2020.pdf")
	//l.registerLinkCommand("lunch", "Link naar weekmenu", "Heb je honger? Bekijk hier het menu voor deze week: https://thomasmore365.sharepoint.com/sites/James/NL/stuvo/Paginas/Weekmenu.aspx?tmbaseCampus=Geel")
	l.registerInfoDagCommand("oho", "Link naar OHO", "Werken en studeren combineren? Dat kan zeker! https://www.thomasmore.be/opleidingen/professionele-bachelor/toegepaste-informatica/toegepaste-informatica-combinatie-werken-en-studeren-oho")
	//l.registerLinkCommand("onlineexamen", "Link naar info over online examens", "Instructies voor het online examen vind je hier: https://thomasmore365.sharepoint.com/sites/s-icts/SitePages/Digitaal-schriftelijk-examen.aspx\nFAQ over digitaal examineren: https://thomasmore365.sharepoint.com/sites/s-icts/SitePages/FAQ-digitaal-examineren.aspx\nTIP: bekijk zeker ook de canvas cursus voor vak specifieke info")
	l.registerLinkCommand("pictures", "Fotoalbum van IT Factory", "De Flickr-link voor IT Factory: https://www.flickr.com/photos/itfactorygeel/albums/with/72157711381764072")
	//l.registerLinkCommand("positief", "Link naar covid 19 meld formulier", "Heb je een bevestigde covid-19 besmetting? Laat dit dan hier weten: https://thomasmore365.sharepoint.com/sites/s-Studentenadministratie/SitePages/Melden-van-een-bevestigde-COVID-19-besmetting.aspx")
	l.registerLinkCommand("printen", "Link naar printen", "Meer informatie nodig over printen? https://thomasmore365.sharepoint.com/sites/s-Leercentrum/SitePages/Printen.aspx Je printkrediet opladen? https://printbeheer.thomasmore.be/")
	l.registerLinkCommand("rooster", "Link naar lessenrooster", "Bekijk hier je lessenrooster: https://rooster.thomasmore.be/")
	l.registerLinkCommand("sharepoint", "Link naar Studentenportaal", "Bekijk hier de 365 sharepoint van de ITFactory: https://thomasmore365.sharepoint.com/sites/s.itfactory/SitePages/Start.aspx")
	l.registerLinkCommand("sinners", "Link naar Sinners", "Wat is Sinners? https://sinners.be/ Hulp nodig omtrend Sinners? Gebruik het #sinners kanaal en we helpen je!")
	l.registerInfoDagCommand("studenten", "Link naar studenten info", "Op zoek naar meer algemene info rondom verder studeren? https://www.thomasmore.be/studenten")
	l.registerLinkCommand("studentenraad", "Contact opnemen met de studentenraad", "Wil je contact opnemen met de studentenraad? Stuur ze een mailtje via: studentenraad.itfactory@thomasmore.be")
	l.registerInfoDagCommand("stuvo", "Link naar Stuvo", "Heb je nood aan een goed gesprek? Neem dan zeker contact op met Stuvo: https://thomasmore365.sharepoint.com/sites/s-Studentenvoorzieningen")
	l.registerLinkCommand("template", "Link naar TM huisstijl templates", "Hier vind je de TM huisstijl templates: \nhttps://static.eyskens.me/tm-template/ppt-new.pptx, \nhttps://static.eyskens.me/tm-template/ppt-old.pptx, \nhttps://static.eyskens.me/tm-template/word-nl.docx, \nhttps://static.eyskens.me/tm-template/word-en.docx")
	//l.registerLinkCommand("twitch", "Link naar ITF Twitch kanaal", "Af en toe livestreamen we wat games op ons Twitch kanaal: https://www.twitch.tv/itfactorygaming")
	l.registerLinkCommand("wallet", "Link naar wallet", "Hoeveel staat er nog op mijn studentenkaart? https://thomasmore.mynetpay.be/")
	//l.registerLinkCommand("webcam", "Campus webcams", "B300 Camera 1: https://www.twitch.tv/maartjeme \nB300 Camera 2: https://www.twitch.tv/b300camera2\nGeitjes: https://www.twitch.tv/tmgeitlive")
	l.registerInfoDagCommand("website", "Link naar Thomas More website", "Bezoek onze website: https://thomasmore.be/opleidingen/professionele-bachelor/it-factory")
	l.registerLinkCommand("atomos", "Link naar Atomos kalender", "Hier vind je de Atomos kalender: https://atomosvzw.be/kalender.php")

	l.registerICELinkCommand("level1", "", "https://donut.sinners.be/on/Level_1.pdf")
	l.registerICELinkCommand("level2", "", "https://donut.sinners.be/ly/Level_2.pdf")
	l.registerICELinkCommand("level3", "", "https://donut.sinners.be/go/Level_3.pdf")
	l.registerICELinkCommand("level4", "", "https://donut.sinners.be/od/Level_4.pdf")
	l.registerICELinkCommand("level5", "", "https://donut.sinners.be/vi/Level_5.pdf")
	l.registerICELinkCommand("level6", "", "https://donut.sinners.be/be/Level_6.pdf")
	l.registerICELinkCommand("level7", "", "https://donut.sinners.be/si/Level_7.pdf")
	l.registerICELinkCommand("level8", "", "https://donut.sinners.be/ce/Level_8.pdf")
	l.registerICELinkCommand("brainstorm", "", "https://donut.sinners.be/word/Document_Brainstorm.docx")
	l.registerICELinkCommand("fase1", "", "https://donut.sinners.be/word/Document_Fase_1.docx")
	l.registerICELinkCommand("fase2", "", "https://donut.sinners.be/word/Document_Fase_2.docx")
	l.registerICELinkCommand("fase3", "", "https://donut.sinners.be/word/Document_Fase_3.docx")
	l.registerICELinkCommand("pitch", "", "https://donut.sinners.be/word/Document_Pitch.docx")
}

func (l *LinkCommands) registerLinkCommand(name, helpText, response string) {
	l.infos = append(l.infos, command.Command{
		Name:        name,
		Category:    command.CategoryLinks,
		Description: helpText,
		Hidden:      false,
	})

	l.output[name] = response

	l.registry.RegisterMessageCreateHandler(name, func(s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, response)
	})
}

func (l *LinkCommands) registerICELinkCommand(name, helpText, response string) {
	l.output[name] = response
	l.iceinfos = append(l.iceinfos, name)
}

func (l *LinkCommands) registerInfoDagCommand(name, helpText, response string) {
	l.infos = append(l.infos, command.Command{
		Name:        name,
		Category:    command.CategoryInfodagen,
		Description: helpText,
		Hidden:      false,
	})

	l.output[name] = response

	l.registry.RegisterMessageCreateHandler(name, func(s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, response)
	})
}

func (l *LinkCommands) slashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	reply := "I do not know that link"

	if len(i.ApplicationCommandData().Options) > 0 {
		if key, ok := i.ApplicationCommandData().Options[0].Value.(string); ok {
			if r, ok := l.output[key]; ok {
				reply = r
			}
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: reply,
		},
	})

	if err != nil {
		log.Println(err)
	}
}

// Info return the commands in this package
func (l *LinkCommands) Info() []command.Command {
	return l.infos
}
