package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/BerjisTech/kra-cli/internal"
	kra "github.com/BerjisTech/kra-connect-go-sdk"
	"github.com/spf13/cobra"
)

var (
	nilReturnPin            string
	nilReturnObligationCode int
	nilReturnPeriod         string
	nilReturnMonth          int
	nilReturnYear           int
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
	fileNilReturnCmd.Flags().IntVar(&nilReturnObligationCode, "obligation-code", 0, "Obligation code (required)")
	fileNilReturnCmd.Flags().StringVar(&nilReturnPeriod, "period", "", "Tax period in YYYYMM format (optional if --month and --year are provided)")
	fileNilReturnCmd.Flags().IntVar(&nilReturnMonth, "month", 0, "Tax period month (1-12)")
	fileNilReturnCmd.Flags().IntVar(&nilReturnYear, "year", 0, "Tax period year (e.g. 2024)")
	fileNilReturnCmd.MarkFlagRequired("pin")
	fileNilReturnCmd.MarkFlagRequired("obligation-code")
}

func runFileNilReturn(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	formatter := internal.NewOutputFormatter(outputFmt)

	month, year, err := resolvePeriod(nilReturnPeriod, nilReturnMonth, nilReturnYear)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "Filing NIL return...\n")
		fmt.Fprintf(os.Stderr, "  PIN: %s\n", nilReturnPin)
		fmt.Fprintf(os.Stderr, "  Obligation Code: %d\n", nilReturnObligationCode)
		fmt.Fprintf(os.Stderr, "  Period: %04d-%02d\n", year, month)
	}

	request := &kra.NILReturnRequest{
		PINNumber:      nilReturnPin,
		ObligationCode: nilReturnObligationCode,
		Month:          month,
		Year:           year,
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

func resolvePeriod(period string, month, year int) (int, int, error) {
	if period != "" {
		if len(period) != 6 {
			return 0, 0, fmt.Errorf("--period must be in YYYYMM format")
		}

		pYear, err := strconv.Atoi(period[:4])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid period year: %w", err)
		}
		pMonth, err := strconv.Atoi(period[4:])
		if err != nil || pMonth < 1 || pMonth > 12 {
			return 0, 0, fmt.Errorf("invalid period month")
		}
		return pMonth, pYear, nil
	}

	if month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("month must be between 1 and 12 or provide --period")
	}
	if year < 2000 {
		return 0, 0, fmt.Errorf("year must be >= 2000 or provide --period")
	}

	return month, year, nil
}
