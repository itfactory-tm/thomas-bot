using DSharpPlus.CommandsNext;
using DSharpPlus.CommandsNext.Attributes;
using System.Threading.Tasks;

namespace ThomasBot.Commands
{
    class UserCommands : BaseCommandModule
    {

        // Link commands ---

        [Command("website")]
        [Description("Link naar Thomas More website")]
        public async Task Website(CommandContext ctx)
        {
            await ctx.RespondAsync("Bezoek onze website: https://thomasmore.be/opleidingen/professionele-bachelor/it-factory");
        }

        [Command("website")]
        [Description("Link naar lessenrooster")]
        public async Task Rooster(CommandContext ctx)
        {
            await ctx.RespondAsync("Bekijk hier je lessenrooster: https://rooster.thomasmore.be/");
        }

        [Command("fb")]
        [Description("Link naar Facebook paginas")]
        public async Task Facebook(CommandContext ctx)
        {
            await ctx.RespondAsync("Bekijk hier onze facebook pagina van Toegepaste informatica: https://www.facebook.com/ToegepasteInformatica.ThomasMoreBE & ELO-ICT: https://www.facebook.com/ElektronicaICT.ThomasMoreBE & ACS: https://www.facebook.com/ACS.ThomasMoreBE");
        }

        [Command("canvas")]
        [Description("Link naar Canvas")]
        public async Task Canvas(CommandContext ctx)
        {
            await ctx.RespondAsync("Bekijk hier je leerplatform (Canvas): https://thomasmore.instructure.com/");
        }

        [Command("ects")]
        [Description("Link naar ECTS fiches")]
        public async Task Ects(CommandContext ctx)
        {
            await ctx.RespondAsync("Bekijk hier de ECTS fiches van ELO-ICT: http://onderwijsaanbodkempen.thomasmore.be/2019/opleidingen/n/SC_51260633.html & Toegepaste Informatica: http://onderwijsaanbodkempen.thomasmore.be/opleidingen/n/SC_51260641.html");
        }

        [Command("lunch")]
        [Description("Link naar weekmenu")]
        public async Task Lunch(CommandContext ctx)
        {
            await ctx.RespondAsync("Heb je honger? Bekijk hier het menu voor deze week: https://thomasmore365.sharepoint.com/sites/James/NL/stuvo/Paginas/Weekmenu.aspx?tmbaseCampus=Geel");
        }

        [Command("sharepoint")]
        [Description("Link naar Studentenportaal")]
        public async Task Sharepoint(CommandContext ctx)
        {
            await ctx.RespondAsync("Bekijk hier de 365 sharepoint van de ITFactory: https://thomasmore365.sharepoint.com/sites/s.itfactory/SitePages/Start.aspx");
        }

        [Command("corona")]
        [Description("Link naar Corona informatie")]
        public async Task Corona(CommandContext ctx)
        {
            await ctx.RespondAsync("Zit je met vragen hoe thomasmore omgaat met corona? Bekijk dan zeker deze pagina: https://thomasmore365.sharepoint.com/sites/James/NL/stuvo/Paginas/Corona.aspx?tmbaseCampus=Geel");
        }

        [Command("stuvo")]
        [Description("Link naar Stuvo")]
        public async Task Stuvo(CommandContext ctx)
        {
            await ctx.RespondAsync("Heb je nood aan een goed gesprek? Neem dan zeker contact op met Stuvo: https://thomasmore365.sharepoint.com/sites/James/NL/stuvo/Paginas/Nood-aan-een-goed-gesprek.aspx?tmbaseCampus=Geel");
        }

        [Command("discord")]
        [Description("Link naar Discord")]
        public async Task Discord(CommandContext ctx)
        {
            await ctx.RespondAsync("Nog een beetje in de war over hoe Discord werkt?: https://support.discordapp.com/hc/nl");
        }

        [Command("kot")]
        [Description("Link naar kot informatie")]
        public async Task Kot(CommandContext ctx)
        {
            await ctx.RespondAsync("Informatie nodig rond op kot gaan? https://www.thomasmore.be/studenten/op-kot");
        }

        [Command("centen")]
        [Description("Link naar financiele info")]
        public async Task Centen(CommandContext ctx)
        {
            await ctx.RespondAsync("Wil je het financiële aspect van verder studeren bekijken? https://centenvoorstudenten.be/");
        }

        [Command("laptop")]
        [Description("Link naar info over laptops")]
        public async Task Laptop(CommandContext ctx)
        {
            await ctx.RespondAsync("Welk materiaal heb ik nodig om in de IT-Factory te kunnen starten? https://www.thomasmore.be/sites/www.thomasmore.be/files/Laptopspecificaties%20voor%20IT%20Factory-studenten%202019-2020.pdf");
        }

        [Command("sinners")]
        [Description("Link naar Sinners")]
        public async Task Sinners(CommandContext ctx)
        {
            await ctx.RespondAsync("Wat is Sinners? https://sinners.be/");
        }

        [Command("emt")]
        [Description("Link naar EMT")]
        public async Task Emt(CommandContext ctx)
        {
            await ctx.RespondAsync("Heeft de IT-Factory een eigen studentenvereniging? Jazeker: https://www.facebook.com/StudentenverenigingEMT");
        }

        [Command("wallet")]
        [Description("Link naar wallet")]
        public async Task Wallet(CommandContext ctx)
        {
            await ctx.RespondAsync("Welk materiaal heb ik nodig om in de IT-Factory te kunnen starten? https://www.thomasmore.be/sites/www.thomasmore.be/files/Laptopspecificaties%20voor%20IT%20Factory-studenten%202019-2020.pdf");
        }

        [Command("kuloket")]
        [Description("Link naar KUloket")]
        public async Task Kuloket(CommandContext ctx)
        {
            await ctx.RespondAsync("Kuloket raadplegen? https://kuloket.be");
        }

        [Command("printen")]
        [Description("Link naar printbeheer")]
        public async Task Printen(CommandContext ctx)
        {
            await ctx.RespondAsync("Je print gegevens bekijken? https://printbeheer.thomasmore.be/");
        }

        [Command("campusshop")]
        [Description("Link naar campusshop")]
        public async Task CampusShop(CommandContext ctx)
        {
            await ctx.RespondAsync("Een kijkje nemen in de campusshop? https://www.campiniamedia.be/mvc/index.jsp");
        }

        [Command("icecube")]
        [Description("Link naar ice-cube")]
        public async Task IceCube(CommandContext ctx)
        {
            await ctx.RespondAsync("Ice-cube, wat is dat? https://www.thomasmore.be/ice-cube");
        }
        [Command("bot")]
        [Description("Link naar de git repo")]
        public async Task Bot(CommandContext ctx)
        {
            await ctx.RespondAsync("Biep Boep, bekijk zeker mijn git repo https://github.com/itfactory-tm/thomas-bot");
        }

        // Fun commands ---

        [Command("hello")]
        [Description("Zeg hallo")]
        public async Task Roll(CommandContext ctx)
        {
            await ctx.RespondAsync("Beep bop boop! Ik ben Thomas Bot, fork me on GitHub!");
        }
    }
}
