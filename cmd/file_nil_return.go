package cmd

import (
	"context"
	"fmt"
	"os"

	kra "github.com/kra-connect/go-sdk"
	"github.com/kra-connect/kra-cli/internal"
	"github.com/spf13/cobra"
)

var (
	nilReturnPin        string
	nilReturnObligation string
	nilReturnPeriod     string
)

var fileNilReturnCmd = &cobra.Command{
	Use:   "file-nil-return",
	Short: "File a NIL return for a tax obligation",
	Long:  `File a NIL return for a tax obligation using the GavaConnect API.`,
	RunE:  runFileNilReturn,
}

func init() {
	rootCmd.AddCommand(fileNilReturnCmd)
	fileNilReturnCmd.Flags().StringVar(&nilReturnPin, "pin", "", "KRA PIN number (required)")
	fileNilReturnCmd.Flags().StringVar(&nilReturnObligation, "obligation", "", "Obligation ID (required)")
	fileNilReturnCmd.Flags().StringVar(&nilReturnPeriod, "period", "", "Tax period in YYYYMM format (required)")
	fileNilReturnCmd.MarkFlagRequired("pin")
	fileNilReturnCmd.MarkFlagRequired("obligation")
	fileNilReturnCmd.MarkFlagRequired("period")
}

func runFileNilReturn(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	formatter := internal.NewOutputFormatter(outputFmt)

	if verbose {
		fmt.Fprintf(os.Stderr, "Filing NIL return...\n")
		fmt.Fprintf(os.Stderr, "  PIN: %s\n", nilReturnPin)
		fmt.Fprintf(os.Stderr, "  Obligation: %s\n", nilReturnObligation)
		fmt.Fprintf(os.Stderr, "  Period: %s\n", nilReturnPeriod)
	}

	request := &kra.NILReturnRequest{
		PINNumber:    nilReturnPin,
		ObligationID: nilReturnObligation,
		Period:       nilReturnPeriod,
	}

	result, err := client.FileNILReturn(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to file NIL return: %w", err)
	}

	if verbose {
		if result.IsAccepted() {
			fmt.Fprintf(os.Stderr, "✓ NIL return accepted\n")
		} else if result.IsPending() {
			fmt.Fprintf(os.Stderr, "⏳ NIL return pending\n")
		} else if result.IsRejected() {
			fmt.Fprintf(os.Stderr, "✗ NIL return rejected\n")
		}
	}

	return formatter.Print(result)
}
