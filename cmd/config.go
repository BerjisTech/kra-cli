package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
	Long: `Manage KRA-CLI configuration settings.

Configuration is stored in ~/.kra-cli.yaml by default.

Available settings:
  - api_key: Your KRA GavaConnect API key
  - base_url: KRA API base URL
  - timeout: Request timeout in seconds
  - output: Default output format (table, json, csv)

Examples:
  # Set API key
  kra-cli config set api-key YOUR_API_KEY

  # Set base URL
  kra-cli config set base-url https://api.kra.go.ke/gavaconnect

  # Set default output format
  kra-cli config set output json

  # Set timeout
  kra-cli config set timeout 60

  # View all configuration
  kra-cli config view

  # Get a specific setting
  kra-cli config get api-key

  # Delete a setting
  kra-cli config delete api-key

  # Show config file location
  kra-cli config path`,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Available keys:
  - api-key: Your KRA GavaConnect API key
  - base-url: KRA API base URL
  - timeout: Request timeout in seconds
  - output: Default output format (table, json, csv)

Examples:
  kra-cli config set api-key YOUR_API_KEY
  kra-cli config set base-url https://api.kra.go.ke/gavaconnect
  kra-cli config set output json
  kra-cli config set timeout 60`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigSet,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Long: `Get a configuration value.

Examples:
  kra-cli config get api-key
  kra-cli config get base-url
  kra-cli config get output`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigGet,
}

var configDeleteCmd = &cobra.Command{
	Use:   "delete <key>",
	Short: "Delete a configuration value",
	Long: `Delete a configuration value.

Examples:
  kra-cli config delete api-key
  kra-cli config delete timeout`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigDelete,
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View all configuration",
	Long: `View all configuration settings.

Example:
  kra-cli config view`,
	RunE: runConfigView,
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	Long: `Show the path to the configuration file.

Example:
  kra-cli config path`,
	RunE: runConfigPath,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configDeleteCmd)
	configCmd.AddCommand(configViewCmd)
	configCmd.AddCommand(configPathCmd)
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	// Convert hyphenated keys to underscored for viper
	viperKey := convertKeyToViperFormat(key)

	// Validate key
	validKeys := map[string]bool{
		"api_key":  true,
		"base_url": true,
		"timeout":  true,
		"output":   true,
	}

	if !validKeys[viperKey] {
		return fmt.Errorf("invalid configuration key: %s (valid keys: api-key, base-url, timeout, output)", key)
	}

	// Set the value
	viper.Set(viperKey, value)

	// Ensure config directory exists
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	configPath := filepath.Join(home, ".kra-cli.yaml")

	// Write config file
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("✓ Configuration updated: %s = %s\n", key, value)
	fmt.Printf("  Config file: %s\n", configPath)

	return nil
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]
	viperKey := convertKeyToViperFormat(key)

	value := viper.Get(viperKey)
	if value == nil {
		fmt.Printf("%s is not set\n", key)
		return nil
	}

	fmt.Printf("%s: %v\n", key, value)
	return nil
}

func runConfigDelete(cmd *cobra.Command, args []string) error {
	key := args[0]
	viperKey := convertKeyToViperFormat(key)

	// Check if key exists
	if !viper.IsSet(viperKey) {
		fmt.Printf("%s is not set\n", key)
		return nil
	}

	// Get all settings
	settings := viper.AllSettings()

	// Delete the key
	delete(settings, viperKey)

	// Clear viper
	viper.Reset()

	// Re-set all settings except the deleted one
	for k, v := range settings {
		viper.Set(k, v)
	}

	// Write config file
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	configPath := filepath.Join(home, ".kra-cli.yaml")

	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("✓ Configuration deleted: %s\n", key)
	return nil
}

func runConfigView(cmd *cobra.Command, args []string) error {
	settings := viper.AllSettings()

	if len(settings) == 0 {
		fmt.Println("No configuration set")
		return nil
	}

	fmt.Println("Configuration:")
	for key, value := range settings {
		// Mask API key for security
		if key == "api_key" {
			valueStr := fmt.Sprintf("%v", value)
			if len(valueStr) > 8 {
				fmt.Printf("  %s: %s...%s\n", key, valueStr[:4], valueStr[len(valueStr)-4:])
			} else {
				fmt.Printf("  %s: %v\n", key, value)
			}
		} else {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	return nil
}

func runConfigPath(cmd *cobra.Command, args []string) error {
	// Check if config file exists
	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		fmt.Printf("Config file: %s\n", configFile)
		return nil
	}

	// Show default path
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %w", err)
	}

	defaultPath := filepath.Join(home, ".kra-cli.yaml")
	fmt.Printf("Config file (default): %s\n", defaultPath)

	// Check if file exists
	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		fmt.Println("  (file does not exist yet)")
	} else {
		fmt.Println("  (file exists)")
	}

	return nil
}

// convertKeyToViperFormat converts hyphenated keys to underscored format
// e.g., "api-key" -> "api_key"
func convertKeyToViperFormat(key string) string {
	switch key {
	case "api-key":
		return "api_key"
	case "base-url":
		return "base_url"
	default:
		return key
	}
}
