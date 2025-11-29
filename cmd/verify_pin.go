package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	kra "github.com/kra-connect/go-sdk"
	"github.com/kra-connect/kra-cli/internal"
	"github.com/spf13/cobra"
)

var (
	pinBatchFile string
)

// verifyPinCmd represents the verify-pin command
var verifyPinCmd = &cobra.Command{
	Use:   "verify-pin [PIN]",
	Short: "Verify a KRA PIN number",
	Long: `Verify a KRA PIN number using the GavaConnect API.

The PIN should be in the format: P followed by 9 digits and a letter (e.g., P051234567A).

Examples:
  # Verify a single PIN
  kra-cli verify-pin P051234567A

  # Verify a PIN with JSON output
  kra-cli verify-pin P051234567A --output json

  # Verify multiple PINs from a CSV file
  kra-cli verify-pin --batch pins.csv

  # The CSV file should have a header row with a "pin" column:
  # pin
  # P051234567A
  # P059876543B

Output Formats:
  - table (default): Human-readable table format
  - json: JSON format for programmatic use
  - csv: CSV format for importing into spreadsheets`,
	Args: func(cmd *cobra.Command, args []string) error {
		if pinBatchFile == "" && len(args) != 1 {
			return fmt.Errorf("requires either a PIN argument or --batch flag")
		}
		if pinBatchFile != "" && len(args) > 0 {
			return fmt.Errorf("cannot use both PIN argument and --batch flag")
		}
		return nil
	},
	RunE: runVerifyPin,
}

func init() {
	rootCmd.AddCommand(verifyPinCmd)
	verifyPinCmd.Flags().StringVar(&pinBatchFile, "batch", "", "CSV file containing PINs to verify")
}

func runVerifyPin(cmd *cobra.Command, args []string) error {
	// Create client
	client, err := createClient()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()

	formatter := internal.NewOutputFormatter(outputFmt)

	// Batch mode
	if pinBatchFile != "" {
		return runVerifyPinBatch(ctx, client, formatter)
	}

	// Single PIN mode
	pin := args[0]

	if verbose {
		fmt.Fprintf(os.Stderr, "Verifying PIN: %s\n", pin)
	}

	result, err := client.VerifyPIN(ctx, pin)
	if err != nil {
		return fmt.Errorf("failed to verify PIN: %w", err)
	}

	return formatter.Print(result)
}

func runVerifyPinBatch(ctx context.Context, client *kra.Client, formatter *internal.OutputFormatter) error {
	// Open CSV file
	file, err := os.Open(pinBatchFile)
	if err != nil {
		return fmt.Errorf("failed to open batch file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	// Find PIN column
	pinCol := -1
	for i, col := range header {
		if col == "pin" || col == "PIN" {
			pinCol = i
			break
		}
	}

	if pinCol == -1 {
		return fmt.Errorf("CSV file must have a 'pin' or 'PIN' column")
	}

	// Read all PINs
	pins := make([]string, 0)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		if pinCol < len(record) {
			pins = append(pins, record[pinCol])
		}
	}

	if len(pins) == 0 {
		return fmt.Errorf("no PINs found in CSV file")
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Verifying %d PINs...\n", len(pins))
	}

	// Verify all PINs
	results, err := client.VerifyPINsBatch(ctx, pins)
	if err != nil {
		return fmt.Errorf("failed to verify PINs: %w", err)
	}

	if verbose {
		validCount := 0
		for _, r := range results {
			if r.IsValid {
				validCount++
			}
		}
		fmt.Fprintf(os.Stderr, "Verified %d PINs (%d valid, %d invalid)\n", len(results), validCount, len(results)-validCount)
	}

	return formatter.Print(results)
}
