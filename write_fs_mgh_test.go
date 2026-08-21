package neuro

// https://pkg.go.dev/testing

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

// mghHeadersMeaningfullyEqual compares the meaningful fields of two MGH headers.
// The Reserved bytes of the header are not compared, as they are documented to
// contain unused/random data.
func mghHeadersMeaningfullyEqual(a, b MghHeader) bool {
	return a.MghVersion == b.MghVersion &&
		a.Dim1Length == b.Dim1Length &&
		a.Dim2Length == b.Dim2Length &&
		a.Dim3Length == b.Dim3Length &&
		a.Dim4Length == b.Dim4Length &&
		a.MghDataType == b.MghDataType &&
		a.DoF == b.DoF &&
		a.RasGoodFlag == b.RasGoodFlag &&
		a.XSize == b.XSize &&
		a.YSize == b.YSize &&
		a.ZSize == b.ZSize &&
		a.Mdc == b.Mdc &&
		a.Pxyz_c == b.Pxyz_c
}

func TestWriteRereadMgh(t *testing.T) {

	mgh, err := ReadFsMgh("testdata/brain.mgh", "no")
	if err != nil {
		t.Errorf("ReadFsMgh failed: %v", err)
	}

	// get a temp file.
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name()) // clean up
	mghFileName := file.Name()
	file.Close()

	err = WriteFsMgh(mghFileName, mgh, "no")
	if err != nil {
		t.Errorf("WriteFsMgh failed: %v", err)
	}

	mghReread, err := ReadFsMgh(mghFileName, "no")
	if err != nil {
		t.Errorf("ReadFsMgh failed: %v", err)
	}

	if !mghHeadersMeaningfullyEqual(mgh.Header, mghReread.Header) {
		t.Errorf("MGH headers differ after roundtrip: %+v vs %+v", mgh.Header, mghReread.Header)
	}
	if !reflect.DeepEqual(mgh.Data.DataMriUchar, mghReread.Data.DataMriUchar) {
		t.Errorf("MGH data differs after roundtrip.")
	}
}

func TestWriteRereadMgz(t *testing.T) {

	mgh, err := ReadFsMgh("testdata/brain.mgh", "no")
	if err != nil {
		t.Errorf("ReadFsMgh failed: %v", err)
	}

	// get a temp file.
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name()) // clean up
	mgzFileName := file.Name() + ".mgz"
	file.Close()

	err = WriteFsMgh(mgzFileName, mgh, "yes")
	if err != nil {
		t.Errorf("WriteFsMgh failed: %v", err)
	}

	mghReread, err := ReadFsMgh(mgzFileName, "yes")
	if err != nil {
		t.Errorf("ReadFsMgh failed: %v", err)
	}

	if !mghHeadersMeaningfullyEqual(mgh.Header, mghReread.Header) {
		t.Errorf("MGZ headers differ after roundtrip: %+v vs %+v", mgh.Header, mghReread.Header)
	}
	if !reflect.DeepEqual(mgh.Data.DataMriUchar, mghReread.Data.DataMriUchar) {
		t.Errorf("MGZ data differs after roundtrip.")
	}
}

func TestWriteRereadMghAuto(t *testing.T) {

	// Writing with isGzipped="auto" should use the file extension to decide
	// between MGH (no compression) and MGZ (gzip compression).
	mgh, err := ReadFsMgh("testdata/brain.mgh", "no")
	if err != nil {
		t.Errorf("ReadFsMgh failed: %v", err)
	}

	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name()) // clean up
	mghFileName := file.Name() + ".mgh"
	file.Close()

	err = WriteFsMgh(mghFileName, mgh, "auto")
	if err != nil {
		t.Errorf("WriteFsMgh failed: %v", err)
	}

	mghReread, err := ReadFsMgh(mghFileName, "auto")
	if err != nil {
		t.Errorf("ReadFsMgh failed: %v", err)
	}
	if !reflect.DeepEqual(mgh.Data.DataMriUchar, mghReread.Data.DataMriUchar) {
		t.Errorf("MGH data differs after roundtrip.")
	}
}

