module github.com/itfactory-tm/thomas-bot

go 1.14

require (
	github.com/bwmarrin/discordgo v0.23.3-0.20210314162722-182d9b48f34b
	github.com/davecgh/go-spew v1.1.1
	github.com/dghubble/go-twitter v0.0.0-20190719072343-39e5462e111f
	github.com/dghubble/oauth1 v0.6.0
	github.com/go-audio/wav v1.0.0
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/uuid v1.2.0 // indirect
	github.com/hraban/opus v0.0.0-20191117073431-57179dff69a6
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/meyskens/discord-ha v0.0.0-20210315192353-c63c44a23a77
	github.com/meyskens/go-hcaptcha v0.0.0-20200428113538-5c28ead635cd
	github.com/sanzaru/go-giphy v0.0.0-20180211202227-c353d5ec6ee8
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.7.1
	go.mongodb.org/mongo-driver v1.4.6
	golang.org/x/crypto v0.0.0-20210314154223-e6e6c4f2bb5b // indirect
	golang.org/x/net v0.0.0-20210315170653-34ac3e1c2000 // indirect
	golang.org/x/sys v0.0.0-20210315160823-c6e025ad8005 // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/genproto v0.0.0-20210315173758-2651cd453018 // indirect
	google.golang.org/grpc v1.36.0 // indirect
	mvdan.cc/xurls/v2 v2.2.0
)

// etcd fix
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
