package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	rootCmd.AddCommand(passwordCmd)
}

var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Get password hash for config file",
	Long:  `Outputs the password hash that needs to be written to the config file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Password: ")
		password, err := terminal.ReadPassword(0)
		if err != nil {
			log.Println(err)
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
		}

		fmt.Println()
		fmt.Println("Hash: " + string(hash))
	},
}
