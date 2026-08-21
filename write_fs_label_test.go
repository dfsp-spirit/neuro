package neuro

// https://pkg.go.dev/testing

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteRereadLabel(t *testing.T) {

	label, err := ReadFsLabel("testdata/lh.cortex.label")
	if err != nil {
		t.Errorf("ReadFsLabel failed: %v", err)
	}

	// get a temp file.
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name()) // clean up
	labelFileName := file.Name()
	file.Close()

	err = WriteFsLabel(labelFileName, label)
	if err != nil {
		t.Errorf("WriteFsLabel failed: %v", err)
	}

	labelReread, err := ReadFsLabel(labelFileName)
	if err != nil {
		t.Errorf("ReadFsLabel failed: %v", err)
	}

	if diff := cmp.Diff(label, labelReread); diff != "" {
		t.Error(diff)
	}
}

func TestWriteRereadLabelSmall(t *testing.T) {

	// Construct a small label by hand (with typical non-trivial float values).
	label := FsLabel{
		ElementIndex: []int32{0, 1, 2, 3, 100, 5000},
		CoordX:       []float32{-3.5, 0.0, 12.25, 0.123456, -100.0, 3.14159},
		CoordY:       []float32{1.0, -2.0, 3.0, 4.0, 5.0, 6.0},
		CoordZ:       []float32{0.0, 0.1, 0.2, 0.3, 0.4, 0.5},
		Value:        []float32{0.0, 1.5, -2.5, 3.25, 4.125, 5.625},
	}

	// get a temp file.
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name()) // clean up
	labelFileName := file.Name()
	file.Close()

	err = WriteFsLabel(labelFileName, label)
	if err != nil {
		t.Errorf("WriteFsLabel failed: %v", err)
	}

	labelReread, err := ReadFsLabel(labelFileName)
	if err != nil {
		t.Errorf("ReadFsLabel failed: %v", err)
	}

	if diff := cmp.Diff(label, labelReread); diff != "" {
		t.Error(diff)
	}
}

func ExampleWriteFsLabel() {

	// Read a label, write it to a temp file, and read it back.
	label, _ := ReadFsLabel("testdata/lh.cortex.label")

	file, _ := os.CreateTemp("", "")
	defer os.Remove(file.Name())
	file.Close()

	_ = WriteFsLabel(file.Name(), label)
	labelReread, _ := ReadFsLabel(file.Name())

	fmt.Printf("Wrote and reread a label containing %d vertices.\n", len(labelReread.ElementIndex))
	// Output: Wrote and reread a label containing 140891 vertices.
}
