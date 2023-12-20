// main_test.go
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatTextWithTestData(t *testing.T) {
	inputDir := filepath.Join("testdata", "input_tests")
	outputDir := filepath.Join("testdata", "output_tests")

	inputFiles, err := os.ReadDir(inputDir)
	if err != nil {
		t.Fatalf("Error reading input test directory: %v", err)
	}

	for _, inputFile := range inputFiles {
		if inputFile.IsDir() {
			continue
		}

		inputFilePath := filepath.Join(inputDir, inputFile.Name())
		expectedOutputFilePath := filepath.Join(outputDir, strings.Replace(inputFile.Name(), "input", "output", 1))

		testName := strings.TrimSuffix(inputFile.Name(), "_input.txt")

		t.Run(testName, func(t *testing.T) {
			// Set up test files

			input, err := os.ReadFile(inputFilePath)
			if err != nil {
				t.Fatalf("Error reading test input file: %v", err)
			}

			expectedOutput, err := os.ReadFile(expectedOutputFilePath)
			if err != nil {
				t.Fatalf("Error reading expected output file: %v", err)
			}

			// Run your text formatting function
			result := FormatText(string(input))

			// Perform assertions on the result
			if result != string(expectedOutput) {
				t.Errorf("Test failed. Expected:\n`%s`\nGot:\n`%s`", expectedOutput, result)
			}
		})
	}
}
