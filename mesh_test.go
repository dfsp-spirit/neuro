package neurogo

import "testing"

func TestNumVertices(t *testing.T){

	var mymesh Mesh
	mymesh.Vertices = make([]float32, 5 * 3)

    got := NumVertices(mymesh)
    want := 5

    if got != want {
        t.Errorf("got %q, wanted %q", got, want)
    }
}

func TestNumVerticesEmpty(t *testing.T){

	var mymesh Mesh

    got := NumVertices(mymesh)
    want := 0

    if got != want {
        t.Errorf("got %q, wanted %q", got, want)
    }
}

func TestNumFaces(t *testing.T){

	var mymesh Mesh
	mymesh.Faces = make([]int32, 5 * 3)

    got := NumFaces(mymesh)
    want := 5

    if got != want {
        t.Errorf("got %q, wanted %q", got, want)
    }
}

func TestNumFacesEmpty(t *testing.T){

	var mymesh Mesh

    got := NumFaces(mymesh)
    want := 0

    if got != want {
        t.Errorf("got %q, wanted %q", got, want)
    }
}
