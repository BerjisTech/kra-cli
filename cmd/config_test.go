package cmd

import "testing"

func TestConvertKeyToViperFormat(t *testing.T) {
	cases := map[string]string{
		"api-key":   "api_key",
		"base-url":  "base_url",
		"other-key": "other-key",
	}

	for input, expected := range cases {
		if got := convertKeyToViperFormat(input); got != expected {
			t.Fatalf("convertKeyToViperFormat(%q) = %q, expected %q", input, got, expected)
		}
	}
}
