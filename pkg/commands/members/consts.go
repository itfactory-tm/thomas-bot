package members

import (
	"regexp"

	"github.com/elliotchance/orderedmap"
)

// TODO: replace me
const roleChannelID = "739512119338467449" // one day we need to stop doing these...

const roleMessage = `We need to assign you a role inside our Discord which will help you gain access to the class specific channels.
Select the following emoji(s) for roles you want to request, note that our moderation team has to approve these first.
1️⃣: 1ITF Student
2️⃣: 2ITF Student
3️⃣: 3ITF Student
👩‍💻: OHO Student
👩‍🎓: Alumni
👩‍🏫: Teacher`

var userIDRoleIDRegex = *regexp.MustCompile(`<@(.*)> wants role <@&(.*)>.*`)

var roleEmoji = orderedmap.NewOrderedMap()

func init() {

	// very upset Discord does not support non-binary emoji
	roleEmoji.Set("1️⃣", "687567949795557386") // 1ITF
	roleEmoji.Set("2️⃣", "687568334379679771") // 2ITF
	roleEmoji.Set("3️⃣", "687568470820388864") // 3ITF
	roleEmoji.Set("👩‍💻", "689844328528478262") // OHO
	roleEmoji.Set("👩‍🎓", "688368287255494702") // Alumni
	roleEmoji.Set("👩‍🏫", "687567374198767617") // Teacher
}
