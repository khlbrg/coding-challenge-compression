package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}

func run() error {
	filename := flag.String("f", "compressed.txt", "Path to the file that should be used")
	decom := flag.Bool("d", true, "If decompress should be used")
	verbose := flag.Bool("v", false, "Echo the result to the terminal")

	flag.Parse()

	fileContent, err := loadFile(*filename)
	if err != nil {
		log.Fatalf("failed to load file: %v", err)
	}

	if *decom {
		decompress(fileContent)
	} else {
		compress(string(fileContent), *verbose)
	}
	return nil
}

func compress(content string, verbose bool) error {
	compressedData, err := Compress(string(content))
	if err != nil {
		return fmt.Errorf("failed to compress: %w", err)
	}

	if err := writeFile("compressed.txt", compressedData); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	if verbose {
		fmt.Printf("\nData compressed from %d bytes to %d bytes\n", len(content), len(compressedData))
		fmt.Printf("saved %.2f%% of memory\n", float64(len(content)-len(compressedData))/float64(len(content))*100)
	}

	return nil
}

func decompress(content []byte) error {
	decompressedData, err := Decompress(content)
	if err != nil {
		return fmt.Errorf("failed to decompress: %w", err)
	}
	if err := writeFile("decompressed.txt", []byte(decompressedData)); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func writeFile(filename string, content []byte) error {
	err := os.WriteFile(filename, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func loadFile(filename string) ([]byte, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}
	return fileContent, nil
}
