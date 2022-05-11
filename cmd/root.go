/*
Copyright © 2022 Aseem Shrey

*/
package cmd

import (
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var (
	threads     int
	timeout     int
	cfgFile     string
	userLicense string
	rootCmd     = &cobra.Command{
		Use:   "upi-recon-cli PHONE_NUMBER",
		Args:  cobra.ArbitraryArgs, // https://github.com/spf13/cobra/issues/42
		Short: "Check UPI ids corresponding to a mobile number",
		Long: `Check virtual payment address corresponding to a mobile number.
Get the user's name as well.`,
		Run: func(cmd *cobra.Command, args []string) {
			api_key := viper.Get("RAZORPAY_LIVE_API_KEY").(string)
			if len(args) > 0 && check_is_a_number(args[0]) {
				vpaSuffixes, err := readLines("data/vpa_suffixes.txt")
				if err != nil {
					log.Error().Msg("Error reading 'data/vpa_suffixes.txt'")
					os.Exit(1)
				}
				checkUpi(args[0], vpaSuffixes, api_key)
			} else {
				cmd.Help()
			}
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file")
	rootCmd.PersistentFlags().IntVarP(&threads, "threads", "t", 100, "No of threads")
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "", 15, "Timeout for requests")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Debug().Msgf("Using config file: %s", viper.ConfigFileUsed())
	}
}
