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
		Use:   "thomas",
		Short: "thomas is the main server binary for Thomas Bot",
		Long:  "thomas is the main server binary for Thomas Bot",
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
