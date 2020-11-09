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
		Use:   "alf",
		Short: "alf is the cleaner server binary for Thomas Bot",
		Long: `alf is the cleaner server binary for Thomas Bot
		Bert and Alf are two cleaners who washed Gordon down in Leaves after he derailed from the turntable in a ditch at Vicarstown Sheds.`,
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
