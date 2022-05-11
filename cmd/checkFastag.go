/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkFastagCmd represents the checkFastag command
var checkFastagCmd = &cobra.Command{
	Use:   "checkFastag VEHICLE_NUMBER",
	Short: "Check FASTag suffixes for vehicle registration number.",
	Run: func(cmd *cobra.Command, args []string) {
		api_key := viper.Get("RAZORPAY_LIVE_API_KEY").(string)
		if len(args) > 0 {
			fastTagSuffixes, err := readLines("data/fastag_suffixes.txt")
			if err != nil {
				log.Error().Msg("Error reading 'data/fastag_suffixes.txt'")
				os.Exit(1)
			}
			vpa := fmt.Sprintf("netc.%s", args[0])
			checkUpi(vpa, fastTagSuffixes, api_key)
		} else {
			log.Error().Msgf("❌ Please enter vehicle registration number")
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(checkFastagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkFastagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkFastagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
