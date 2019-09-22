package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func readConfigFile() {
	// config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/lightbull/")
	viper.AddConfigPath(".")

	viper.SetDefault("api.listen", 8080)
	viper.SetDefault("api.authentication", "")

	viper.SetDefault("ethernet", "")

	viper.SetDefault("directories.config", "/lightbull")
	viper.SetDefault("directories.tmp", "/var/cache/lightbull")

	viper.SetDefault("leds.brightnessCap", 80)
	viper.SetDefault("leds.spiMHz", 1)
	viper.SetDefault("leds.fps", 25)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("Fatal error config file: %s", err))
	}
}
