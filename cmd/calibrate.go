package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/light-bull/lightbull/hardware"
	"github.com/spf13/cobra"
)

var numberLeds int

func init() {
	rootCmd.AddCommand(calibrateCmd)

	calibrateCmd.Flags().IntVarP(&numberLeds, "number", "n", 750, "Total number of connected LEDs")
}

var calibrateCmd = &cobra.Command{
	Use:   "calibrate",
	Short: "Turn single LEDs on",
	Long: `Interactively turn single LEDs on to find out which LED is in which LED part.
	The part definition in ignored here and the original/raw LED IDs are used.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Init LEDs
		led := hardware.NewLED()
		led.AddPart("calibrate", 0, numberLeds)
		if err := led.Init(); err != nil {
			log.Fatal(err)
		}

		// Read number for led
		reader := bufio.NewReader(os.Stdin)

		for {
			fmt.Print("Enter LED ID: ")
			// read input, trim newline and convert to int
			idStr, _ := reader.ReadString('\n')
			idStr = strings.TrimSuffix(idStr, "\n")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("Invalid number")
				continue
			}

			// check range
			if id < 0 || id >= numberLeds {
				fmt.Println("ID out of range")
				continue
			}

			// turn single LED on
			led.SetColorAll(0, 0, 0)
			led.SetColor("calibrate", id, 255, 0, 0)
			led.Update()
		}
	},
}
