package internal

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	original := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = original

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("failed to copy stdout: %v", err)
	}
	_ = r.Close()

	return buf.String()
}

func TestOutputFormatterJSON(t *testing.T) {
	formatter := NewOutputFormatter("json")
	payload := map[string]string{"status": "ok"}

	out := captureStdout(t, func() {
		if err := formatter.Print(payload); err != nil {
			t.Fatalf("Print returned error: %v", err)
		}
	})

	if want := "\"status\": \"ok\""; !bytes.Contains([]byte(out), []byte(want)) {
		t.Fatalf("expected JSON output to contain %q, got %q", want, out)
	}
}

func TestOutputFormatterUnsupportedFormat(t *testing.T) {
	formatter := NewOutputFormatter("unsupported")
	err := formatter.Print(map[string]string{})
	if err == nil {
		t.Fatalf("expected error for unsupported format")
	}
}
