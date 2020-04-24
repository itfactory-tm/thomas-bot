package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/itfactory-tm/thomas-bot/pkg/command"
)

func init() {
	registerInfoDagCommand("website", "Link naar Thomas More website", "Bezoek onze website: https://thomasmore.be/opleidingen/professionele-bachelor/it-factory")
	registerLinkCommand("rooster", "Link naar lessenrooster", "Bekijk hier je lessenrooster: https://rooster.thomasmore.be/")
	registerInfoDagCommand("fb", "Link naar Facebook paginas", "Bekijk hier onze facebook pagina van Toegepaste informatica: https://www.facebook.com/ToegepasteInformatica.ThomasMoreBE & ELO-ICT: https://www.facebook.com/ElektronicaICT.ThomasMoreBE & ACS: https://www.facebook.com/ACS.ThomasMoreBE")
	registerLinkCommand("canvas", "Link naar Canvas", "Bekijk hier je leerplatform (Canvas): https://thomasmore.instructure.com/")
	registerInfoDagCommand("ects", "Link naar ECTS fiches", "Bekijk hier de ECTS fiches van ELO-ICT: http://onderwijsaanbodkempen.thomasmore.be/2019/opleidingen/n/SC_51260633.html & Toegepaste Informatica: http://onderwijsaanbodkempen.thomasmore.be/opleidingen/n/SC_51260641.html")
	registerLinkCommand("lunch", "Link naar weekmenu", "Heb je honger? Bekijk hier het menu voor deze week: https://thomasmore365.sharepoint.com/sites/James/NL/stuvo/Paginas/Weekmenu.aspx?tmbaseCampus=Geel")
	registerLinkCommand("sharepoint", "Link naar Studentenportaal", "Bekijk hier de 365 sharepoint van de ITFactory: https://thomasmore365.sharepoint.com/sites/s.itfactory/SitePages/Start.aspx")
	registerLinkCommand("corona", "Link naar Corona informatie", "Zit je met vragen hoe thomasmore omgaat met corona? Bekijk dan zeker deze pagina: https://thomasmore365.sharepoint.com/sites/James/NL/stuvo/Paginas/Corona.aspx?tmbaseCampus=Geel")
	registerInfoDagCommand("stuvo", "Link naar Stuvo", "Heb je nood aan een goed gesprek? Neem dan zeker contact op met Stuvo: https://thomasmore365.sharepoint.com/sites/James/NL/stuvo/Paginas/Nood-aan-een-goed-gesprek.aspx?tmbaseCampus=Geel")
	registerInfoDagCommand("discord", "Link naar Discord documentatie", "Nog een beetje in de war over hoe Discord werkt?: https://support.discordapp.com/hc/nl")
	registerLinkCommand("kot", "Link naar kot informatie", "Informatie nodig rond op kot gaan? https://www.thomasmore.be/studenten/op-kot")
	registerInfoDagCommand("laptop", "Link naar info over laptops", "Welk materiaal heb ik nodig om in de IT-Factory te kunnen starten? https://www.thomasmore.be/sites/www.thomasmore.be/files/Laptopspecificaties%20voor%20IT%20Factory-studenten%202019-2020.pdf")
	registerLinkCommand("sinners", "Link naar Sinners", "Wat is Sinners? https://sinners.be/")
	registerInfoDagCommand("emt", "Link naar EMT", "Heeft de IT-Factory een eigen studentenvereniging? Jazeker: https://www.facebook.com/StudentenverenigingEMT")
	registerLinkCommand("wallet", "Link naar wallet", "Hoeveel staat er nog op mijn studentenkaart? https://wallet.thomasmore.be/Account/Login?ReturnUrl=%2F")
	registerLinkCommand("kuloket", "Link naar KUloket", "Kuloket raadplegen? https://kuloket.be")
	registerLinkCommand("printen", "Link naar printbeheer", "Je print gegevens bekijken? https://printbeheer.thomasmore.be/")
	registerLinkCommand("campusshop", "Link naar campusshop", "Een kijkje nemen in de campusshop? https://www.campiniamedia.be/mvc/index.jsp")
	registerLinkCommand("icecube", "Link naar ice-cube", "Ice-cube, wat is dat? https://www.thomasmore.be/ice-cube")
	registerLinkCommand("bot", "Link naar de git repo van deze bot", "Biep Boep, bekijk zeker mijn git repo https://github.com/itfactory-tm/thomas-bot")
	registerInfoDagCommand("inschrijven", "Link naar inschrijven", "Wil je je inschrijven? Dat kan hier! https://www.thomasmore.be/inschrijven")
	registerInfoDagCommand("junior", "Link naar Junior College", "Benieuwd wat Junior College is? Bekijk het hier: https://www.thomasmore.be/site/junior-university-college")
	registerInfoDagCommand("oho", "Link naar OHO", "Werken en studeren combineren? Dat kan zeker! https://www.thomasmore.be/opleidingen/professionele-bachelor/toegepaste-informatica/toegepaste-informatica-combinatie-werken-en-studeren-oho")
	registerInfoDagCommand("centen", "Link naar financiële informatie", "Wil je het financiële aspect van verder studeren bekijken? https://centenvoorstudenten.be/")
	registerInfoDagCommand("studenten", "Link naar studenten info", "Op zoek naar meer algemene info rondom verder studeren? https://www.thomasmore.be/studenten")
}

func registerLinkCommand(name, helpText, response string) {
	registerCommand(command.Command{
		Name:        name,
		Category:    command.CategoryLinks,
		Description: helpText,
		Hidden:      false,
		Handler: func(s *discordgo.Session, m *discordgo.MessageCreate) {
			s.ChannelMessageSend(m.ChannelID, response)
		},
	})
}

func registerInfoDagCommand(name, helpText, response string) {
	registerCommand(command.Command{
		Name:        name,
		Category:    command.CategoryLinks,
		Description: helpText,
		Hidden:      false,
		Handler: func(s *discordgo.Session, m *discordgo.MessageCreate) {
			s.ChannelMessageSend(m.ChannelID, response)
		},
	})
}
