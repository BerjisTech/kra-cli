// Package cmd provides the command-line interface for KRA-CLI
package cmd

import (
	"fmt"
	"os"
	"time"

	kra "github.com/BerjisTech/kra-connect-go-sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	apiKey    string
	baseURL   string
	timeout   int
	outputFmt string
	verbose   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kra-cli",
	Short: "KRA GavaConnect CLI Tool",
	Long: `KRA-CLI is a command-line tool for interacting with Kenya Revenue Authority's GavaConnect API.

It provides easy-to-use commands for:
  - PIN verification
  - TCC (Tax Compliance Certificate) checking
  - E-slip validation
  - NIL return filing
  - Taxpayer details lookup

Examples:
  kra-cli verify-pin P051234567A
  kra-cli check-tcc TCC123456
  kra-cli validate-slip 1234567890
  kra-cli file-nil-return --pin P051234567A --obligation OBL123 --period 202401

Configuration:
  Set your API key using:
    kra-cli config set api-key YOUR_API_KEY

  Or use the --api-key flag:
    kra-cli verify-pin P051234567A --api-key YOUR_API_KEY

  Or set the KRA_API_KEY environment variable:
    export KRA_API_KEY=YOUR_API_KEY`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kra-cli.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "KRA API key (overrides config)")
	rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "https://api.kra.go.ke/gavaconnect", "KRA API base URL")
	rootCmd.PersistentFlags().IntVar(&timeout, "timeout", 30, "request timeout in seconds")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "table", "output format: table, json, csv")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Bind flags to viper
	viper.BindPFlag("api_key", rootCmd.PersistentFlags().Lookup("api-key"))
	viper.BindPFlag("base_url", rootCmd.PersistentFlags().Lookup("base-url"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not find home directory: %v\n", err)
			return
		}

		// Search config in home directory with name ".kra-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kra-cli")
	}

	// Read environment variables
	viper.SetEnvPrefix("KRA")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
		}
	}

	// Apply values from config/env to flags if flags weren't explicitly set
	if !rootCmd.PersistentFlags().Changed("api-key") {
		apiKey = viper.GetString("api_key")
	}
	if !rootCmd.PersistentFlags().Changed("base-url") {
		baseURL = viper.GetString("base_url")
	}
	if !rootCmd.PersistentFlags().Changed("timeout") {
		timeout = viper.GetInt("timeout")
	}
}

// getAPIKey retrieves the API key from flags, config, or environment
func getAPIKey() (string, error) {
	if apiKey != "" {
		return apiKey, nil
	}

	key := viper.GetString("api_key")
	if key == "" {
		return "", fmt.Errorf("API key not set. Use --api-key flag, set KRA_API_KEY environment variable, or run: kra-cli config set api-key YOUR_KEY")
	}

	return key, nil
}

// createClient creates a KRA client with the configured options
func createClient() (*kra.Client, error) {
	key, err := getAPIKey()
	if err != nil {
		return nil, err
	}

	opts := []kra.Option{
		kra.WithAPIKey(key),
		kra.WithBaseURL(baseURL),
		kra.WithTimeout(time.Duration(timeout) * time.Second),
	}

	if verbose {
		opts = append(opts, kra.WithDebug(true))
	}

	return kra.NewClient(opts...)
}
