package links

import (
	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

// LinkCommands contains the tm!hello command
type LinkCommands struct {
	infos    []command.Command
	registry command.Registry
}

// NewLinkCommands gives a new LinkCommands
func NewLinkCommands() *LinkCommands {
	return &LinkCommands{
		infos: []command.Command{},
	}
}

// Register registers the handlers
func (l *LinkCommands) Register(registry command.Registry, server command.Server) {
	l.registry = registry

	l.buildLinks()
}

func (l *LinkCommands) buildLinks() {
	l.registerInfoDagCommand("website", "Link naar Thomas More website", "Bezoek onze website: https://thomasmore.be/opleidingen/professionele-bachelor/it-factory")
	l.registerLinkCommand("rooster", "Link naar lessenrooster", "Bekijk hier je lessenrooster: https://rooster.thomasmore.be/")
	l.registerInfoDagCommand("fb", "Link naar Facebook paginas", "Bekijk hier onze facebook pagina van Toegepaste informatica: https://www.facebook.com/ToegepasteInformatica.ThomasMoreBE & ELO-ICT: https://www.facebook.com/ElektronicaICT.ThomasMoreBE & ACS: https://www.facebook.com/ACS.ThomasMoreBE")
	l.registerLinkCommand("canvas", "Link naar Canvas", "Bekijk hier je leerplatform (Canvas): https://thomasmore.instructure.com/")
	l.registerInfoDagCommand("ects", "Link naar ECTS fiches", "Bekijk hier de ECTS fiches van ELO-ICT: http://onderwijsaanbodkempen.thomasmore.be/2019/opleidingen/n/SC_51260633.html & Toegepaste Informatica: http://onderwijsaanbodkempen.thomasmore.be/opleidingen/n/SC_51260641.html")
	l.registerLinkCommand("lunch", "Link naar weekmenu", "Heb je honger? Bekijk hier het menu voor deze week: https://thomasmore365.sharepoint.com/sites/James/NL/stuvo/Paginas/Weekmenu.aspx?tmbaseCampus=Geel")
	l.registerLinkCommand("sharepoint", "Link naar Studentenportaal", "Bekijk hier de 365 sharepoint van de ITFactory: https://thomasmore365.sharepoint.com/sites/s.itfactory/SitePages/Start.aspx")
	l.registerLinkCommand("corona", "Link naar Corona informatie", "Zit je met vragen hoe thomasmore omgaat met corona? Bekijk dan zeker deze pagina: https://thomasmore365.sharepoint.com/sites/James/NL/stuvo/Paginas/Corona.aspx?tmbaseCampus=Geel")
	l.registerInfoDagCommand("stuvo", "Link naar Stuvo", "Heb je nood aan een goed gesprek? Neem dan zeker contact op met Stuvo: https://thomasmore365.sharepoint.com/sites/James/NL/stuvo/Paginas/Nood-aan-een-goed-gesprek.aspx?tmbaseCampus=Geel")
	l.registerInfoDagCommand("discord", "Link naar Discord documentatie", "Nog een beetje in de war over hoe Discord werkt?: https://support.discordapp.com/hc/nl")
	l.registerLinkCommand("kot", "Link naar kot informatie", "Informatie nodig rond op kot gaan? https://www.thomasmore.be/studenten/op-kot")
	l.registerInfoDagCommand("laptop", "Link naar info over laptops", "Welk materiaal heb ik nodig om in de IT-Factory te kunnen starten? https://www.thomasmore.be/sites/www.thomasmore.be/files/Laptopspecificaties%20voor%20IT%20Factory-studenten%202019-2020.pdf")
	l.registerLinkCommand("sinners", "Link naar Sinners", "Wat is Sinners? https://sinners.be/")
	l.registerInfoDagCommand("emt", "Link naar EMT", "Heeft de IT-Factory een eigen studentenvereniging? Jazeker: https://www.facebook.com/StudentenverenigingEMT")
	l.registerLinkCommand("wallet", "Link naar wallet", "Hoeveel staat er nog op mijn studentenkaart? https://wallet.thomasmore.be/Account/Login?ReturnUrl=%2F")
	l.registerLinkCommand("kuloket", "Link naar KUloket", "Kuloket raadplegen? https://kuloket.be")
	l.registerLinkCommand("printen", "Link naar printbeheer", "Je print gegevens bekijken? https://printbeheer.thomasmore.be/")
	l.registerLinkCommand("campusshop", "Link naar campusshop", "Een kijkje nemen in de campusshop? https://www.campiniamedia.be/mvc/index.jsp")
	l.registerLinkCommand("icecube", "Link naar ice-cube", "Ice-cube, wat is dat? https://www.thomasmore.be/ice-cube")
	l.registerLinkCommand("bot", "Link naar de git repo van deze bot", "Biep Boep, bekijk zeker mijn git repo https://github.com/itfactory-tm/thomas-bot")
	l.registerInfoDagCommand("inschrijven", "Link naar inschrijven", "Wil je je inschrijven? Dat kan hier! https://www.thomasmore.be/inschrijven")
	l.registerInfoDagCommand("junior", "Link naar Junior College", "Benieuwd wat Junior College is? Bekijk het hier: https://www.thomasmore.be/site/junior-university-college")
	l.registerInfoDagCommand("oho", "Link naar OHO", "Werken en studeren combineren? Dat kan zeker! https://www.thomasmore.be/opleidingen/professionele-bachelor/toegepaste-informatica/toegepaste-informatica-combinatie-werken-en-studeren-oho")
	l.registerInfoDagCommand("centen", "Link naar financiële informatie", "Wil je het financiële aspect van verder studeren bekijken? https://centenvoorstudenten.be/")
	l.registerInfoDagCommand("studenten", "Link naar studenten info", "Op zoek naar meer algemene info rondom verder studeren? https://www.thomasmore.be/studenten")
	l.registerLinkCommand("template", "Link naar TM huisstijl templates", "Hier vind je de TM huisstijl templates: https://thomasmore365.sharepoint.com/sites/James/NL/marcom/Paginas/Huisstijl-templates.aspx")
	l.registerLinkCommand("onlineexamen", "Link naar info over online examens", "Instructies voor het online examen vind je hier: https://thomasmore365.sharepoint.com/sites/James/NL/ict/Paginas/Instructies-bij-digitaal-examineren.aspx\nFAQ over digitaal examineren: https://thomasmore365.sharepoint.com/sites/James/NL/ict/Paginas/FAQ-digitaal-examineren.aspx\nTIP: bekijk zeker ook de canvas cursus voor vak specifieke info")
}

func (l *LinkCommands) registerLinkCommand(name, helpText, response string) {
	l.infos = append(l.infos, command.Command{
		Name:        name,
		Category:    command.CategoryLinks,
		Description: helpText,
		Hidden:      false,
	})

	l.registry.RegisterMessageCreateHandler(name, func(s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, response)
	})
}

func (l *LinkCommands) registerInfoDagCommand(name, helpText, response string) {
	l.infos = append(l.infos, command.Command{
		Name:        name,
		Category:    command.CategoryInfodagen,
		Description: helpText,
		Hidden:      false,
	})

	l.registry.RegisterMessageCreateHandler(name, func(s *discordgo.Session, m *discordgo.MessageCreate) {
		s.ChannelMessageSend(m.ChannelID, response)
	})
}

// Info return the commands in this package
func (l *LinkCommands) Info() []command.Command {
	return l.infos
}