func TestWriteMghAllDataTypes(t *testing.T) {

	// Create a small Mgh struct for each supported data type and verify the roundtrip.
	testCases := []struct {
		name     string
		dtype    int32
		makeData func() Mgh
	}{
		{
			name:  "MRI_UCHAR",
			dtype: MRI_UCHAR,
			makeData: func() Mgh {
				data := []uint8{0, 1, 2, 255, 128, 64, 32, 16}
				return Mgh{
					Header: MghHeader{
						MghVersion:  1,
						Dim1Length:  int32(len(data)),
						Dim2Length:  1,
						Dim3Length:  1,
						Dim4Length:  1,
						MghDataType: MRI_UCHAR,
						RasGoodFlag: 0,
					},
					Data: MghData{MghDataType: MRI_UCHAR, DataMriUchar: data},
				}
			},
		},
		{
			name:  "MRI_INT",
			dtype: MRI_INT,
			makeData: func() Mgh {
				data := []int32{0, -1, 2, 100000, -100000, 42}
				return Mgh{
					Header: MghHeader{
						MghVersion:  1,
						Dim1Length:  int32(len(data)),
						Dim2Length:  1,
						Dim3Length:  1,
						Dim4Length:  1,
						MghDataType: MRI_INT,
						RasGoodFlag: 0,
					},
					Data: MghData{MghDataType: MRI_INT, DataMriInt: data},
				}
			},
		},
		{
			name:  "MRI_FLOAT",
			dtype: MRI_FLOAT,
			makeData: func() Mgh {
				data := []float32{0.0, -1.5, 2.25, 3.14159, 1e-6, 1e6}
				return Mgh{
					Header: MghHeader{
						MghVersion:  1,
						Dim1Length:  int32(len(data)),
						Dim2Length:  1,
						Dim3Length:  1,
						Dim4Length:  1,
						MghDataType: MRI_FLOAT,
						RasGoodFlag: 0,
					},
					Data: MghData{MghDataType: MRI_FLOAT, DataMriFloat: data},
				}
			},
		},
		{
			name:  "MRI_SHORT",
			dtype: MRI_SHORT,
			makeData: func() Mgh {
				data := []int16{0, -1, 2, 32767, -32768}
				return Mgh{
					Header: MghHeader{
						MghVersion:  1,
						Dim1Length:  int32(len(data)),
						Dim2Length:  1,
						Dim3Length:  1,
						Dim4Length:  1,
						MghDataType: MRI_SHORT,
						RasGoodFlag: 0,
					},
					Data: MghData{MghDataType: MRI_SHORT, DataMriShort: data},
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mgh := tc.makeData()

			file, err := os.CreateTemp("", "")
			if err != nil {
				t.Errorf("CreateTemp failed: %v", err)
			}
			defer os.Remove(file.Name()) // clean up
			mghFileName := file.Name()
			file.Close()

			if err := WriteFsMgh(mghFileName, mgh, "no"); err != nil {
				t.Errorf("WriteFsMgh failed: %v", err)
			}

			mghReread, err := ReadFsMgh(mghFileName, "no")
			if err != nil {
				t.Errorf("ReadFsMgh failed: %v", err)
			}

			if mghReread.Header.MghDataType != mgh.Header.MghDataType {
				t.Errorf("got data type %d, wanted %d", mghReread.Header.MghDataType, mgh.Header.MghDataType)
			}
			if !reflect.DeepEqual(mgh.Data, mghReread.Data) {
				t.Errorf("MGH data differs after roundtrip: %+v vs %+v", mgh.Data, mghReread.Data)
			}
		})
	}
}

func TestWriteMghDataSizeMismatch(t *testing.T) {

	// Header declares 100 values but data has only 1 -> writing must fail.
	mgh := Mgh{
		Header: MghHeader{
			MghVersion:  1,
			Dim1Length:  100,
			Dim2Length:  1,
			Dim3Length:  1,
			Dim4Length:  1,
			MghDataType: MRI_FLOAT,
			RasGoodFlag: 0,
		},
		Data: MghData{MghDataType: MRI_FLOAT, DataMriFloat: []float32{1.0}},
	}

	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name()) // clean up
	mghFileName := file.Name()
	file.Close()

	err = WriteFsMgh(mghFileName, mgh, "no")
	if err == nil {
		t.Errorf("WriteFsMgh should have failed for mismatched data size, but it succeeded.")
	}
}

func ExampleWriteFsMgh() {

	// Read an MGH file, write it to a temp file, and read it back.
	mgh, _ := ReadFsMgh("testdata/brain.mgh", "no")

	file, _ := os.CreateTemp("", "")
	defer os.Remove(file.Name())
	file.Close()

	_ = WriteFsMgh(file.Name(), mgh, "no")
	mghReread, _ := ReadFsMgh(file.Name(), "no")

	fmt.Printf("Read Mgh with dimensions (%d, %d, %d, %d), roundtrip preserved %d data values.\n",
		int(mgh.Header.Dim1Length), int(mgh.Header.Dim2Length), int(mgh.Header.Dim3Length), int(mgh.Header.Dim4Length),
		len(mghReread.Data.DataMriUchar))
	// Output: Read Mgh with dimensions (256, 256, 256, 1), roundtrip preserved 16777216 data values.
}
