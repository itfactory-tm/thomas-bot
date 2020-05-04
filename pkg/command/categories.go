package command

import (
	"strings"
)

// Category is an enumeration type of the categories
type Category int

const (
	// CategoryInfodagen is the category for Infodagen commands
	CategoryInfodagen Category = iota
	// CategoryLinks is the category for link commands
	CategoryLinks Category = iota
	// CategoryAlgemeen is the category for general commands
	CategoryAlgemeen Category = iota
	// CategoryFun is the category for funny commands
	CategoryFun Category = iota
	// CategoryModeratie is the category for moderation commands
	CategoryModeratie Category = iota
	// CategoryStudenten is the category for student commands
	CategoryStudenten Category = iota
	// CategoryOverige is the category for other commands
	CategoryOverige Category = iota
)

// CategoryToString converts a Category to the string of the name
func CategoryToString(in Category) string {
	switch in {
	case CategoryFun:
		return "Fun"
	case CategoryLinks:
		return "Links"
	case CategoryInfodagen:
		return "Infodagen"
	case CategoryAlgemeen:
		return "Algemeen"
	case CategoryStudenten:
		return "Studenten"
	case CategoryModeratie:
		return "Moderatie"
	case CategoryOverige:
		return "Overige"
	}

	return "" // empty is not found
}

// StringToCategory gets the Category of string, case insensitive
func StringToCategory(in string) Category {
	switch strings.ToLower(in) {
	case "fun":
		return CategoryFun
	case "links":
		return CategoryLinks
	case "infodagen":
		return CategoryInfodagen
	case "algemeen":
		return CategoryAlgemeen
	case "studenten":
		return CategoryStudenten
	case "moderatie":
		return CategoryModeratie
	case "overige":
		return CategoryOverige
	}

	return CategoryOverige // Overige is not found
}
