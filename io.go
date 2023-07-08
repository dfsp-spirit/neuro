package neuro

import (
	"fmt"
	"os"
	"bufio"
)

// Write a string to a text file.
//
// Parameters:
//  - s: the string to write
//  - filepath: the path to the file
//
// Returns:
//  - error: an error if one occurred
func strToTextFile(s string, filepath string) (error) {

	if Verbosity >= 2 {
		fmt.Printf("Writing string of length %d to file '%s'.\n", len(s), filepath)
	}

	f, err := os.Create(filepath)
    if err != nil {
		err = fmt.Errorf("strToTextFile: could not create new text file '%s': '%s'.", filepath, err)
		return err
	}

	defer f.Close()

	numBytesWritten, err := f.WriteString(s)
    if err != nil {
		err = fmt.Errorf("strToTextFile: could not write to output text file '%s': '%s'.", filepath, err)
		return err
	}

	if Verbosity >= 2 {
    	fmt.Printf("Wrote %d bytes to text file '%s'\n", numBytesWritten, filepath)
	}

	f.Sync()

	return nil
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
//
// Parameters:
//  - path: the path to the file
//
// Returns:
//  - lines: a slice of strings, each string is a line in the file
//  - error: an error if one occurred
func readLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

// writeLines writes the lines to the given file.
//
// Parameters:
//  - lines: a slice of strings, each string is a line in the file
//  - path: the path to the file
//
// Returns:
//  - error: an error if one occurred
func writeLines(lines []string, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    for _, line := range lines {
        fmt.Fprintln(w, line)
    }
    return w.Flush()
}
