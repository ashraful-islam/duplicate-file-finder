package main

import "testing"

var testDir string = "test_data"
var processResult ProcessResult

func TestProcess(t *testing.T) {

	expected := ProcessResult{
		NumScannedFiles: 3,
		SizeScannedFiles: 112,
		NumDuplicateFiles:2,
		SizeDuplicateFiles: 72,
	}

	found := Process(testDir)

	isCountMatch := expected.NumScannedFiles == found.NumScannedFiles
	isCountMatch = isCountMatch && expected.NumDuplicateFiles == expected.NumDuplicateFiles

	isSizeMatch := expected.SizeScannedFiles == found.SizeScannedFiles
	isSizeMatch = isSizeMatch && expected.SizeDuplicateFiles == found.SizeDuplicateFiles

	if !isCountMatch {
		t.Errorf(
			"Expected Count to be %v, %v Found %v, %v",
			expected.NumScannedFiles,
			expected.NumDuplicateFiles,
			found.NumScannedFiles,
			found.NumDuplicateFiles,
			)
	}

	if !isSizeMatch {
		t.Errorf(
			"Expected Size to be %v, %v Found %v, %v",
			expected.SizeScannedFiles,
			expected.SizeDuplicateFiles,
			found.SizeScannedFiles,
			found.SizeDuplicateFiles,
		)
	}

}

func BenchmarkProcess(b *testing.B) {

	var result ProcessResult

	for n := 0; n < b.N; n++ {
		// storing result here will eliminate compiler removing
		// this function call altogether
		result = Process(testDir)
	}

	// storing the result in package level variable will
	// mitigate elimination of the benchmark function itself by compiler
	// src: https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go#compiler-optimisation
	processResult = result
}