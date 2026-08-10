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

func TestBuildKernelSampleNestedVault(t *testing.T) {
	tmpDir := t.TempDir()

	// Build sample vault layout
	// tmpDir/
	//   root_notes/index-note.md
	//   projects/quantum/quantum-intro.md
	//   archive/2025/old-notes.md
	//   .obsidian/config-note.md
	dirs := []string{
		filepath.Join(tmpDir, "root_notes"),
		filepath.Join(tmpDir, "projects", "quantum"),
		filepath.Join(tmpDir, "archive", "2025"),
		filepath.Join(tmpDir, ".obsidian"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatal(err)
		}
	}

	os.WriteFile(filepath.Join(tmpDir, "root_notes", "index-note.md"), []byte("[[quantum-intro]] and [[old-notes]]"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "projects", "quantum", "quantum-intro.md"), []byte("[[index-note]] and [[old-notes]]"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "archive", "2025", "old-notes.md"), []byte("[[quantum-intro]]"), 0644)
	os.WriteFile(filepath.Join(tmpDir, ".obsidian", "config-note.md"), []byte("[[index-note]]"), 0644)

	exclude := makeExcludeMap(nil)

	// Flat scan on root -> no pages error
	_, _, _, _, err := buildKernel(tmpDir, false, exclude)
	if err == nil {
		t.Fatalf("expected error on flat scan with no root .md files, got nil")
	}

	// Recursive scan -> 3 pages found, .obsidian ignored
	kern, _, pages, _, err := buildKernel(tmpDir, true, exclude)
	if err != nil {
		t.Fatalf("recursive buildKernel failed: %v", err)
	}

	if len(pages) != 3 {
		t.Fatalf("expected 3 pages, got %d (%v)", len(pages), pages)
	}

	expectedPages := []string{"index-note", "old-notes", "quantum-intro"}
	for i, p := range pages {
		if p != expectedPages[i] {
			t.Errorf("expected page %d to be %s, got %s", i, expectedPages[i], p)
		}
	}

	pi, err := kern.Stationary(1e-12, 5000)
	if err != nil {
		t.Fatalf("stationary distribution failed: %v", err)
	}

	// quantum-intro has inbound links from both index-note and old-notes, so it should have highest pi
	if pi[2] <= pi[0] || pi[2] <= pi[1] {
		t.Errorf("expected quantum-intro to be most central, got pi=%v", pi)
	}
}
