package db

type Configuration struct {
	GuildID          string   `json:"guildID"`
	WelcomeChannelID string   `json:"welcomeChannelID"`
	WelcomeText      string   `json:"welcomeText"`
	WelcomeDM        []string `json:"welcomeDM"`

	RoleManagement RoleManagementConfiguration `json:"roleManagement"`

	Hives             []HiveConfiguration              `json:"hives"`
	LookingForPlayers []LookingForPlayersConfiguration `json:"lookingForPlayers"`
	Schedules         []ScheduleConfiguration          `json:"schedules"`
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
	RoleAdminChannelID string `json:"roleAdminChannelID"`
	DefaultRole        string `json:"defaultRole"`

	RoleSets []RoleSet `json:"roleSets"`
}

type RoleSet struct {
	Message string `json:"message"`
	Roles   []Role `json:"roles"`
}

type Role struct {
	ID    string `json:"id"`
	Emoji string `json:"emoji"`
}

type LookingForPlayersConfiguration struct {
	RequestChannelIDs  []string `json:"requestChannelIDs"`
	AdvertiseChannelID string   `json:"advertiseChannelID"`
	HiveChannelID      string   `json:"hiveChannelID"`
}

type ScheduleConfiguration struct {
	ClassName string `json:"className"`
	URL       string `json:"url"`
}
