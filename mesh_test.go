package neurogo

import "testing"

func TestNumVertices(t *testing.T){

	var mymesh Mesh
	mymesh.Vertices = make([]float32, 5 * 3)

    got := NumVertices(mymesh)
    want := 5

    if got != want {
        t.Errorf("got %d, wanted %d", got, want)
    }
}

func TestNumVerticesEmpty(t *testing.T){

	var mymesh Mesh

    got := NumVertices(mymesh)
    want := 0

    if got != want {
        t.Errorf("got %d, wanted %d", got, want)
    }
}

func TestNumFaces(t *testing.T){

	var mymesh Mesh
	mymesh.Faces = make([]int32, 5 * 3)

    got := NumFaces(mymesh)
    want := 5

    if got != want {
        t.Errorf("got %d, wanted %d", got, want)
    }
}

func TestNumFacesEmpty(t *testing.T){

	var mymesh Mesh

    got := NumFaces(mymesh)
    want := 0

    if got != want {
        t.Errorf("got %d, wanted %d", got, want)
    }
}

func TestCubeFaces(t *testing.T) {
	var mycube Mesh = GenerateCube()

	got := NumFaces(mycube)
	want := 12

	if got != want {
        t.Errorf("got %d, wanted %d", got, want)
    }
}


func TestCubeVertices(t *testing.T) {
	var mycube Mesh = GenerateCube()

	got := NumVertices(mycube)
	want := 8

	if got != want {
        t.Errorf("got %d, wanted %d", got, want)
    }
}
