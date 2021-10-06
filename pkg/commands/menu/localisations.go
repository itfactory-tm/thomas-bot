package menu

type ResponseTexts struct {
	Language       string
	LanguageCode   string
	TryLater       string
	NoWeekMenu     string
	NoDayMenu      string
	PoliteResponse string
	NoItem         func(itemName string) string
}

var supportedLanguages = [2]string{"nl", "en"}

// GetResponseTexts returns the different localisation options
func GetResponseTexts(language string) (responses ResponseTexts) {
	switch language {
	case "nl":
		responses.Language = "Nederlands"
		responses.LanguageCode = "nl"
		responses.TryLater = "We kunnen momenteel niet aan het menu, probeer het later nog eens"
		responses.NoWeekMenu = "Deze campus heeft nog geen menu voor deze week!"
		responses.NoDayMenu = "Er is geen menu op deze dag"
		responses.PoliteResponse = "Hier is het menu: "
		responses.NoItem = func(itemName string) string { return "Er is geen " + itemName + " beschikbaar vandaag" }
		break

	case "en":
		responses.Language = "English"
		responses.LanguageCode = "en"
		responses.TryLater = "We can't get the menu at this time, try again later"
		responses.NoWeekMenu = "That campus does not have a menu for this week yet!"
		responses.NoDayMenu = "There is no menu available this day"
		responses.PoliteResponse = "Here is the menu: "
		responses.NoItem = func(itemName string) string { return "There is no " + itemName + " available today" }
		break

	default:
		responses.Language = "Nederlands"
		responses.LanguageCode = ""
		responses.TryLater = "We kunnen momenteel niet aan het menu, probeer het later nog eens"
		responses.NoWeekMenu = "Deze campus heeft nog geen menu voor deze week!"
		responses.NoDayMenu = "Er is geen menu op deze dag"
		responses.PoliteResponse = "Hier is het menu: "
		responses.NoItem = func(itemName string) string { return "Er is geen " + itemName + " beschikbaar vandaag" }
		break
	}

	return responses
}
