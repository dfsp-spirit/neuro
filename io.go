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

	n3, err := f.WriteString("writes\n")
    if err != nil {
		err = fmt.Errorf("strToTextFile: could not write to output text file '%s': '%s'.", filepath, err)
		return err
	}

    fmt.Printf("Wrote %d bytes\n", n3)

	f.Sync()



	return nil
}
