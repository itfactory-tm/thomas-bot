package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	rootCmd = &cobra.Command{
		Use:   "edward",
		Short: "edward is the companion CLI to assist with Thomas Bot tasks",
		Long: `Edward the Blue Engine (No. 2) is a 4-4-0 ex-Furness Railway K2 class locomotive. He is the first character to appear in The Railway Series. He is painted blue with red stripes. He was built in 1896 and arrived on Sodor in 1915.
Edward is one of the oldest engines on the railway, as well as very kind, always tries to help, and is a friend to everyone. He likes pulling trains as well as shunting trucks, which he is very knowledgeable about. However, the bigger engines sometimes make fun of him because they say that "' Tender Engines don't shunt'", and that he is old fashioned, but the Fat Controller still says he's a Useful Engine.
He also has a station named after him and it is also where he lives.

Oh sorry, it's just a CLI tool to do things like mass send messages.'`,
	}
)

func initConfig() {
	viper.AutomaticEnv()
}

func main() {
	flag.Parse()
	cobra.OnInitialize(initConfig)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	err := rootCmd.Execute()
	if err != nil {
		glog.Error(err)
	}
}
