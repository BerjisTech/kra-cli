package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/kra-connect/kra-cli/internal"
	"github.com/spf13/cobra"
)

var (
	showObligations bool
)

var getTaxpayerCmd = &cobra.Command{
	Use:   "get-taxpayer [PIN]",
	Short: "Get taxpayer details",
	Long:  `Retrieve comprehensive taxpayer details using the GavaConnect API.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runGetTaxpayer,
}

func init() {
	rootCmd.AddCommand(getTaxpayerCmd)
	getTaxpayerCmd.Flags().BoolVar(&showObligations, "show-obligations", false, "Show tax obligations")
}

func runGetTaxpayer(cmd *cobra.Command, args []string) error {
	client, err := createClient()
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := context.Background()
	formatter := internal.NewOutputFormatter(outputFmt)

	pin := args[0]

	if verbose {
		fmt.Fprintf(os.Stderr, "Retrieving taxpayer details for PIN: %s\n", pin)
	}

	details, err := client.GetTaxpayerDetails(ctx, pin)
	if err != nil {
		return fmt.Errorf("failed to get taxpayer details: %w", err)
	}

	if verbose {
		fmt.Fprintf(os.Stderr, "✓ Taxpayer details retrieved\n")
		if details.IsActive() {
			fmt.Fprintf(os.Stderr, "  Status: Active\n")
		} else {
			fmt.Fprintf(os.Stderr, "  Status: Inactive\n")
		}
		fmt.Fprintf(os.Stderr, "  Display Name: %s\n", details.GetDisplayName())
	}

	if showObligations && len(details.Obligations) > 0 {
		if outputFmt == "table" {
			fmt.Println("=== Taxpayer Information ===")
			fmt.Printf("PIN: %s\n", details.PINNumber)
			fmt.Printf("Name: %s\n", details.GetDisplayName())
			if details.TaxpayerType != "" {
				fmt.Printf("Type: %s\n", details.TaxpayerType)
			}
			if details.Status != "" {
				fmt.Printf("Status: %s\n", details.Status)
			}
			fmt.Println()

			fmt.Println("=== Tax Obligations ===")
			return formatter.Print(details.Obligations)
		} else {
			return formatter.Print(details)
		}
	}

	return formatter.Print(details)
}
