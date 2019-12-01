package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "E-ssh -j <bastion> -r <aws region> -p <profile>",
	Short: "Easy way to choose instance that you want to access through bastion.",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var (
	Jumper  string
	Region  string
	Profile string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&Profile, "profile", "p", "", "Use a specific profile from your credential file.")
	rootCmd.Flags().StringVarP(&Jumper, "jumper", "j", "", "Use a specific host as jumper.")
	rootCmd.Flags().StringVarP(&Region, "region", "r", "", "The region to use. Overrides config/env settings.")
}
