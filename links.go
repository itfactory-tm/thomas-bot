package main

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	registerCommand("website", "links", "Link naar Thomas More website", sayWebsite)
	registerCommand("rooster", "links", "Link naar ", sayRooster)
	registerCommand("fb", "links", "Link naar Facebook paginas", SayFb)
	registerCommand("Canvas", "links", "Link naar Canvas", sayCanvas)
	registerCommand("ects", "links", "Link naar ECTS fiches", sayEcts)
	registerCommand("lunch", "links", "Link naar weekmenu", sayLunch)
	registerCommand("sharepoint", "links", "Link naar Studentenportaal", saySharepoint)
	registerCommand("corona", "links", "Link naar Corona informatie", sayCorona)
	registerCommand("stuvo", "links", "Link naar Stuvo", sayStuvo)
	registerCommand("discord", "links", "Link naar Discord", sayDiscord)
	registerCommand("kot", "links", "Link naar kot informatie", sayKot)
	registerCommand("centen", "links", "Link naar financiële info", sayCenten)
	registerCommand("laptop", "links", "Link naar info over laptops", sayLaptop)
	registerCommand("sinners", "links", "Link naar Sinners", saySinners)
	registerCommand("emt", "links", "Link naar EMT", sayEmt)
	registerCommand("wallet", "Link naar wallet", sayWallet)
	registerCommand("kuloket", "Link naar KUloket", sayKuloket)
	registerCommand("printen", "Link naar printbeheer", sayPrinten)
	registerCommand("campusshop", "Link naar campusshop", sayCampusshop)
	registerCommand("icecube", "Link naar ice-cube", sayIcecube)
}

func sayWebsite(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Bezoek onze website: https://thomasmore.be/opleidingen/professionele-bachelor/it-factory")
}

func sayRooster(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Bekijk hier je lessenrooster: https://rooster.thomasmore.be/")
}

func SayFb(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Bekijk hier onze facebook pagina van Toegepaste informatica: https://www.facebook.com/ToegepasteInformatica.ThomasMoreBE & ELO-ICT: https://www.facebook.com/ElektronicaICT.ThomasMoreBE & ACS: https://www.facebook.com/ACS.ThomasMoreBE")
}

func sayCanvas(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Bekijk hier je leerplatform (Canvas): https://thomasmore.instructure.com/")
}

func sayEcts(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Bekijk hier de ECTS fiches van ELO-ICT: http://onderwijsaanbodkempen.thomasmore.be/2019/opleidingen/n/SC_51260633.html & Toegepaste Informatica: http://onderwijsaanbodkempen.thomasmore.be/opleidingen/n/SC_51260641.html")
}

func sayLunch(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Heb je honger? Bekijk hier de menu voor deze week: https://thomasmore365.sharepoint.com/sites/James/NL/stuvo/Paginas/Weekmenu.aspx?tmbaseCampus=Geel")
}

func saySharepoint(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Bekijk hier de 365 sharepoint van de ITFactory: https://thomasmore365.sharepoint.com/sites/s.itfactory/SitePages/Start.aspx")

}

func sayCorona(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Zit je met vragen hoe thomasmore omgaat met corona? Bekijk dan zeker deze pagina: https://thomasmore365.sharepoint.com/sites/James/NL/stuvo/Paginas/Corona.aspx?tmbaseCampus=Geel")
}

func sayStuvo(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Heb je nood aan een goed gesprek? Neem dan zeker contact op met Stuvo: https://thomasmore365.sharepoint.com/sites/James/NL/stuvo/Paginas/Nood-aan-een-goed-gesprek.aspx?tmbaseCampus=Geel")
}

func sayDiscord(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Nog een beetje in de war over hoe Discord werkt?: https://support.discordapp.com/hc/nl")
}

func sayKot(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Informatie nodig rond op kot gaan? https://www.thomasmore.be/studenten/op-kot")
}

func sayCenten(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Wil je het financiële aspect van verder studeren bekijken? https://www.thomasmore.be/studenten/centen-voor-studenten")
}

func sayLaptop(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Welk materiaal heb ik nodig om in de IT-Factory te kunnen starten? https://www.thomasmore.be/sites/www.thomasmore.be/files/Laptopspecificaties%20voor%20IT%20Factory-studenten%202019-2020.pdf")
}

func saySinners(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Wat is Sinners? https://sinners.be/")
}

func sayEmt(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Heeft de IT-Factory een eigen studentenvereniging? Jazeker: https://www.facebook.com/StudentenverenigingEMT")
}

func sayWallet(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Hoeveel staat er nog op mijn studentenkaart? https://wallet.thomasmore.be/Account/Login?ReturnUrl=%2F")
}

func sayKuloket(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Kuloket raadplegen? https://idp.kuleuven.be/idp/profile/SAML2/Redirect/SSO?execution=e2s1")
}

func sayPrinten(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Je print gegevens bekijken? https://printbeheer.thomasmore.be/")
}

func sayCampusshop(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Een kijkje nemen in de campusshop? https://www.campiniamedia.be/mvc/index.jsp")
}
func sayIcecube(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "Ice-cube, wat is dat? https://www.thomasmore.be/ice-cube")
}