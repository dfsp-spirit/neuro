package neurogo

import "testing"

func TestMean(t *testing.T){

	data := []float32{1.0, 2.0, 3.0, 4.0, 5.0}

    got, _ := mean(data)
    want := float32(3.0)

    if got != want {
        t.Errorf("got %f, wanted %f", got, want)
    }
}


func TestMax(t *testing.T){

	data := []float32{1.0, 2.0, 3.0, 4.0, 5.0}

    got, _ := max(data)
    want := float32(5.0)

    if got != want {
        t.Errorf("got %f, wanted %f", got, want)
    }
}


func TestMin(t *testing.T){

	data := []float32{1.0, 2.0, 3.0, 4.0, 5.0}

    got, _ := min(data)
    want := float32(1.0)

    if got != want {
        t.Errorf("got %f, wanted %f", got, want)
    }
}