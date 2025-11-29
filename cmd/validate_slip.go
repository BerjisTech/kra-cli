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
	eslipBatchFile string
)

var validateSlipCmd = &cobra.Command{
	Use:   "validate-slip [ESLIP]",
	Short: "Validate an electronic payment slip",
	Long:  `Validate an electronic payment slip (e-slip) using the GavaConnect API.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if eslipBatchFile == "" && len(args) != 1 {
			return fmt.Errorf("requires either an e-slip argument or --batch flag")
		}
		if eslipBatchFile != "" && len(args) > 0 {
			return fmt.Errorf("cannot use both e-slip argument and --batch flag")
		}
		return nil
	},
	RunE: runValidateSlip,
}

func init() {
	rootCmd.AddCommand(validateSlipCmd)
	validateSlipCmd.Flags().StringVar(&eslipBatchFile, "batch", "", "CSV file containing e-slips to validate")
}

func runValidateSlip(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	formatter := internal.NewOutputFormatter(outputFmt)

	if eslipBatchFile != "" {
		return runValidateSlipBatch(ctx, client, formatter)
	}

	eslip := args[0]

	if verbose {
		fmt.Fprintf(os.Stderr, "Validating e-slip: %s\n", eslip)
	}

	result, err := client.ValidateEslip(ctx, eslip)
	if err != nil {
		return fmt.Errorf("failed to validate e-slip: %w", err)
	}

	return formatter.Print(result)
}

func runValidateSlipBatch(ctx context.Context, client *kra.Client, formatter *internal.OutputFormatter) error {
	file, err := os.Open(eslipBatchFile)
	if err != nil {
		return fmt.Errorf("failed to open batch file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	eslipCol := -1
	for i, col := range header {
		if col == "eslip" || col == "ESLIP" || col == "e-slip" || col == "E-SLIP" {
			eslipCol = i
			break
		}
	}

	if eslipCol == -1 {
		return fmt.Errorf("CSV file must have an 'eslip' or 'e-slip' column")
	}

	eslips := make([]string, 0)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		if eslipCol < len(record) {
			eslips = append(eslips, record[eslipCol])
		}
	}

	if len(eslips) == 0 {
		return fmt.Errorf("no e-slips found in CSV file")
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Validating %d e-slips...\n", len(eslips))
	}

	// Process each eslip individually (no batch method available)
	results := make([]*kra.EslipValidationResult, 0, len(eslips))
	for _, eslip := range eslips {
		result, err := client.ValidateEslip(ctx, eslip)
		if err != nil {
			if verbose {
				fmt.Fprintf(os.Stderr, "Warning: failed to validate e-slip %s: %v\n", eslip, err)
			}
			continue
		}
		results = append(results, result)
	}

	if verbose {
		validCount := 0
		for _, r := range results {
			if r.IsValid {
				validCount++
			}
		}
		fmt.Fprintf(os.Stderr, "Validated %d e-slips (%d valid, %d invalid)\n", len(results), validCount, len(results)-validCount)
	}

	return formatter.Print(results)
}
