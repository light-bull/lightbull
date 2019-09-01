package cmd

import (
	"fmt"
	"log"

	"github.com/light-bull/lightbull/lightbull"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the control server",
	Long:  `The controller offers a webinteface (by default on port 8080) and manages the complete hardware.`,
	Run: func(cmd *cobra.Command, args []string) {
		// config file
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/lightbull/")
		viper.AddConfigPath(".")

		viper.SetDefault("listen", 8080)

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

		// run
		_, err = lightbull.New()
		if err != nil {
			log.Fatal(err)
		}

		waitQuit()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
