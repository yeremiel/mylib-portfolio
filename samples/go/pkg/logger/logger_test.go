package logger

import (
	"bytes"
	"strings"
	"testing"
)

// Both the package-level functions and a logger injected via Base() must report
// the real call site (this test file).
func TestCallerIsReportedAtCallSite(t *testing.T) {
	var buf bytes.Buffer
	if err := InitWithWriter(&buf, "info"); err != nil {
		t.Fatalf("InitWithWriter: %v", err)
	}
	t.Cleanup(func() { Log = nil })

	Info("via package function")
	Base().Info("via injected logger")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("log lines = %d, want 2: %q", len(lines), buf.String())
	}
	for _, line := range lines {
		if !strings.Contains(line, "logger/logger_test.go") {
			t.Errorf("caller should be the test file, got: %s", line)
		}
	}
}

func TestBase_NilLoggerReturnsNop(t *testing.T) {
	Log = nil
	Base().Info("must not panic")
}
