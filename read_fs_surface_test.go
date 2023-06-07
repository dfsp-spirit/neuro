package neuro

// https://pkg.go.dev/testing

import (
	"fmt"
	"testing"
)

func TestReadFsSurface(t *testing.T){
	var surfFile string = "testdata/lh.white"

	// Read the surface file
	surf, _ := ReadFsSurface(surfFile)

	got := NumVertices(surf)
	want := 149244

	if got != want {
		t.Errorf("got %d vertices in surface file, wanted %d", got, want)
	}

	got = NumFaces(surf)
	want = 298484

	if got != want {
		t.Errorf("got %d faces in surface file, wanted %d", got, want)
	}
}

func ExampleReadFsSurface() {
	var surfaceFile string = "testdata/lh.white"

	// Read the curv file
	mesh, _ := ReadFsSurface(surfaceFile)

	fmt.Printf("Read mesh with %d vertices and %d faces from surface file '%s'.\n", len(mesh.Vertices)/3, len(mesh.Faces)/3, surfaceFile)
	// Output: Read mesh with 149244 vertices and 298484 faces from surface file 'testdata/lh.white'.
}
