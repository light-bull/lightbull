package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/light-bull/lightbull/hardware"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the LEDs",
	Long:  `Run different test programs for the LED stripe`,
	Run: func(cmd *cobra.Command, args []string) {
		readConfigFile()

		hardware, err := hardware.New()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Setting color of all parts to red")
		hardware.Led.SetColorAll(255, 0, 0)
		hardware.Update()
		time.Sleep(1000000000 * 5)

		fmt.Println("Setting color of all parts to green")
		hardware.Led.SetColorAll(0, 255, 0)
		hardware.Update()
		time.Sleep(1000000000 * 5)

		fmt.Println("Setting color of all parts to blue")
		hardware.Led.SetColorAll(0, 0, 255)
		hardware.Update()
		time.Sleep(1000000000 * 5)

		fmt.Println("Setting color of all parts to white")
		hardware.Led.SetColorAll(255, 255, 255)
		hardware.Update()
		time.Sleep(1000000000 * 5)

		fmt.Println("Turn off all parts")
		hardware.Led.SetColorAll(0, 0, 0)
		hardware.Update()
		time.Sleep(1000000000 * 5)

		for _, part := range hardware.Led.GetParts() {
			fmt.Println("Settings color of", part, "to red")
			hardware.Led.SetColorAll(0, 0, 0)
			hardware.Led.SetColorAllPart(part, 255, 0, 0)
			hardware.Update()
			time.Sleep(1000000000 * 5)
		}

		fmt.Printf("\n")
	},
}
