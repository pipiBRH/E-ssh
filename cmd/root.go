package cmd

import (
	"fmt"
	"os"

	interactive "github.com/pipiBRH/E-ssh/interactivet"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "E-ssh",
	Short: "Easy way to choose instance that you want to access through bastion.",
	Run: func(cmd *cobra.Command, args []string) {
		interactive.InteractiveWithStdin(Profile, Region, Jumper, User)
	},
}

var (
	Jumper  string
	Region  string
	Profile string
	User    string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&Profile, "profile", "p", "", "Use a specific profile from your credential file. Overrides config/env settings.")
	rootCmd.Flags().StringVarP(&Jumper, "jumper", "j", "", "Use a specific host as jumper. e.g. <host>:<port>")
	rootCmd.Flags().StringVarP(&Region, "region", "r", "", "The region to use. Overrides config/env settings.")
	rootCmd.Flags().StringVarP(&User, "user", "u", os.Getenv("USER"), "Use a specific user instead of system user")
}
