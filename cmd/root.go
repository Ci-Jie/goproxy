package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "GoProxy",
	Short: "GoProxy is a module proxy server for golang.",
	Long:  "GoProxy is a module proxy server for golang.",
}

var log = logrus.New()
var configPath string

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(newStartCmd())
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "/etc/goproxy/goproxy.yaml", "Provide configuration file (default: /etc/goproxy/goproxy.yaml)")
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func initConfig() {
	log.Out = os.Stdout
	viper.SetConfigType("yaml")
	conf, err := os.Open(configPath)
	defer conf.Close()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	viper.ReadConfig(conf)
}
