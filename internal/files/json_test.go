package files

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testPayload struct {
	Name string `json:"name"`
}

func TestReadJSONValid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "valid.json")
	if err := os.WriteFile(path, []byte(`{"name":"railyard"}`), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	value, err := ReadJSON[testPayload](path, "test payload", JSONReadOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value.Name != "railyard" {
		t.Fatalf("unexpected payload: %#v", value)
	}
}

func TestReadJSONMissingStrictErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")

	_, err := ReadJSON[testPayload](path, "mods index", JSONReadOptions{})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if !strings.Contains(err.Error(), "failed to read mods index") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestReadJSONMissingAllowedReturnsZero(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")

	value, err := ReadJSON[testPayload](path, "app config", JSONReadOptions{AllowMissing: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != (testPayload{}) {
		t.Fatalf("expected zero value, got %#v", value)
	}
}

func TestReadJSONEmptyStrictErrors(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(path, []byte("   \n\t  "), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := ReadJSON[testPayload](path, "maps index", JSONReadOptions{})
	if err == nil {
		t.Fatal("expected error for empty file")
	}
	if !strings.Contains(err.Error(), "failed to parse maps index") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestReadJSONEmptyAllowedReturnsZero(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.json")
	if err := os.WriteFile(path, []byte("\n"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	value, err := ReadJSON[testPayload](path, "app config", JSONReadOptions{AllowEmpty: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != (testPayload{}) {
		t.Fatalf("expected zero value, got %#v", value)
	}
}

func TestReadJSONInvalidJSONIncludesLabel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.json")
	if err := os.WriteFile(path, []byte(`{"name":`), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := ReadJSON[testPayload](path, "app config", JSONReadOptions{})
	if err == nil {
		t.Fatal("expected parse error")
	}
	if !strings.Contains(err.Error(), "failed to parse app config") {
		t.Fatalf("unexpected error message: %v", err)
	}
}
