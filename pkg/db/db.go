package db

type Database interface {
	ConfigForGuild(guildID string) (*Configuration, error)
}
