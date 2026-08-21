package neuro

// https://pkg.go.dev/testing

import (
	"fmt"
	"os"
	"testing"
)

func TestReadNifti(t *testing.T) {
	var niiFile string = "testdata/brain.nii"

	nii, err := ReadNifti(niiFile)
	if err != nil {
		t.Errorf("ReadNifti failed: %v", err)
	}

	hdr := nii.Header

	// Dimensions: a 256x256x256 volume.
	if hdr.Dim[0] != 3 {
		t.Errorf("got dim[0]=%d, wanted 3", hdr.Dim[0])
	}
	wantDims := [4]int16{256, 256, 256, 1}
	gotDims := [4]int16{hdr.Dim[1], hdr.Dim[2], hdr.Dim[3], hdr.Dim[4]}
	if gotDims != wantDims {
		t.Errorf("got dims %v, wanted %v", gotDims, wantDims)
	}

	// Data type: uint8.
	if hdr.Datatype != NIFTI_DT_UINT8 {
		t.Errorf("got datatype %d, wanted %d (NIFTI_DT_UINT8)", hdr.Datatype, NIFTI_DT_UINT8)
	}
	if len(nii.Data.DataUint8) != 256*256*256 {
		t.Errorf("got %d data values, wanted %d", len(nii.Data.DataUint8), 256*256*256)
	}

	// The data sum is known to match the corresponding MGH file (brain.mgh).
	var sum int64 = 0
	for _, v := range nii.Data.DataUint8 {
		sum += int64(v)
	}
	if sum != 121035479 {
		t.Errorf("got NIfTI data sum %d, wanted %d", sum, 121035479)
	}

	// Header fields of the test file.
	if hdr.VoxOffset != 352 {
		t.Errorf("got vox_offset %f, wanted 352", hdr.VoxOffset)
	}
	if hdr.Pixdim[0] != -1 {
		t.Errorf("got pixdim[0]=%f, wanted -1 (qfac=-1)", hdr.Pixdim[0])
	}
	if hdr.QformCode != 1 || hdr.SformCode != 1 {
		t.Errorf("got qform_code=%d sform_code=%d, wanted 1 and 1", hdr.QformCode, hdr.SformCode)
	}
	if !approxEqual(hdr.QuaternB, 0.0, 1e-5) || !approxEqual(hdr.QuaternC, 0.70710677, 1e-5) || !approxEqual(hdr.QuaternD, -0.70710677, 1e-5) {
		t.Errorf("got quaternion (%f, %f, %f), wanted (0, 0.70710677, -0.70710677)", hdr.QuaternB, hdr.QuaternC, hdr.QuaternD)
	}
	if !approxEqual(hdr.QoffsetX, 127.500046, 1e-3) || !approxEqual(hdr.QoffsetY, -98.62726, 1e-3) || !approxEqual(hdr.QoffsetZ, 79.09527, 1e-3) {
		t.Errorf("got qoffset (%f, %f, %f), wanted (127.500046, -98.62726, 79.09527)", hdr.QoffsetX, hdr.QoffsetY, hdr.QoffsetZ)
	}
	wantSrowX := [4]float32{-1, 0, 0, 127.500046}
	wantSrowY := [4]float32{0, 0, 1, -98.62726}
	wantSrowZ := [4]float32{0, -1, 0, 79.09527}
	if !approxSlice(hdr.SrowX[:], wantSrowX[:], 1e-3) || !approxSlice(hdr.SrowY[:], wantSrowY[:], 1e-3) || !approxSlice(hdr.SrowZ[:], wantSrowZ[:], 1e-3) {
		t.Errorf("got sform rows (%v, %v, %v), wanted (%v, %v, %v)", hdr.SrowX, hdr.SrowY, hdr.SrowZ, wantSrowX, wantSrowY, wantSrowZ)
	}
	if hdr.Magic != [4]byte{'n', '+', '1', 0} {
		t.Errorf("got magic %q, wanted 'n+1'", hdr.Magic)
	}
}

func TestReadNiftiNonexistentFile(t *testing.T) {
	_, err := ReadNifti("testdata/doesnotexist.nii")
	if err == nil {
		t.Errorf("ReadNifti should have failed for a nonexistent file, but it succeeded.")
	}
}

func TestReadNiftiTooSmall(t *testing.T) {
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name())
	if _, err := file.Write([]byte("too small")); err != nil {
		t.Errorf("Write failed: %v", err)
	}
	file.Close()

	_, err = ReadNifti(file.Name())
	if err == nil {
		t.Errorf("ReadNifti should have failed for a too-small file, but it succeeded.")
	}
}

func TestReadNiftiInvalidMagic(t *testing.T) {

	// Build a file with a valid sizeof_hdr but an invalid magic string.
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name())
	file.Close()

	// Write a 352-byte file: little-endian sizeof_hdr=348 and a bad magic string.
	bs := make([]byte, 352)
	bs[0] = 0x5c // 348 little-endian: 0x5c 0x01 0x00 0x00
	bs[1] = 0x01
	copy(bs[344:348], []byte("xxxx"))
	if err := os.WriteFile(file.Name(), bs, 0644); err != nil {
		t.Errorf("WriteFile failed: %v", err)
	}

	_, err = ReadNifti(file.Name())
	if err == nil {
		t.Errorf("ReadNifti should have failed for an invalid magic string, but it succeeded.")
	}
}

func ExampleReadNifti() {
	var niiFile string = "testdata/brain.nii"

	// Read the NIfTI file.
	nii, _ := ReadNifti(niiFile)

	fmt.Printf("Read NIfTI with dimensions (%d, %d, %d, %d) and %d voxels from file '%s'.\n",
		int(nii.Header.Dim[1]), int(nii.Header.Dim[2]), int(nii.Header.Dim[3]), int(nii.Header.Dim[4]),
		len(nii.Data.DataUint8), niiFile)
	// Output: Read NIfTI with dimensions (256, 256, 256, 1) and 16777216 voxels from file 'testdata/brain.nii'.
}
