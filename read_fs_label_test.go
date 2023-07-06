package neuro

// https://pkg.go.dev/testing

import (
	"fmt"
	"testing"
)

func TestReadFsLabel(t *testing.T){
	var curvFile string = "testdata/lh.cortex.label"

	// Read the label file
	label, _ := ReadFsLabel(curvFile)

	got := len(label.CoordX)
	want := 140891

	if got != want {
		t.Errorf("got data for %d cortex vertices in label file, wanted %d", got, want)
	}
}

func ExampleReadFsLabel() {
	var labelFile string = "testdata/lh.cortex.label"

	// Read the curv file
	label, _ := ReadFsLabel(labelFile)

	fmt.Printf("Read label containing %d vertices from label file '%s'.\n", len(label.ElementIndex), labelFile)
	// Output: Read label containing 140891 vertices from label file 'testdata/lh.cortex.label'.
}
