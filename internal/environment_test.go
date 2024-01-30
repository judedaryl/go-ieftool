package internal

import (
	"path/filepath"
	"testing"
)

func Test_NewEnvironmentsWithoutFilter(t *testing.T) {
	p, _ := filepath.Abs("../test/fixtures/config.yaml")
	e, err := NewEnvironmentsFromConfig(p, "")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if e.Len() == 1 {
		t.Fatalf("expect more than one environment, got %d", e.Len())
	}
}

func Test_NewEnvironmentsWithFilter(t *testing.T) {
	p, _ := filepath.Abs("../test/fixtures/config.yaml")
	e, err := NewEnvironmentsFromConfig(p, "test")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if e.Len() != 1 {
		t.Fatalf("expect more than one environment, got %d", e.Len())
	}
}
