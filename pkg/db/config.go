package db

type Configuration struct {
	WelcomeChannelID string                       `json:"welcomeChannelID"`
	WelcomeText      string                       `json:"welcomeText"`
	Hives            map[string]HiveConfiguration `json:"hives"`
}

type HiveConfiguration struct {
	RequestChannelID   string `json:"requestChannelID"`
	JunkyardCategoryID string `json:"junkyardCategoryID"`
	TextCategoryID     string `json:"textCategoryID"`
	VoiceCategoryID    string `json:"voiceCategoryID"`
	Prefix             string `json:"prefix"`
}
