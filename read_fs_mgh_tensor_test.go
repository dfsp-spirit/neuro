package neuro

// https://pkg.go.dev/testing

import (
	"fmt"

	"gorgonia.org/tensor"
)

func ExampleReadFsMgh_tensor() {
	// Illustrates how to read the 4D MRI image returned by ReadFsMgh into a tensor data structure
	// from the "gorgonia.org/tensor" package for convenient access to voxel values.
	// This example requires 'import "gorgonia.org/tensor"'.
	var mgzFile string = "testdata/brain.mgz"

	mgh, _ := ReadFsMgh(mgzFile, "yes")
	var h MghHeader = mgh.Header

	data := tensor.New(tensor.WithShape(int(h.Dim1Length), int(h.Dim2Length), int(h.Dim3Length), int(h.Dim4Length)), tensor.WithBacking(mgh.Data.DataMriUchar))
	val1, _ := data.At(99, 99, 99, 0)    // 77
	val2, _ := data.At(109, 109, 109, 0) // 71
	val3, _ := data.At(0, 0, 0, 0)       // 0

	// values known from external tests with FreeSurfer software, try on command line: mri_info --voxel 99 99 99 testdata/brain.mgz

	fmt.Printf("Voxel values=%d, %d, %d", val1, val2, val3)
	// Output: Voxel values=77, 71, 0
}
