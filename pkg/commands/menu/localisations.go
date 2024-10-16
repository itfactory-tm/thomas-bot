package menu

import (
	"math/rand"
	"time"
)

type ResponseTexts struct {
	Language         string
	LanguageCode     string
	TryLater         string
	NoWeekMenu       string
	NoWeekMenuSunday string
	NoDayMenu        string
	PoliteResponse   string
	NoItem           func(itemName string) string
}

var supportedLanguages = [2]string{"nl", "en"}

// GetResponseTexts returns the different localisation options
func GetResponseTexts(language string) (responses ResponseTexts) {
	rand.Seed(time.Now().UnixNano())
	switch language {
	case "nl":
		responses.Language = "Nederlands"
		responses.LanguageCode = "nl"
		responses.TryLater = "We kunnen momenteel niet aan het menu, probeer het later nog eens"
		responses.NoWeekMenu = "Deze campus heeft nog geen menu voor deze week!"
		responses.NoWeekMenuSunday = "Het is nog maar zondag! Nog een nachtje slapen en dan kan je weer dromen over eten op school!"
		responses.NoDayMenu = "Er is geen menu op deze dag"
		responses.PoliteResponse = "Hier is het menu: "
		responses.NoItem = func(itemName string) string { return "Er is geen " + itemName + " beschikbaar vandaag" }
		// 1 in 50 chance
		if rand.Intn(50) == 1 {
			responses.PoliteResponse = "Hier is het menu (moest je de maintainer van deze code tegenkomen, trakteer haar eens op een dagmenu): "
		}
		break

	case "en":
		responses.Language = "English"
		responses.LanguageCode = "en"
		responses.TryLater = "We can't get the menu at this time, try again later"
		responses.NoWeekMenu = "That campus does not have a menu for this week yet!"
		responses.NoWeekMenuSunday = "It's only sunday! Sleep one more night and you can dream about food at school again!"
		responses.NoDayMenu = "There is no menu available this day"
		responses.PoliteResponse = "Here is the menu: "
		responses.NoItem = func(itemName string) string { return "There is no " + itemName + " available today" }
		break

	default:
		responses.Language = "Nederlands"
		responses.LanguageCode = ""
		responses.TryLater = "We kunnen momenteel niet aan het menu, probeer het later nog eens"
		responses.NoWeekMenu = "Deze campus heeft nog geen menu voor deze week!"
		responses.NoWeekMenuSunday = "Het is nog maar zondag! Nog een nachtje slapen en dan kan je weer dromen over eten op school!"
		responses.NoDayMenu = "Er is geen menu op deze dag"
		responses.PoliteResponse = "Hier is het menu: "
		responses.NoItem = func(itemName string) string { return "Er is geen " + itemName + " beschikbaar vandaag" }
		// 1 in 50 chance
		if rand.Intn(50) == 1 {
			responses.PoliteResponse = "Hier is het menu (moest je de maintainer van deze code tegenkomen, trakteer haar eens op een dagmenu): "
		}
		break
	}

	return responses
}
