package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "broker-gateway",
	Short: "Matthiola",
	Long:  "The broker gateway of project Matthiola, a distributed commodities OTC electronic trading system",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("[cmd.root.rootCmd] [FETAL] %s", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config/config.toml", "config file")
	rootCmd.AddCommand(receiverCmd)
	rootCmd.AddCommand(senderCmd)
	rootCmd.AddCommand(matcherCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(broadcasterCmd)

	viper.BindPFlags(rootCmd.PersistentFlags())
	viper.SetDefault("author", "rudeigerc <rudeigerc@gmail.com>")
	viper.SetDefault("license", "MIT")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("config")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("[cmd.root.initConfig] [FETAL] Fatal error config file: %s \n", err)
	}

}
