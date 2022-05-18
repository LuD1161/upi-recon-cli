/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// checkGpayCmd represents the checkGpay command
var checkGpayCmd = &cobra.Command{
	Use:   "checkGpay EMAIL_ID",
	Short: "Check gmail id corresponding to GPay suffixes.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			gpaySuffixes, err := readLines("data/gpay_suffixes.txt")
			if err != nil {
				log.Error().Msg("Error reading 'data/gpay_suffixes.txt'")
				os.Exit(1)
			}
			vpa_suffix, email_id := args[0], args[0]
			if strings.HasSuffix(email_id, "@gmail.com") {
				vpa_suffix = email_id[:len(email_id)-10]
			}
			checkUpi(vpa_suffix, gpaySuffixes)
		} else {
			log.Error().Msgf("❌ Please enter vehicle registration number")
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(checkGpayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkGpayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkGpayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
