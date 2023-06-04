package neurogo

import (
	"fmt"
	"os"
)

// Write a string to a text file.
func strToTextFile(s string, filepath string) (error) {
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
