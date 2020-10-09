package main

import (
	"flag"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// to be overwritten in build
var revision = "dev"

var (
	// Used for flags.
	rootCmd = &cobra.Command{
		Use:   "james",
		Short: "james is a seperate bot for all gaming related commands in the itf server",
		Long: `james is the fifth red engine! He's the brightest red engine on Sodor, admire his beautiful paintwork!
				He's also our seperate gaming bot for all gaming related commands in the itf server, but whatever`,
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
