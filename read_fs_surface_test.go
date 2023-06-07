package neuro

// https://pkg.go.dev/testing

import (
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

