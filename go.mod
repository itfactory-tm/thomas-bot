module github.com/itfactory-tm/thomas-bot

go 1.14

require (
	github.com/bwmarrin/discordgo v0.23.3-0.20210506151729-0f05488fa0b3
	github.com/dghubble/go-twitter v0.0.0-20190719072343-39e5462e111f
	github.com/dghubble/oauth1 v0.6.0
	github.com/go-audio/wav v1.0.0
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/uuid v1.2.0 // indirect
	github.com/hraban/opus v0.0.0-20191117073431-57179dff69a6
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/meyskens/discord-ha v0.0.0-20210510092547-733c6df9e810
	github.com/meyskens/go-hcaptcha v0.0.0-20200428113538-5c28ead635cd
	github.com/sanzaru/go-giphy v0.0.0-20180211202227-c353d5ec6ee8
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.7.1
	go.mongodb.org/mongo-driver v1.4.6
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20210506145944-38f3c27a63bf // indirect
	golang.org/x/net v0.0.0-20210508051633-16afe75a6701 // indirect
	golang.org/x/sys v0.0.0-20210507161434-a76c4d0a0096 // indirect
	google.golang.org/genproto v0.0.0-20210506142907-4a47615972c2 // indirect
	google.golang.org/grpc v1.37.0 // indirect
	mvdan.cc/xurls/v2 v2.2.0
)

replace (
	// adding beta slash commands
	github.com/bwmarrin/discordgo => github.com/meyskens/discordgo v0.23.3-0.20210210083539-d11a0797e600
	// etcd fix
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
