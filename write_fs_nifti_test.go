package neuro

// https://pkg.go.dev/testing

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

// niftiHeadersMeaningfullyEqual compares the header fields of two Nifti structs that
// are relevant for a read/write roundtrip.
func niftiHeadersMeaningfullyEqual(a, b Nifti1Header) bool {
	return a.SizeofHdr == b.SizeofHdr &&
		a.Dim == b.Dim &&
		a.Datatype == b.Datatype &&
		a.Bitpix == b.Bitpix &&
		a.Pixdim == b.Pixdim &&
		a.VoxOffset == b.VoxOffset &&
		a.SclSlope == b.SclSlope &&
		a.SclInter == b.SclInter &&
		a.QformCode == b.QformCode &&
		a.SformCode == b.SformCode &&
		a.QuaternB == b.QuaternB &&
		a.QuaternC == b.QuaternC &&
		a.QuaternD == b.QuaternD &&
		a.QoffsetX == b.QoffsetX &&
		a.QoffsetY == b.QoffsetY &&
		a.QoffsetZ == b.QoffsetZ &&
		a.SrowX == b.SrowX &&
		a.SrowY == b.SrowY &&
		a.SrowZ == b.SrowZ &&
		a.Magic == b.Magic
}

func TestWriteRereadNifti(t *testing.T) {

	nii, err := ReadNifti("testdata/brain.nii")
	if err != nil {
		t.Errorf("ReadNifti failed: %v", err)
	}

	// get a temp file.
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name()) // clean up
	niiFileName := file.Name()
	file.Close()

	err = WriteNifti(niiFileName, nii, "no")
	if err != nil {
		t.Errorf("WriteNifti failed: %v", err)
	}

	niiReread, err := ReadNifti(niiFileName)
	if err != nil {
		t.Errorf("ReadNifti failed: %v", err)
	}

	if !niftiHeadersMeaningfullyEqual(nii.Header, niiReread.Header) {
		t.Errorf("NIfTI headers differ after roundtrip: %+v vs %+v", nii.Header, niiReread.Header)
	}
	if !reflect.DeepEqual(nii.Data.DataUint8, niiReread.Data.DataUint8) {
		t.Errorf("NIfTI data differs after roundtrip.")
	}
}

func TestWriteRereadNiftiGz(t *testing.T) {

	nii, err := ReadNifti("testdata/brain.nii")
	if err != nil {
		t.Errorf("ReadNifti failed: %v", err)
	}

	// get a temp file with a .nii.gz extension.
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name()) // clean up
	niiGzFileName := file.Name() + ".nii.gz"
	file.Close()

	err = WriteNifti(niiGzFileName, nii, "yes")
	if err != nil {
		t.Errorf("WriteNifti failed: %v", err)
	}

	niiReread, err := ReadNifti(niiGzFileName)
	if err != nil {
		t.Errorf("ReadNifti failed: %v", err)
	}

	if !reflect.DeepEqual(nii.Data.DataUint8, niiReread.Data.DataUint8) {
		t.Errorf("NIfTI data differs after gzip roundtrip.")
	}
	if niiReread.Header.Dim != nii.Header.Dim {
		t.Errorf("NIfTI dimensions differ after gzip roundtrip.")
	}
}

