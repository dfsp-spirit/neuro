package neuro

// https://pkg.go.dev/testing

import (
	"fmt"
	"testing"
)

func TestReadFsCurv(t *testing.T){
	var curvFile string = "testdata/lh.thickness"

	// Read the curv file
	pvdata, _ := ReadFsCurv(curvFile)

	got := len(pvdata)
	want := 149244

	if got != want {
		t.Errorf("got data for %d vertices in curv file, wanted %d", got, want)
	}

	// Compute some stats on the data and check them
	gotMax, _ := max(pvdata)
	if gotMax < 2.0 || gotMax > 5.0 {
		t.Errorf("got max value %f in curv file, wanted between 2.0 and 5.0", gotMax)
	}

	gotMin, _ := min(pvdata)
	if gotMin < 0.0 || gotMin > 2.0 {
		t.Errorf("got min value %f in curv file, wanted between 0.0 and 2.0", gotMin)
	}

	gotMean, _ := mean(pvdata)
	if gotMean < 1.0 || gotMean > 3.0 {
		t.Errorf("got mean value %f in curv file, wanted between 1.0 and 3.0", gotMean)
	}

}

func ExampleReadFsCurv() {
	var curvFile string = "testdata/lh.thickness"

	// Read the curv file
	pvdata, _ := ReadFsCurv(curvFile)

	fmt.Printf("Read %d values from curv file '%s'.\n", len(pvdata), curvFile)
	// Output: Read 149244 values from curv file 'testdata/lh.thickness'.
}
