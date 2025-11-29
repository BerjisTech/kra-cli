// Package internal provides internal utilities for the KRA-CLI tool
package internal

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// OutputFormatter handles output formatting for different formats (table, JSON, CSV)
type OutputFormatter struct {
	Format string // "table", "json", or "csv"
}

// NewOutputFormatter creates a new output formatter
func NewOutputFormatter(format string) *OutputFormatter {
	return &OutputFormatter{Format: format}
}

// Print outputs data in the specified format
func (f *OutputFormatter) Print(data interface{}) error {
	switch strings.ToLower(f.Format) {
	case "json":
		return f.printJSON(data)
	case "csv":
		return f.printCSV(data)
	case "table":
		return f.printTable(data)
	default:
		return fmt.Errorf("unsupported output format: %s (supported: table, json, csv)", f.Format)
	}
}

// printJSON outputs data as JSON
func (f *OutputFormatter) printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// printCSV outputs data as CSV
func (f *OutputFormatter) printCSV(data interface{}) error {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	// Handle slice of structs or maps
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			return nil
		}

		// Get first element to determine structure
		first := v.Index(0)
		if first.Kind() == reflect.Struct || (first.Kind() == reflect.Ptr && first.Elem().Kind() == reflect.Struct) {
			return f.printStructSliceCSV(data)
		} else if first.Kind() == reflect.Map {
			return f.printMapSliceCSV(data)
		}
	} else if v.Kind() == reflect.Struct || (v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct) {
		return f.printStructSliceCSV([]interface{}{data})
	}

	return fmt.Errorf("unsupported data type for CSV output")
}

// printStructSliceCSV prints a slice of structs as CSV
func (f *OutputFormatter) printStructSliceCSV(data interface{}) error {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	v := reflect.ValueOf(data)
	if v.Len() == 0 {
		return nil
	}

	first := v.Index(0)
	if first.Kind() == reflect.Ptr {
		first = first.Elem()
	}

	t := first.Type()
	headers := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Use JSON tag if available, otherwise use field name
		tag := field.Tag.Get("json")
		if tag != "" && tag != "-" {
			headers[i] = strings.Split(tag, ",")[0]
		} else {
			headers[i] = field.Name
		}
	}

	if err := writer.Write(headers); err != nil {
		return err
	}

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		row := make([]string, t.NumField())
		for j := 0; j < t.NumField(); j++ {
			field := item.Field(j)
			row[j] = fmt.Sprintf("%v", field.Interface())
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// printMapSliceCSV prints a slice of maps as CSV
func (f *OutputFormatter) printMapSliceCSV(data interface{}) error {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	v := reflect.ValueOf(data)
	if v.Len() == 0 {
		return nil
	}

	// Get headers from first map
	first := v.Index(0).Interface().(map[string]interface{})
	headers := make([]string, 0, len(first))
	for k := range first {
		headers = append(headers, k)
	}

	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write rows
	for i := 0; i < v.Len(); i++ {
		m := v.Index(i).Interface().(map[string]interface{})
		row := make([]string, len(headers))
		for j, h := range headers {
			row[j] = fmt.Sprintf("%v", m[h])
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// printTable outputs data as a formatted table
func (f *OutputFormatter) printTable(data interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	v := reflect.ValueOf(data)

	// Handle slice
	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			fmt.Println("No data to display")
			return nil
		}

		first := v.Index(0)
		if first.Kind() == reflect.Struct || (first.Kind() == reflect.Ptr && first.Elem().Kind() == reflect.Struct) {
			return f.printStructSliceTable(table, data)
		} else if first.Kind() == reflect.Map {
			return f.printMapSliceTable(table, data)
		}
	} else if v.Kind() == reflect.Struct || (v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct) {
		return f.printSingleStructTable(table, data)
	} else if v.Kind() == reflect.Map {
		return f.printSingleMapTable(table, data)
	}

	return fmt.Errorf("unsupported data type for table output")
}

// printSingleStructTable prints a single struct as a vertical table
func (f *OutputFormatter) printSingleStructTable(table *tablewriter.Table, data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Use JSON tag if available
		name := field.Name
		tag := field.Tag.Get("json")
		if tag != "" && tag != "-" {
			name = strings.Split(tag, ",")[0]
		}

		// Skip nil pointers
		if value.Kind() == reflect.Ptr && value.IsNil() {
			continue
		}

		table.Append([]string{name, fmt.Sprintf("%v", value.Interface())})
	}

	table.Render()
	return nil
}

// printSingleMapTable prints a single map as a vertical table
func (f *OutputFormatter) printSingleMapTable(table *tablewriter.Table, data interface{}) error {
	m := data.(map[string]interface{})
	for k, v := range m {
		table.Append([]string{k, fmt.Sprintf("%v", v)})
	}

	table.Render()
	return nil
}

// printStructSliceTable prints a slice of structs as a horizontal table
func (f *OutputFormatter) printStructSliceTable(table *tablewriter.Table, data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Len() == 0 {
		fmt.Println("No data to display")
		return nil
	}

	first := v.Index(0)
	if first.Kind() == reflect.Ptr {
		first = first.Elem()
	}

	t := first.Type()
	headers := make([]string, 0)
	visibleFields := make([]int, 0)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag != "" && tag != "-" {
			headers = append(headers, strings.Split(tag, ",")[0])
			visibleFields = append(visibleFields, i)
		} else {
			headers = append(headers, field.Name)
			visibleFields = append(visibleFields, i)
		}
	}

	table.SetHeader(headers)

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		row := make([]string, len(visibleFields))
		for j, fieldIdx := range visibleFields {
			field := item.Field(fieldIdx)
			row[j] = fmt.Sprintf("%v", field.Interface())
		}

		table.Append(row)
	}

	table.Render()
	return nil
}

// printMapSliceTable prints a slice of maps as a horizontal table
func (f *OutputFormatter) printMapSliceTable(table *tablewriter.Table, data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Len() == 0 {
		fmt.Println("No data to display")
		return nil
	}

	// Get headers from first map
	first := v.Index(0).Interface().(map[string]interface{})
	headers := make([]string, 0, len(first))
	for k := range first {
		headers = append(headers, k)
	}

	table.SetHeader(headers)

	// Write rows
	for i := 0; i < v.Len(); i++ {
		m := v.Index(i).Interface().(map[string]interface{})
		row := make([]string, len(headers))
		for j, h := range headers {
			row[j] = fmt.Sprintf("%v", m[h])
		}

		table.Append(row)
	}

	table.Render()
	return nil
}

// PrintError prints an error message to stderr
func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

// PrintSuccess prints a success message to stdout
func PrintSuccess(message string) {
	fmt.Println("✓", message)
}

// PrintWarning prints a warning message to stdout
func PrintWarning(message string) {
	fmt.Println("⚠", message)
}