func TestWriteNiftiAllDataTypes(t *testing.T) {

	// Create a small Nifti struct for each supported data type and verify the roundtrip.
	testCases := []struct {
		name      string
		datatype  int16
		makeNifti func() Nifti
	}{
		{
			name:     "NIFTI_DT_UINT8",
			datatype: NIFTI_DT_UINT8,
			makeNifti: func() Nifti {
				data := []uint8{0, 1, 2, 255, 128, 64}
				return Nifti{
					Header: Nifti1Header{Dim: [8]int16{4, int16(len(data)), 1, 1, 1, 1, 1, 1}, Datatype: NIFTI_DT_UINT8},
					Data:   NiftiData{Datatype: NIFTI_DT_UINT8, DataUint8: data},
				}
			},
		},
		{
			name:     "NIFTI_DT_INT16",
			datatype: NIFTI_DT_INT16,
			makeNifti: func() Nifti {
				data := []int16{0, -1, 2, 32767, -32768}
				return Nifti{
					Header: Nifti1Header{Dim: [8]int16{4, int16(len(data)), 1, 1, 1, 1, 1, 1}, Datatype: NIFTI_DT_INT16},
					Data:   NiftiData{Datatype: NIFTI_DT_INT16, DataInt16: data},
				}
			},
		},
		{
			name:     "NIFTI_DT_INT32",
			datatype: NIFTI_DT_INT32,
			makeNifti: func() Nifti {
				data := []int32{0, -1, 2, 100000, -100000}
				return Nifti{
					Header: Nifti1Header{Dim: [8]int16{4, int16(len(data)), 1, 1, 1, 1, 1, 1}, Datatype: NIFTI_DT_INT32},
					Data:   NiftiData{Datatype: NIFTI_DT_INT32, DataInt32: data},
				}
			},
		},
		{
			name:     "NIFTI_DT_FLOAT32",
			datatype: NIFTI_DT_FLOAT32,
			makeNifti: func() Nifti {
				data := []float32{0.0, -1.5, 2.25, 3.14159}
				return Nifti{
					Header: Nifti1Header{Dim: [8]int16{4, int16(len(data)), 1, 1, 1, 1, 1, 1}, Datatype: NIFTI_DT_FLOAT32},
					Data:   NiftiData{Datatype: NIFTI_DT_FLOAT32, DataFloat32: data},
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nii := tc.makeNifti()

			file, err := os.CreateTemp("", "")
			if err != nil {
				t.Errorf("CreateTemp failed: %v", err)
			}
			defer os.Remove(file.Name()) // clean up
			niiFileName := file.Name()
			file.Close()

			if err := WriteNifti(niiFileName, nii, "no"); err != nil {
				t.Errorf("WriteNifti failed: %v", err)
			}

			niiReread, err := ReadNifti(niiFileName)
			if err != nil {
				t.Errorf("ReadNifti failed: %v", err)
			}

			if niiReread.Header.Datatype != nii.Header.Datatype {
				t.Errorf("got datatype %d, wanted %d", niiReread.Header.Datatype, nii.Header.Datatype)
			}
			if !reflect.DeepEqual(nii.Data, niiReread.Data) {
				t.Errorf("NIfTI data differs after roundtrip: %+v vs %+v", nii.Data, niiReread.Data)
			}
		})
	}
}

func TestWriteNiftiDataSizeMismatch(t *testing.T) {

	// Header declares 100 values but data has only 1 -> writing must fail.
	nii := Nifti{
		Header: Nifti1Header{Dim: [8]int16{4, 100, 1, 1, 1, 1, 1, 1}, Datatype: NIFTI_DT_FLOAT32},
		Data:   NiftiData{Datatype: NIFTI_DT_FLOAT32, DataFloat32: []float32{1.0}},
	}

	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name()) // clean up
	niiFileName := file.Name()
	file.Close()

	err = WriteNifti(niiFileName, nii, "no")
	if err == nil {
		t.Errorf("WriteNifti should have failed for mismatched data size, but it succeeded.")
	}
}

func ExampleWriteNifti() {

	// Read an MGH file, convert it to a Nifti struct, write it to a temp file, and read it back.
	mgh, _ := ReadFsMgh("testdata/brain.mgh", "no")
	nii, _ := MghToNifti(mgh)

	file, _ := os.CreateTemp("", "")
	defer os.Remove(file.Name())
	file.Close()

	_ = WriteNifti(file.Name(), nii, "no")
	niiReread, _ := ReadNifti(file.Name())

	fmt.Printf("Wrote and reread a NIfTI file with dimensions (%d, %d, %d, %d).\n",
		int(niiReread.Header.Dim[1]), int(niiReread.Header.Dim[2]), int(niiReread.Header.Dim[3]), int(niiReread.Header.Dim[4]))
	// Output: Wrote and reread a NIfTI file with dimensions (256, 256, 256, 1).
}
