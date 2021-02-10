package db

type Configuration struct {
	WelcomeChannelID string              `json:"welcomeChannelID"`
	WelcomeText      string              `json:"welcomeText"`
	Hives            []HiveConfiguration `json:"hives"`
}

type HiveConfiguration struct {
	RequestChannelIDs  []string `json:"requestChannelIDs"`
	JunkyardCategoryID string   `json:"junkyardCategoryID"`
	TextCategoryID     string   `json:"textCategoryID"`
	VoiceCategoryID    string   `json:"voiceCategoryID"`
	Prefix             string   `json:"prefix"`
	VoiceBitrate       int      `json:"voiceBitrate"`
}
