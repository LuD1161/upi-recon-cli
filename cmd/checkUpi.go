/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkUpiCmd represents the checkUpi command
var checkUpiCmd = &cobra.Command{
	Use:   "checkUpi",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		api_key := viper.Get("RAZORPAY_LIVE_API_KEY").(string)
		if check_is_a_number(args[0]) {
			checkUpi(args[0], api_key)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkUpiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkUpiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkUpiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func checkUpi(number string, api_key string) {
	maxGoroutines := 1000
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Info().Msg("Got signal to close the program")
		os.Exit(0)
	}()

	vpaSuffixes, err := readLines("data/vpa_suffixes.txt")
	if err != nil {
		log.Error().Msg("Error reading 'data/vpa_suffixes.txt'")
		os.Exit(1)
	}

	vpas := make([]string, 0)
	for _, vpaSuffix := range vpaSuffixes {
		vpa := fmt.Sprintf("%s@%s", number, vpaSuffix)
		vpas = append(vpas, vpa)
	}
	log.Info().Msgf("Unique VPAs : %d", len(vpas))
	vpasChan := make(chan string, maxGoroutines)
	resultsChan := make(chan VPAResponse)
	for i := 0; i < maxGoroutines; i++ {
		go MakeRequest(vpasChan, resultsChan, api_key)
	}

	go func() {
		for _, vpa := range vpas {
			log.Debug().Msgf("Working on  : %s", vpa)
			vpasChan <- vpa
		}
	}()

	for i := 0; i < len(vpas); i++ {
		result := <-resultsChan
		if result.Error == nil && result.Success == true && result.CustomerName != "" {
			log.Info().Msgf("Customer Name : %s | VPA : %s", result.CustomerName, result.VPA)
		}
	}
}
