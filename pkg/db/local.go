package db

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type LocalDatabase struct {
	config map[string]Configuration
}

func (l *LocalDatabase) ConfigForGuild(guildID string) (*Configuration, error) {
	config, ok := l.config[guildID]
	if !ok {
		return nil, errors.New("guild not in database")
	}
	return &config, nil
}

func NewLocalDB(path string) (Database, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var conf map[string]Configuration
	err = json.Unmarshal(b, &conf)
	if err != nil {
		return nil, err
	}

	return &LocalDatabase{
		config: conf,
	}, nil
}
