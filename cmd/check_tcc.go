package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/BerjisTech/kra-cli/internal"
	kra "github.com/BerjisTech/kra-connect-go-sdk"
	"github.com/spf13/cobra"
)

var (
	tccBatchFile string
	tccPIN       string
)

var checkTccCmd = &cobra.Command{
	Use:   "check-tcc [TCC]",
	Short: "Check a Tax Compliance Certificate",
	Long:  `Check the validity of a Tax Compliance Certificate (TCC) using the GavaConnect API.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if tccBatchFile == "" && len(args) != 1 {
			return fmt.Errorf("requires either a TCC argument or --batch flag")
		}
		if tccBatchFile != "" && len(args) > 0 {
			return fmt.Errorf("cannot use both TCC argument and --batch flag")
		}
		if tccBatchFile == "" && tccPIN == "" {
			return fmt.Errorf("--pin is required when checking a single TCC")
		}
		return nil
	},
	RunE: runCheckTcc,
}

func init() {
	rootCmd.AddCommand(checkTccCmd)
	checkTccCmd.Flags().StringVar(&tccBatchFile, "batch", "", "CSV file containing TCCs to check")
	checkTccCmd.Flags().StringVar(&tccPIN, "pin", "", "Taxpayer PIN associated with the TCC (required when not using --batch)")
}

func runCheckTcc(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	formatter := internal.NewOutputFormatter(outputFmt)

	if tccBatchFile != "" {
		return runCheckTccBatch(ctx, client, formatter)
	}

	tcc := args[0]

	if verbose {
		fmt.Fprintf(os.Stderr, "Checking TCC: %s\n", tcc)
	}

	req := &kra.TCCVerificationRequest{
		KraPIN:    tccPIN,
		TCCNumber: tcc,
	}

	result, err := client.VerifyTCC(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to check TCC: %w", err)
	}

	return formatter.Print(result)
}

func runCheckTccBatch(ctx context.Context, client *kra.Client, formatter *internal.OutputFormatter) error {
	file, err := os.Open(tccBatchFile)
	if err != nil {
		return fmt.Errorf("failed to open batch file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	tccCol := -1
	pinCol := -1
	for i, col := range header {
		switch strings.ToLower(col) {
		case "tcc":
			tccCol = i
		case "pin":
			pinCol = i
		}
	}

	if tccCol == -1 || pinCol == -1 {
		return fmt.Errorf("CSV file must have 'tcc' and 'pin' columns")
	}

	requests := make([]*kra.TCCVerificationRequest, 0)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		if tccCol < len(record) && pinCol < len(record) {
			tccValue := strings.TrimSpace(record[tccCol])
			pinValue := strings.TrimSpace(record[pinCol])
			if tccValue == "" || pinValue == "" {
				continue
			}
			requests = append(requests, &kra.TCCVerificationRequest{
				KraPIN:    pinValue,
				TCCNumber: tccValue,
			})
		}
	}

	if len(requests) == 0 {
		return fmt.Errorf("no TCC/PIN pairs found in CSV file")
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Checking %d TCCs...\n", len(requests))
	}

	results, err := client.VerifyTCCsBatch(ctx, requests)
	if err != nil {
		return fmt.Errorf("failed to check TCCs: %w", err)
	}

	if verbose {
		validCount := 0
		for _, r := range results {
			if r.IsValid {
				validCount++
			}
		}
		fmt.Fprintf(os.Stderr, "Checked %d TCCs (%d valid, %d invalid)\n", len(results), validCount, len(results)-validCount)
	}

	return formatter.Print(results)
}
