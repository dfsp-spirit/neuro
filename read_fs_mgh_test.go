package neuro

// https://pkg.go.dev/testing

import (
	"testing"
)

func TestReadFsMghHeader(t *testing.T){
	var mghFile string = "testdata/brain.mgh"

	hdr, _ := ReadFsMghHeader(mghFile, "auto")

	got := hdr.MghDataType
	want, _ := getMghDataTypeCode("MRI_UCHAR")

	if got != want {
		t.Errorf("got MGH data type %d, wanted %d", got, want)
	}
}

func TestReadFsMghHeaderRAS(t *testing.T){
	var mghFile string = "testdata/brain.mgh"

	hdr, _ := ReadFsMghHeader(mghFile, "no")

	got := hdr.RasGoodFlag
	var want int16 = 1

	if got != want {
		t.Errorf("got MGH RasGoodFlag=%d, wanted %d", got, want)
	}
}

func TestReadFsMghFull(t *testing.T){
	var mghFile string = "testdata/brain.mgh"

	mgh, _ := ReadFsMgh(mghFile, "auto")

	got := mgh.header.RasGoodFlag
	var want int16 = 1

	if got != want {
		t.Errorf("got MGH RasGoodFlag=%d, wanted %d", got, want)
	}
}

func TestReadFsMghFullAt0000(t *testing.T){
	var mghFile string = "testdata/brain.mgh"

	mgh, _ := ReadFsMgh(mghFile, "no")

	got := mgh.data.DataMriUchar[0]
	var want uint8 = 0

	if got != want {
		t.Errorf("got data value=%d, wanted %d", got, want)
	}
}

func TestReadFsMghFullSum(t *testing.T){
	var mghFile string = "testdata/brain.mgh"

	mgh, _ := ReadFsMgh(mghFile, "no")

	var sum int = 0
	for _, voxel_val := range mgh.data.DataMriUchar {
		sum += int(voxel_val)
	}
	var got = sum
	var want int = 121035479  // known from external tests with standard software.

	if got != want {
		t.Errorf("got MGH data sum=%d, wanted %d", got, want)
	}
}

func TestReadFsMgzFullSum(t *testing.T){
	var mgzFile string = "testdata/brain.mgz"

	mgh, _ := ReadFsMgh(mgzFile, "yes")

	var sum int = 0
	for _, voxel_val := range mgh.data.DataMriUchar {
		sum += int(voxel_val)
	}
	var got = sum
	var want int = 121035479  // known from external tests with standard software.

	if got != want {
		t.Errorf("got MGH data sum=%d, wanted %d", got, want)
	}
}


// The indices in the following lines from R code are 1-based, so substract 1 from them for Golang.
//  expect_equal(vd[100, 100, 100, 1], 77);      # try on command line: mri_info --voxel 99 99 99 inst/extdata/brain.mgz
//  expect_equal(vd[110, 110, 110, 1], 71);
//  expect_equal(vd[1, 1, 1, 1], 0);
//  expect_equal(sum(vd), 121035479);