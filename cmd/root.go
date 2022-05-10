/*
Copyright © 2022 Aseem Shrey

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var (
	cfgFile     string
	userLicense string
	rootCmd     = &cobra.Command{
		Use:   "upi-recon-cli",
		Args:  cobra.ArbitraryArgs, // https://github.com/spf13/cobra/issues/42
		Short: "Check UPI ids corresponding to a mobile number",
		Long: `Check virtual payment address corresponding to a mobile number.
Get the user's name as well.`,
		Run: func(cmd *cobra.Command, args []string) {
			api_key := viper.Get("RAZORPAY_LIVE_API_KEY").(string)
			if check_is_a_number(args[0]) {
				checkUpi(args[0], api_key)
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().Int32P("threads", "t", 1000, "No of threads")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}