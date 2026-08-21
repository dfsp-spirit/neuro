package neuro

// https://pkg.go.dev/testing

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteRereadSurfaceCube(t *testing.T) {

	mesh := GenerateCube()

	// get a temp file.
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name()) // clean up
	surfFileName := file.Name()
	file.Close()

	err = WriteFsSurface(surfFileName, mesh)
	if err != nil {
		t.Errorf("WriteFsSurface failed: %v", err)
	}

	meshReread, err := ReadFsSurface(surfFileName)
	if err != nil {
		t.Errorf("ReadFsSurface failed: %v", err)
	}

	if diff := cmp.Diff(mesh, meshReread); diff != "" {
		t.Error(diff)
	}
}

func TestWriteRereadSurfaceBrain(t *testing.T) {

	// Read the real surface, write it, read it back, and compare.
	mesh, err := ReadFsSurface("testdata/lh.white")
	if err != nil {
		t.Errorf("ReadFsSurface failed: %v", err)
	}

	// get a temp file.
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name()) // clean up
	surfFileName := file.Name()
	file.Close()

	err = WriteFsSurface(surfFileName, mesh)
	if err != nil {
		t.Errorf("WriteFsSurface failed: %v", err)
	}

	meshReread, err := ReadFsSurface(surfFileName)
	if err != nil {
		t.Errorf("ReadFsSurface failed: %v", err)
	}

	if len(mesh.Vertices) != len(meshReread.Vertices) {
		t.Errorf("got %d vertex coordinates after roundtrip, wanted %d", len(meshReread.Vertices), len(mesh.Vertices))
	}
	if len(mesh.Faces) != len(meshReread.Faces) {
		t.Errorf("got %d face indices after roundtrip, wanted %d", len(meshReread.Faces), len(mesh.Faces))
	}
	if !reflect.DeepEqual(mesh.Vertices, meshReread.Vertices) {
		t.Errorf("vertex coordinates differ after surface roundtrip.")
	}
	if !reflect.DeepEqual(mesh.Faces, meshReread.Faces) {
		t.Errorf("face indices differ after surface roundtrip.")
	}
}

func ExampleWriteFsSurface() {

	// Create a small mesh and write it to a temp file, then read it back.
	mesh := GenerateCube()

	file, _ := os.CreateTemp("", "")
	defer os.Remove(file.Name())
	file.Close()

	_ = WriteFsSurface(file.Name(), mesh)
	meshReread, _ := ReadFsSurface(file.Name())

	fmt.Printf("Wrote and reread a mesh with %d vertices and %d faces.\n",
		len(meshReread.Vertices)/3, len(meshReread.Faces)/3)
	// Output: Wrote and reread a mesh with 8 vertices and 12 faces.
}
