package command

import (
	"strings"
)

type Category int

const (
	CategoryFun       Category = iota
	CategoryLinks     Category = iota
	CategoryInfodagen Category = iota
	CategoryInfo      Category = iota
	CategoryStudenten Category = iota
	CategoryModeratie Category = iota
	CategoryOverige   Category = iota
)

func CategoryToString(in Category) string {
	switch in {
	case CategoryFun:
		return "Fun"
	case CategoryLinks:
		return "Links"
	case CategoryInfodagen:
		return "Infodagen"
	case CategoryInfo:
		return "Info"
	case CategoryStudenten:
		return "Studenten"
	case CategoryModeratie:
		return "Moderatie"
	case CategoryOverige:
		return "Overige"
	}

	return "" // empty is not found
}

func StringToCategory(in string) Category {
	switch strings.ToLower(in) {
	case "fun":
		return CategoryFun
	case "links":
		return CategoryLinks
	case "infodagen":
		return CategoryInfodagen
	case "info":
		return CategoryInfo
	case "studenten":
		return CategoryStudenten
	case "moderatie":
		return CategoryModeratie
	case "overige":
		return CategoryOverige
	}

	return CategoryOverige // Overige is not found
}
