package cmd

import (
	"log"

	"github.com/light-bull/lightbull/lightbull"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the control server",
	Long:  `The controller offers a webinteface (by default on port 8080) and manages the complete hardware.`,
	Run: func(cmd *cobra.Command, args []string) {
		readConfigFile()

		_, err := lightbull.New()
		if err != nil {
			log.Fatal(err)
		}

		waitQuit()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
