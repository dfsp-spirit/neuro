package neuro

// https://pkg.go.dev/testing

import (
	"testing"
)

func TestReadFsMghHeader(t *testing.T) {
	var mghFile string = "testdata/brain.mgh"

	hdr, _ := ReadFsMghHeader(mghFile, "auto")

	got := hdr.MghDataType
	want, _ := getMghDataTypeCode("MRI_UCHAR")

	if got != want {
		t.Errorf("got MGH data type %d, wanted %d", got, want)
	}
}

func TestReadFsMghHeaderRAS(t *testing.T) {
	var mghFile string = "testdata/brain.mgh"

	hdr, _ := ReadFsMghHeader(mghFile, "no")

	got := hdr.RasGoodFlag
	var want int16 = 1

	if got != want {
		t.Errorf("got MGH RasGoodFlag=%d, wanted %d", got, want)
	}
}

func TestReadFsMghFull(t *testing.T) {
	var mghFile string = "testdata/brain.mgh"

	mgh, _ := ReadFsMgh(mghFile, "mgh")

	got := mgh.Header.RasGoodFlag
	var want int16 = 1

	if got != want {
		t.Errorf("got MGH RasGoodFlag=%d, wanted %d", got, want)
	}
}

func TestReadFsMghFullAt0000(t *testing.T) {
	var mghFile string = "testdata/brain.mgh"

	mgh, _ := ReadFsMgh(mghFile, "no")

	got := mgh.Data.DataMriUchar[0]
	var want uint8 = 0

	if got != want {
		t.Errorf("got data value=%d, wanted %d", got, want)
	}
}

func TestReadFsMghFullSum(t *testing.T) {
	var mghFile string = "testdata/brain.mgh"

	mgh, _ := ReadFsMgh(mghFile, "no")

	var sum int = 0
	for _, voxel_val := range mgh.Data.DataMriUchar {
		sum += int(voxel_val)
	}
	var got = sum
	var want int = 121035479 // known from external tests with standard software.

	if got != want {
		t.Errorf("got MGH data sum=%d, wanted %d", got, want)
	}
}

func TestReadFsMgzFullSum(t *testing.T) {
	var mgzFile string = "testdata/brain.mgz"

	mgh, _ := ReadFsMgh(mgzFile, "yes")

	var sum int = 0
	for _, voxel_val := range mgh.Data.DataMriUchar {
		sum += int(voxel_val)
	}
	var got = sum
	var want int = 121035479 // known from external tests with standard software.

	if got != want {
		t.Errorf("got MGH data sum=%d, wanted %d", got, want)
	}
}

func TestReadFsMgzFullSumisGzippedMgz(t *testing.T) {
	var mgzFile string = "testdata/brain.mgz"

	mgh, _ := ReadFsMgh(mgzFile, "mgz") // Use 'mgz' for isGzipped

	var sum int = 0
	for _, voxel_val := range mgh.Data.DataMriUchar {
		sum += int(voxel_val)
	}
	var got = sum
	var want int = 121035479 // known from external tests with standard software.

	if got != want {
		t.Errorf("got MGH data sum=%d, wanted %d", got, want)
	}
}

func TestReadFsMghPervertex(t *testing.T) {
	var mgzFile string = "testdata/lh.thickness.fwhm5.fsaverage.mgh"

	mgh, _ := ReadFsMgh(mgzFile, "auto")

	mean_thickness, _ := mean(mgh.Data.DataMriFloat)

	lower_border := 2.31
	upper_border := 2.33
	if mean_thickness < float32(lower_border) || mean_thickness > float32(upper_border) {
		t.Errorf("got mean thickness=%f, wanted between %f and %f", mean_thickness, lower_border, upper_border)
	}
}
