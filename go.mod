module github.com/itfactory-tm/thomas-bot

go 1.14

require (
	github.com/arran4/golang-ical v0.0.0-20210807024147-770fa87aff1d
	github.com/bwmarrin/discordgo v0.23.2
	github.com/dghubble/go-twitter v0.0.0-20190719072343-39e5462e111f
	github.com/dghubble/oauth1 v0.6.0
	github.com/go-audio/wav v1.0.0
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/hraban/opus v0.0.0-20191117073431-57179dff69a6
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/meyskens/discord-ha v0.0.0-20210723094030-8791e408bab7
	github.com/meyskens/go-hcaptcha v0.0.0-20200428113538-5c28ead635cd
	github.com/sanzaru/go-giphy v0.0.0-20211118160211-e9e78e55bc7a
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	go.mongodb.org/mongo-driver v1.5.1
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20210506142907-4a47615972c2 // indirect
	mvdan.cc/xurls/v2 v2.2.0
)

replace (
	// pull select code
	github.com/bwmarrin/discordgo v0.23.2 => github.com/meyskens/discordgo v0.23.3-0.20210723093830-80a9f1364942
	// etcd fix
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
