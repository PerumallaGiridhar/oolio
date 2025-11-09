package index

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTxtFile(t *testing.T, dir, name string, lines []string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	content := ""
	for _, l := range lines {
		content += l + "\n"
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write txt file: %v", err)
	}
	return path
}

func TestEnsurePebble_CreatesAndLoadsFromTxt(t *testing.T) {
	dir := t.TempDir()

	txtPath := writeTxtFile(t, dir, "codes.txt", []string{
		"abc123",
		"  DEF456  ",
		"",
	})

	store, err := EnsurePebble(txtPath)
	if err != nil {
		t.Fatalf("EnsurePebble returned error: %v", err)
	}
	defer store.Close()

	// basic metadata
	if store.Txt != txtPath {
		t.Errorf("store.Txt = %q, want %q", store.Txt, txtPath)
	}
	expectedDbDir := txtPath + ".peb"
	if store.DbDir != expectedDbDir {
		t.Errorf("store.DbDir = %q, want %q", store.DbDir, expectedDbDir)
	}

	// manifest should now exist
	if !hasManifest(store.DbDir) {
		t.Errorf("expected manifest in %q after EnsurePebble", store.DbDir)
	}

	// inserted codes should be found, case-insensitive & trimmed
	found, err := store.Has("abc123")
	if err != nil {
		t.Fatalf("Has(abc123) error: %v", err)
	}
	if !found {
		t.Errorf("expected Has(abc123) = true")
	}

	found, err = store.Has("  def456  ")
	if err != nil {
		t.Fatalf("Has(def456) error: %v", err)
	}
	if !found {
		t.Errorf("expected Has(def456) = true (case-insensitive, trimmed)")
	}

	// missing code
	found, err = store.Has("ZZZ999")
	if err != nil {
		t.Fatalf("Has(ZZZ999) error: %v", err)
	}
	if found {
		t.Errorf("expected Has(ZZZ999) = false for missing code")
	}

	// empty / whitespace should be false, no error
	found, err = store.Has("   ")
	if err != nil {
		t.Fatalf("Has(whitespace) error: %v", err)
	}
	if found {
		t.Errorf("expected Has(whitespace) = false")
	}
}

func TestEnsurePebble_ReusesExistingDB(t *testing.T) {
	dir := t.TempDir()

	txtPath := writeTxtFile(t, dir, "codes.txt", []string{"FIRST123"})

	store1, err := EnsurePebble(txtPath)
	if err != nil {
		t.Fatalf("first EnsurePebble error: %v", err)
	}
	if !hasManifest(store1.DbDir) {
		t.Fatalf("expected manifest after first EnsurePebble")
	}
	if ok, _ := store1.Has("FIRST123"); !ok {
		t.Fatalf("expected FIRST123 to exist after first load")
	}
	store1.Close()

	txtPath = writeTxtFile(t, dir, "codes.txt", []string{"FIRST123", "NEWCODE99"})

	store2, err := EnsurePebble(txtPath)
	if err != nil {
		t.Fatalf("second EnsurePebble error: %v", err)
	}
	defer store2.Close()

	if ok, _ := store2.Has("FIRST123"); !ok {
		t.Errorf("expected FIRST123 to still exist in reused DB")
	}

	if ok, _ := store2.Has("NEWCODE99"); ok {
		t.Errorf("expected NEWCODE99 to NOT exist; DB should have been reused, not reloaded")
	}
}

func TestHasManifest(t *testing.T) {
	dir := t.TempDir()

	if hasManifest(dir) {
		t.Fatalf("expected hasManifest(%q) = false for empty dir", dir)
	}

	manifestPath := filepath.Join(dir, "MANIFEST-000001")
	if err := os.WriteFile(manifestPath, []byte("dummy"), 0o644); err != nil {
		t.Fatalf("failed to write manifest file: %v", err)
	}

	if !hasManifest(dir) {
		t.Fatalf("expected hasManifest(%q) = true after manifest file created", dir)
	}
}
