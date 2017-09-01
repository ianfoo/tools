package imgcat

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestIsSupported(t *testing.T) {
	defer func(old string) { os.Setenv("TERM_PROGRAM", old) }(os.Getenv("TERM_PROGRAM"))
	os.Setenv("TERM_PROGRAM", "foo")
	if _, err := New(nil); err == nil {
		t.Fatal("imgcat should not be supported now")
	}
}

func TestNew(t *testing.T) {
	// Change is supported to be always true and restore at the end.
	defer func(old func() bool) { isSupported = old }(isSupported)
	isSupported = func() bool { return true }

	tc := []struct {
		name    string
		in      string
		options []Option
		out     string
	}{
		// {"empty", "", nil, "\x1b]1337;File=:\a\n"},
		{"test", "test", nil, "\x1b]1337;File=:dGVzdA==\a\n"},
		{"test inline", "test", []Option{Inline(true)}, "\x1b]1337;File=inline=1:dGVzdA==\a\n"},
		{"test outline", "test", []Option{Inline(false)}, "\x1b]1337;File=inline=0:dGVzdA==\a\n"},
		{"test with name", "test", []Option{Name("test")}, "\x1b]1337;File=name=dGVzdA==:dGVzdA==\a\n"},
		{"test width 10 cells", "test", []Option{Width(Cells(10))}, "\x1b]1337;File=width=10:dGVzdA==\a\n"},
		{"test width 10px", "test", []Option{Width(Pixels(10))}, "\x1b]1337;File=width=10px:dGVzdA==\a\n"},
		{"test width 10%", "test", []Option{Width(Percent(10))}, "\x1b]1337;File=width=10%:dGVzdA==\a\n"},
		{"test width auto", "test", []Option{Width(Auto())}, "\x1b]1337;File=width=auto:dGVzdA==\a\n"},
		{"test with size", "test", []Option{Size(42)}, "\x1b]1337;File=size=42:dGVzdA==\a\n"},
		{"test height 10 cells", "test", []Option{Height(Cells(10))}, "\x1b]1337;File=height=10:dGVzdA==\a\n"},
		{"test height 10px", "test", []Option{Height(Pixels(10))}, "\x1b]1337;File=height=10px:dGVzdA==\a\n"},
		{"test height 10%", "test", []Option{Height(Percent(10))}, "\x1b]1337;File=height=10%:dGVzdA==\a\n"},
		{"test preserve aspect ration", "test", []Option{PreserveAspectRatio(true)}, "\x1b]1337;File=preserveAspectRatio=1:dGVzdA==\a\n"},
		{"test don't preserve aspect ration", "test", []Option{PreserveAspectRatio(false)}, "\x1b]1337;File=preserveAspectRatio=0:dGVzdA==\a\n"},
		{"all options together", "test", []Option{
			Inline(true), Name("test"), Width(Percent(10)), Height(Percent(10)), PreserveAspectRatio(false), Size(42),
		}, "\x1b]1337;File=inline=1;name=dGVzdA==;width=10%;height=10%;preserveAspectRatio=0;size=42:dGVzdA==\a\n"},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			w, err := New(&buf, tt.options...)
			if err != nil {
				t.Fatalf("could not create writer: %v", err)
			}
			_, err = fmt.Fprint(w, tt.in)
			if err != nil {
				t.Fatalf("could not write: %v", err)
			}
			if err := w.Close(); err != nil {
				t.Fatalf("could not close: %v", err)
			}
			if got := buf.String(); got != tt.out {
				t.Fatalf("expected output %q; got %q", tt.out, got)
			}
		})
	}
}

type badWriter struct{}

func (badWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("bad writer")
}

func TestWriter(t *testing.T) {
	// Change is supported to be always true and restore at the end.
	defer func(old func() bool) { isSupported = old }(isSupported)
	isSupported = func() bool { return true }

	w, err := New(badWriter{})
	if err != nil {
		t.Fatalf("could not create writer: %v", err)
	}
	_, err = fmt.Fprint(w, "test")
	if err == nil {
		t.Fatalf("expected error; got nothing")
	}
	if err.Error() != "bad writer" {
		t.Fatalf("expected error bad writer; got %v", err)
	}
}
