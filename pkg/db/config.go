package db

type Configuration struct {
	WelcomeChannelID string   `json:"welcomeChannelID"`
	WelcomeText      string   `json:"welcomeText"`
	WelcomeDM        []string `json:"welcomeDM"`

	RoleManagement RoleManagementConfiguration `json:"roleManagement"`

	Hives             []HiveConfiguration              `json:"hives"`
	LookingForPlayers []LookingForPlayersConfiguration `json:"lookingForPlayers"`
}

type HiveConfiguration struct {
	RequestChannelIDs  []string `json:"requestChannelIDs"`
	JunkyardCategoryID string   `json:"junkyardCategoryID"`
	TextCategoryID     string   `json:"textCategoryID"`
	VoiceCategoryID    string   `json:"voiceCategoryID"`
	Prefix             string   `json:"prefix"`
	VoiceBitrate       int      `json:"voiceBitrate"`
}

type RoleManagementConfiguration struct {
	Message            string `json:"message"`
	RoleAdminChannelID string `json:"roleAdminChannelID"`
	DefaultRole        string `json:"defaultRole"`

	Roles []Role `json:"roles"`
}

type Role struct {
	ID    string `json:"id"`
	Emoji string `json:"emoji"`
	//Name  string `json:"name"`
}

type LookingForPlayersConfiguration struct {
	RequestChannelIDs  []string `json:"requestChannelIDs"`
	AdvertiseChannelID string   `json:"advertiseChannelID"`
	HiveChannelID      string   `json:"hiveChannelID"`
}
