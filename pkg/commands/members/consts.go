package members

import (
	"regexp"
)

var userIDRoleIDRegex = *regexp.MustCompile(`<@(.*)> wants role <@&(.*)>.*`)
