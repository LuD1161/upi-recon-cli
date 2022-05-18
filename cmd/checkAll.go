/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// checkAllCmd represents the checkAll command
var checkAllCmd = &cobra.Command{
	Use:   "checkAll",
	Short: "Check a particular number against all UPI identifiers.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 && check_is_a_number(args[0]) {
			vpaSuffixes, err := readLines("data/all_suffixes.txt")
			if err != nil {
				log.Error().Msg("Error reading 'data/all_suffixes.txt'")
				os.Exit(1)
			}
			checkUpi(args[0], vpaSuffixes)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(checkAllCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkAllCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkAllCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
