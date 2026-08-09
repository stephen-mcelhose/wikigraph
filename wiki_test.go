package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPagesFlatAndRecursive(t *testing.T) {
	tmpDir := t.TempDir()

	// Create structure:
	// tmpDir/
	//   page1.md
	//   sub/
	//     page2.md
	//   .hidden/
	//     page3.md
	if err := os.WriteFile(filepath.Join(tmpDir, "page1.md"), []byte("[[page2]]"), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(tmpDir, "sub")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "page2.md"), []byte("[[page1]]"), 0644); err != nil {
		t.Fatal(err)
	}
	hiddenDir := filepath.Join(tmpDir, ".hidden")
	if err := os.Mkdir(hiddenDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(hiddenDir, "page3.md"), []byte("[[page1]]"), 0644); err != nil {
		t.Fatal(err)
	}

	exclude := makeExcludeMap(nil)

	// Flat scan: should only find page1
	pagesFlat, _, _, err := loadPages(tmpDir, false, exclude)
	if err != nil {
		t.Fatalf("flat loadPages failed: %v", err)
	}
	if len(pagesFlat) != 1 || pagesFlat[0] != "page1" {
		t.Errorf("expected [page1], got %v", pagesFlat)
	}

	// Recursive scan: should find page1 and page2, skipping .hidden
	pagesRec, _, _, err := loadPages(tmpDir, true, exclude)
	if err != nil {
		t.Fatalf("recursive loadPages failed: %v", err)
	}
	if len(pagesRec) != 2 || pagesRec[0] != "page1" || pagesRec[1] != "page2" {
		t.Errorf("expected [page1, page2], got %v", pagesRec)
	}
}

func TestLoadPagesDuplicateSlugError(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "note.md"), []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}
	subDir := filepath.Join(tmpDir, "folder")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "note.md"), []byte("world"), 0644); err != nil {
		t.Fatal(err)
	}

	exclude := makeExcludeMap(nil)
	_, _, _, err := loadPages(tmpDir, true, exclude)
	if err == nil {
		t.Fatalf("expected error on duplicate slug, got nil")
	}
}
