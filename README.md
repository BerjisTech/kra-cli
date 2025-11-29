# KRA-CLI

> Command-line interface for Kenya Revenue Authority's GavaConnect API

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue.svg)](https://golang.org)

## Overview

KRA-CLI is a powerful command-line tool that provides easy access to the Kenya Revenue Authority's GavaConnect API for tax compliance verification, PIN validation, TCC checking, and more.

## Features

- ✅ **PIN Verification** - Verify KRA PIN numbers with detailed taxpayer information
- ✅ **TCC Checking** - Validate Tax Compliance Certificates
- ✅ **E-slip Validation** - Verify electronic payment slips
- ✅ **NIL Return Filing** - File NIL returns programmatically
- ✅ **Taxpayer Details** - Retrieve comprehensive taxpayer information
- ✅ **Batch Operations** - Process multiple requests from CSV files
- ✅ **Multiple Output Formats** - Table, JSON, and CSV outputs
- ✅ **Configuration Management** - Store API keys and settings
- ✅ **Cross-Platform** - Works on Windows, macOS, and Linux

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/kra-connect/kra-cli.git
cd kra-cli

# Build
go build -o kra-cli

# Install
go install
```

### From Release (Coming Soon)

Download pre-built binaries from the [releases page](https://github.com/kra-connect/kra-cli/releases).

**macOS:**
```bash
brew install kra-cli
```

**Linux:**
```bash
# Debian/Ubuntu
wget https://github.com/kra-connect/kra-cli/releases/download/v0.1.0/kra-cli_0.1.0_linux_amd64.deb
sudo dpkg -i kra-cli_0.1.0_linux_amd64.deb

# Red Hat/CentOS
wget https://github.com/kra-connect/kra-cli/releases/download/v0.1.0/kra-cli_0.1.0_linux_amd64.rpm
sudo rpm -i kra-cli_0.1.0_linux_amd64.rpm
```

**Windows:**
Download the `.exe` file from releases and add to your PATH.

## Quick Start

### 1. Set Up Your API Key

```bash
# Option 1: Using config command
kra-cli config set api-key YOUR_API_KEY

# Option 2: Using environment variable
export KRA_API_KEY=YOUR_API_KEY

# Option 3: Using flag (not recommended for security)
kra-cli verify-pin P051234567A --api-key YOUR_API_KEY
```

### 2. Verify a PIN

```bash
kra-cli verify-pin P051234567A
```

### 3. Check a TCC

```bash
kra-cli check-tcc TCC123456
```

## Usage

### Global Flags

All commands support these global flags:

```
--api-key string    KRA API key (overrides config)
--base-url string   KRA API base URL (default "https://api.kra.go.ke/gavaconnect")
--config string     Config file (default is $HOME/.kra-cli.yaml)
--output, -o        Output format: table, json, csv (default "table")
--timeout int       Request timeout in seconds (default 30)
--verbose, -v       Verbose output
--help, -h          Help for any command
```

### PIN Verification

Verify a single PIN or multiple PINs from a CSV file.

```bash
# Single PIN
kra-cli verify-pin P051234567A

# With JSON output
kra-cli verify-pin P051234567A --output json

# Batch verification from CSV
kra-cli verify-pin --batch pins.csv

# Batch with CSV output
kra-cli verify-pin --batch pins.csv --output csv > results.csv
```

**CSV Format for Batch:**
```csv
pin
P051234567A
P059876543B
P051111111C
```

**Example Output (Table):**
```
PIN           VALID  NAME                TYPE        STATUS
P051234567A   true   John Doe Ltd        Company     active
P059876543B   true   Jane Smith          Individual  active
P051111111C   false  -                   -           -
```

### TCC Checking

Verify Tax Compliance Certificates.

```bash
# Single TCC
kra-cli check-tcc TCC123456

# Batch checking
kra-cli check-tcc --batch tccs.csv

# JSON output
kra-cli check-tcc TCC123456 --output json
```

**CSV Format:**
```csv
tcc
TCC123456
TCC789012
```

### E-slip Validation

Validate electronic payment slips.

```bash
# Single e-slip
kra-cli validate-slip 1234567890

# Batch validation
kra-cli validate-slip --batch eslips.csv

# CSV output
kra-cli validate-slip 1234567890 --output csv
```

**CSV Format:**
```csv
eslip
1234567890
0987654321
```

### NIL Return Filing

File NIL returns for tax obligations.

```bash
# File a NIL return
kra-cli file-nil-return \
  --pin P051234567A \
  --obligation OBL123456 \
  --period 202401

# With verbose output
kra-cli file-nil-return \
  --pin P051234567A \
  --obligation OBL123456 \
  --period 202401 \
  --verbose

# JSON output for automation
kra-cli file-nil-return \
  --pin P051234567A \
  --obligation OBL123456 \
  --period 202401 \
  --output json
```

### Taxpayer Details

Retrieve comprehensive taxpayer information.

```bash
# Get taxpayer details
kra-cli get-taxpayer P051234567A

# JSON output
kra-cli get-taxpayer P051234567A --output json

# Show all obligations
kra-cli get-taxpayer P051234567A --show-obligations
```

### Configuration Management

Manage CLI configuration settings.

```bash
# Set API key
kra-cli config set api-key YOUR_API_KEY

# Set base URL
kra-cli config set base-url https://api.kra.go.ke/gavaconnect

# Set default output format
kra-cli config set output json

# View current configuration
kra-cli config view

# Get a specific setting
kra-cli config get api-key

# Delete a setting
kra-cli config delete api-key

# Show config file location
kra-cli config path
```

## Output Formats

### Table Format (Default)

Human-readable table output for terminal use.

```bash
kra-cli verify-pin P051234567A
# or
kra-cli verify-pin P051234567A --output table
```

### JSON Format

Structured JSON for programmatic use and piping.

```bash
kra-cli verify-pin P051234567A --output json | jq '.taxpayer_name'
```

**Example JSON Output:**
```json
{
  "pin_number": "P051234567A",
  "is_valid": true,
  "taxpayer_name": "John Doe Ltd",
  "taxpayer_type": "Company",
  "status": "active",
  "registration_date": "2020-01-15",
  "tax_office": "Nairobi Tax Office",
  "verified_at": "2025-01-28T10:30:00Z"
}
```

### CSV Format

CSV output for importing into spreadsheets.

```bash
kra-cli verify-pin --batch pins.csv --output csv > results.csv
```

## Configuration File

KRA-CLI stores configuration in `~/.kra-cli.yaml`:

```yaml
api_key: your-api-key-here
base_url: https://api.kra.go.ke/gavaconnect
timeout: 30
output: table
```

You can also specify a custom config file:

```bash
kra-cli --config /path/to/config.yaml verify-pin P051234567A
```

## Environment Variables

KRA-CLI supports these environment variables:

- `KRA_API_KEY` - Your GavaConnect API key
- `KRA_BASE_URL` - API base URL (optional)
- `KRA_TIMEOUT` - Request timeout in seconds (optional)

Environment variables are overridden by config file settings, which are overridden by command-line flags.

**Precedence (highest to lowest):**
1. Command-line flags (`--api-key`)
2. Environment variables (`KRA_API_KEY`)
3. Config file (`~/.kra-cli.yaml`)

## Batch Operations

### CSV File Format

All batch operations expect CSV files with a header row and appropriate column names.

**PIN Verification (`pins.csv`):**
```csv
pin
P051234567A
P059876543B
```

**TCC Checking (`tccs.csv`):**
```csv
tcc
TCC123456
TCC789012
```

**E-slip Validation (`eslips.csv`):**
```csv
eslip
1234567890
0987654321
```

### Processing Large Batches

For large CSV files, use verbose mode to track progress:

```bash
kra-cli verify-pin --batch large-pins.csv --verbose
```

Output results to a file:

```bash
kra-cli verify-pin --batch pins.csv --output csv > results.csv
kra-cli verify-pin --batch pins.csv --output json > results.json
```

## Examples

### Example 1: Verify Suppliers

Verify a list of supplier PINs and save results:

```bash
# Create suppliers.csv with PIN column
cat > suppliers.csv <<EOF
pin
P051234567A
P059876543B
P051111111C
EOF

# Verify all PINs
kra-cli verify-pin --batch suppliers.csv --output csv > verified-suppliers.csv

# Count valid suppliers
kra-cli verify-pin --batch suppliers.csv --output json | jq '[.[] | select(.is_valid == true)] | length'
```

### Example 2: Check Expiring TCCs

Check TCCs and filter expiring ones:

```bash
kra-cli check-tcc --batch tccs.csv --output json | \
  jq '.[] | select(.days_until_expiry < 30 and .days_until_expiry > 0)'
```

### Example 3: Automation with Scripts

```bash
#!/bin/bash
# verify-and-notify.sh

PINS_FILE="suppliers.csv"
OUTPUT_FILE="results.json"

# Verify PINs
kra-cli verify-pin --batch "$PINS_FILE" --output json > "$OUTPUT_FILE"

# Count invalid PINs
INVALID_COUNT=$(jq '[.[] | select(.is_valid == false)] | length' "$OUTPUT_FILE")

if [ "$INVALID_COUNT" -gt 0 ]; then
  echo "Warning: $INVALID_COUNT invalid PINs found"
  # Send notification or email
fi
```

### Example 4: Integration with jq

```bash
# Extract taxpayer names
kra-cli verify-pin --batch pins.csv --output json | jq '.[].taxpayer_name'

# Filter active companies
kra-cli verify-pin --batch pins.csv --output json | \
  jq '.[] | select(.status == "active" and .taxpayer_type == "Company")'

# Create summary report
kra-cli verify-pin --batch pins.csv --output json | \
  jq '{total: length, valid: [.[] | select(.is_valid == true)] | length}'
```

## Error Handling

KRA-CLI provides clear error messages:

```bash
# Invalid PIN format
$ kra-cli verify-pin INVALID
Error: failed to verify PIN: invalid PIN format

# Missing API key
$ kra-cli verify-pin P051234567A
Error: API key not set. Use --api-key flag, set KRA_API_KEY environment variable, or run: kra-cli config set api-key YOUR_KEY

# Network timeout
$ kra-cli verify-pin P051234567A --timeout 1
Error: failed to verify PIN: context deadline exceeded

# Authentication error
$ kra-cli verify-pin P051234567A --api-key invalid
Error: failed to verify PIN: authentication failed: invalid API key
```

## Development

### Building from Source

```bash
# Clone repository
git clone https://github.com/kra-connect/kra-cli.git
cd kra-cli

# Install dependencies
go mod download

# Build
go build -o kra-cli

# Run tests
go test ./...

# Build for all platforms
./scripts/build-all.sh
```

### Project Structure

```
kra-cli/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command and global flags
│   ├── verify_pin.go      # PIN verification command
│   ├── check_tcc.go       # TCC checking command
│   ├── validate_slip.go   # E-slip validation command
│   ├── file_nil_return.go # NIL return filing command
│   ├── get_taxpayer.go    # Taxpayer details command
│   └── config.go          # Configuration management
├── internal/              # Internal utilities
│   └── output.go         # Output formatting (table, JSON, CSV)
├── main.go               # Entry point
├── go.mod                # Go module definition
└── README.md             # This file
```

## Troubleshooting

### Command Not Found

If `kra-cli` command is not found after installation:

```bash
# Check if $GOPATH/bin is in your PATH
echo $PATH | grep $GOPATH/bin

# Add to PATH if missing (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:$(go env GOPATH)/bin
```

### API Key Issues

```bash
# Verify API key is set
kra-cli config get api-key

# Test with a simple command
kra-cli verify-pin P051234567A --verbose
```

### Certificate Issues

If you encounter SSL/TLS errors:

```bash
# Update certificates (Ubuntu/Debian)
sudo apt-get update && sudo apt-get install ca-certificates

# macOS
brew install ca-certificates
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](../../CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/kra-connect/kra-cli/issues)
- **Documentation**: [docs.kra-connect.dev](https://docs.kra-connect.dev)
- **Email**: support@kra-connect.dev

## Related Projects

- [Go SDK](../go-sdk) - Go library for KRA GavaConnect API
- [Python SDK](../python-sdk) - Python library
- [Node.js SDK](../node-sdk) - Node.js/TypeScript library
- [PHP SDK](../php-sdk) - PHP library
- [Flutter SDK](../flutter-sdk) - Dart/Flutter library

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.

---

**Made with ❤️ for Kenyan developers and businesses**
