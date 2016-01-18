package main

import (
	"path/filepath"
	"testing"
)

func TestGetAllFilesNoDir(t *testing.T) {
	fp, err := getAllFiles("DIR_DOESNT_EXIST")
	if len(fp) != 0 {
		t.Fatalf("expected empty file path slice")
	}
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetAllFilesTestData(t *testing.T) {
	fps, err := getAllFiles("./testdata")
	if err != nil {
		t.Fatalf("expected error, got (%s)", err)
	}
	if len(fps) != 2 {
		t.Fatalf("expected 2 files, got %d", len(fps))
	}

	expected := map[string]int{
		"file1.txt": 0,
		"file2.txt": 0,
	}
	for _, fp := range fps {
		base := filepath.Base(fp.RelativePath)
		_, ok := expected[base]
		if !ok {
			t.Errorf("found %s, but not expected", base)
		}
		expected[base]++
		if expected[base] != 1 {
			t.Errorf("found %d instances of %s, expected 1", expected[base], base)
		}
	}
}
