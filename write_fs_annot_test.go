package neuro

// https://pkg.go.dev/testing

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestWriteRereadAnnot(t *testing.T) {

	annot, err := ReadFsAnnot("testdata/lh.aparc.annot")
	if err != nil {
		t.Errorf("ReadFsAnnot failed: %v", err)
	}

	// get a temp file.
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name()) // clean up
	annotFileName := file.Name()
	file.Close()

	err = WriteFsAnnot(annotFileName, annot)
	if err != nil {
		t.Errorf("WriteFsAnnot failed: %v", err)
	}

	annotReread, err := ReadFsAnnot(annotFileName)
	if err != nil {
		t.Errorf("ReadFsAnnot failed: %v", err)
	}

	if !reflect.DeepEqual(annot.VertexLabels, annotReread.VertexLabels) {
		t.Errorf("vertex labels differ after annot roundtrip.")
	}
	if !reflect.DeepEqual(annot.VertexIndices, annotReread.VertexIndices) {
		t.Errorf("vertex indices differ after annot roundtrip.")
	}
	if !reflect.DeepEqual(annot.Colortable, annotReread.Colortable) {
		t.Errorf("colortable differs after annot roundtrip: %+v vs %+v", annot.Colortable, annotReread.Colortable)
	}
	if !reflect.DeepEqual(annot.VertexColors(false), annotReread.VertexColors(false)) {
		t.Errorf("vertex colors differ after annot roundtrip.")
	}
}

func TestWriteRereadAnnotSmall(t *testing.T) {

	// Construct a small annotation by hand (e.g., two regions on 4 vertices).
	// Note that the Label values must be consistent with the RGBA values, as the
	// reader recomputes them on read (label = r + g*256 + b*65536 + a*16777216).
	annot := Annot{
		VertexIndices: []int32{0, 1, 2, 3},
		VertexLabels:  []int32{255, 255, 1639705, 255}, // red (255,0,0) and unknown (25,5,25)
		Colortable: Colortable{
			ID:    []int32{0, 1},
			Name:  []string{"unknown", "myregion"},
			R:     []int32{25, 255},
			G:     []int32{5, 0},
			B:     []int32{25, 0},
			A:     []int32{0, 0},
			Label: []int32{1639705, 255},
		},
	}

	// get a temp file.
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name()) // clean up
	annotFileName := file.Name()
	file.Close()

	err = WriteFsAnnot(annotFileName, annot)
	if err != nil {
		t.Errorf("WriteFsAnnot failed: %v", err)
	}

	annotReread, err := ReadFsAnnot(annotFileName)
	if err != nil {
		t.Errorf("ReadFsAnnot failed: %v", err)
	}

	numVertices, err := annotReread.NumVertices()
	if err != nil {
		t.Errorf("Annot.NumVertices failed: %v", err)
	}
	if numVertices != 4 {
		t.Errorf("got %d vertices after roundtrip, wanted 4", numVertices)
	}
	if !reflect.DeepEqual(annot.VertexLabels, annotReread.VertexLabels) {
		t.Errorf("vertex labels differ after annot roundtrip: %+v vs %+v", annot.VertexLabels, annotReread.VertexLabels)
	}
	if !reflect.DeepEqual(annot.Colortable, annotReread.Colortable) {
		t.Errorf("colortable differs after annot roundtrip: %+v vs %+v", annot.Colortable, annotReread.Colortable)
	}
}

func ExampleWriteFsAnnot() {

	// Read an annotation, write it to a temp file, and read it back.
	annot, _ := ReadFsAnnot("testdata/lh.aparc.annot")

	file, _ := os.CreateTemp("", "")
	defer os.Remove(file.Name())
	file.Close()

	_ = WriteFsAnnot(file.Name(), annot)
	annotReread, _ := ReadFsAnnot(file.Name())

	numVertices, _ := annotReread.NumVertices()
	fmt.Printf("Wrote and reread an annot with %d vertices and %d regions.\n",
		numVertices, annotReread.Colortable.NumEntries())
	// Output: Wrote and reread an annot with 149244 vertices and 36 regions.
}
