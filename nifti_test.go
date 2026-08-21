package neuro

// https://pkg.go.dev/testing

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"testing"
	"unsafe"
)

// approxEqual reports whether two float32 values are approximately equal.
func approxEqual(a, b float32, tol float32) bool {
	return math.Abs(float64(a-b)) <= float64(tol)
}

// approxSlice reports whether two float32 slices are approximately equal element-wise.
func approxSlice(a, b []float32, tol float32) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !approxEqual(a[i], b[i], tol) {
			return false
		}
	}
	return true
}

func TestNiftiHeaderSizeIs348(t *testing.T) {
	// The NIfTI-1 header must be exactly 348 bytes on disk. Since encoding/binary
	// serializes struct fields in declaration order without padding, this guards
	// the Go struct against accidental field reordering.
	if size := unsafe.Sizeof(Nifti1Header{}); size != 348 {
		t.Errorf("got unsafe.Sizeof(Nifti1Header{}) = %d, wanted 348", size)
	}
}

func TestRotationQuaternionRoundtrip(t *testing.T) {

	// A set of rotation matrices that must be reconstructed exactly by the
	// quaternion conversion (within a small tolerance).
	testMatrices := []struct {
		name string
		r    [3][3]float64
	}{
		{"identity", [3][3]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}},
		{"rotX90", [3][3]float64{{1, 0, 0}, {0, 0, -1}, {0, 1, 0}}},
		{"rotY90", [3][3]float64{{0, 0, 1}, {0, 1, 0}, {-1, 0, 0}}},
		{"rotZ90", [3][3]float64{{0, -1, 0}, {1, 0, 0}, {0, 0, 1}}},
		// The (determinant +1, after qfac flip) rotation of the brain.mgh vox2ras.
		{"brainMgh", [3][3]float64{{-1, 0, 0}, {0, 0, 1}, {0, 1, 0}}},
		// The (determinant +1, after qfac flip) rotation of the brain.nii sform.
		{"brainNii", [3][3]float64{{-1, 0, 0}, {0, 0, -1}, {0, -1, 0}}},
	}

	const tol = 1e-6
	for _, tc := range testMatrices {
		t.Run(tc.name, func(t *testing.T) {
			a, b, c, d := rotationMatrixToQuaternion(tc.r)
			if a < 0 {
				t.Errorf("quaternion scalar component a=%f is negative, NIfTI convention requires a >= 0", a)
			}
			norm := math.Sqrt(a*a + b*b + c*c + d*d)
			if math.Abs(norm-1.0) > 1e-5 {
				t.Errorf("quaternion is not normalized: |q|=%f", norm)
			}
			recon := quaternionToRotationMatrix(a, b, c, d)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if math.Abs(recon[i][j]-tc.r[i][j]) > tol {
						t.Errorf("reconstructed matrix [%d][%d]=%f, wanted %f", i, j, recon[i][j], tc.r[i][j])
					}
				}
			}
		})
	}
}

func TestMghToNifti(t *testing.T) {

	mgh, err := ReadFsMgh("testdata/brain.mgh", "no")
	if err != nil {
		t.Errorf("ReadFsMgh failed: %v", err)
	}

	nii, err := MghToNifti(mgh)
	if err != nil {
		t.Errorf("MghToNifti failed: %v", err)
	}

	hdr := nii.Header

	// Dimensions and data type.
	if hdr.Dim[1] != 256 || hdr.Dim[2] != 256 || hdr.Dim[3] != 256 || hdr.Dim[4] != 1 {
		t.Errorf("got dims (%d,%d,%d,%d), wanted (256,256,256,1)", hdr.Dim[1], hdr.Dim[2], hdr.Dim[3], hdr.Dim[4])
	}
	if hdr.Datatype != NIFTI_DT_UINT8 {
		t.Errorf("got datatype %d, wanted %d (NIFTI_DT_UINT8)", hdr.Datatype, NIFTI_DT_UINT8)
	}
	if hdr.Bitpix != 8 {
		t.Errorf("got bitpix %d, wanted 8", hdr.Bitpix)
	}
	if !reflect.DeepEqual(nii.Data.DataUint8, mgh.Data.DataMriUchar) {
		t.Errorf("NIfTI data differs from MGH data after conversion.")
	}

	// Header fields.
	if hdr.VoxOffset != 352 {
		t.Errorf("got vox_offset %f, wanted 352", hdr.VoxOffset)
	}
	if hdr.SclSlope != 1 || hdr.SclInter != 0 {
		t.Errorf("got scl_slope=%f scl_inter=%f, wanted 1 and 0", hdr.SclSlope, hdr.SclInter)
	}
	if hdr.QformCode != 1 || hdr.SformCode != 1 {
		t.Errorf("got qform_code=%d sform_code=%d, wanted 1 and 1", hdr.QformCode, hdr.SformCode)
	}

	// The sform must equal the MGH vox2ras matrix (Mdc rows + Pxyz_c).
	wantSrowX := [4]float32{mgh.Header.Mdc[0], mgh.Header.Mdc[1], mgh.Header.Mdc[2], mgh.Header.Pxyz_c[0]}
	wantSrowY := [4]float32{mgh.Header.Mdc[3], mgh.Header.Mdc[4], mgh.Header.Mdc[5], mgh.Header.Pxyz_c[1]}
	wantSrowZ := [4]float32{mgh.Header.Mdc[6], mgh.Header.Mdc[7], mgh.Header.Mdc[8], mgh.Header.Pxyz_c[2]}
	if !approxSlice(hdr.SrowX[:], wantSrowX[:], 1e-6) || !approxSlice(hdr.SrowY[:], wantSrowY[:], 1e-6) || !approxSlice(hdr.SrowZ[:], wantSrowZ[:], 1e-6) {
		t.Errorf("got sform rows (%v, %v, %v), wanted (%v, %v, %v)", hdr.SrowX, hdr.SrowY, hdr.SrowZ, wantSrowX, wantSrowY, wantSrowZ)
	}

	// The qform must reconstruct the MGH vox2ras matrix.
	// For brain.mgh, this is the known expected quaternion.
	if !approxEqual(hdr.QuaternB, 0.0, 1e-5) || !approxEqual(hdr.QuaternC, 0.70710677, 1e-5) || !approxEqual(hdr.QuaternD, 0.70710677, 1e-5) {
		t.Errorf("got quaternion (%f, %f, %f), wanted (0, 0.70710677, 0.70710677)", hdr.QuaternB, hdr.QuaternC, hdr.QuaternD)
	}
	if !approxEqual(hdr.QoffsetX, mgh.Header.Pxyz_c[0], 1e-5) || !approxEqual(hdr.QoffsetY, mgh.Header.Pxyz_c[1], 1e-5) || !approxEqual(hdr.QoffsetZ, mgh.Header.Pxyz_c[2], 1e-5) {
		t.Errorf("got qoffset (%f, %f, %f), wanted (%f, %f, %f)", hdr.QoffsetX, hdr.QoffsetY, hdr.QoffsetZ, mgh.Header.Pxyz_c[0], mgh.Header.Pxyz_c[1], mgh.Header.Pxyz_c[2])
	}
	// The brain.mgh vox2ras has a negative determinant, so qfac = pixdim[0] = -1.
	if hdr.Pixdim[0] != -1 {
		t.Errorf("got pixdim[0]=%f, wanted -1 (qfac=-1)", hdr.Pixdim[0])
	}
	if !approxEqual(hdr.Pixdim[1], 1, 1e-4) || !approxEqual(hdr.Pixdim[2], 1, 1e-4) || !approxEqual(hdr.Pixdim[3], 1, 1e-4) {
		t.Errorf("got pixdim[1..3]=(%f,%f,%f), wanted (1,1,1)", hdr.Pixdim[1], hdr.Pixdim[2], hdr.Pixdim[3])
	}
}

func TestMghToNiftiNoRas(t *testing.T) {

	// An MGH volume without valid RAS information.
	mgh := Mgh{
		Header: MghHeader{
			MghVersion:  1,
			Dim1Length:  2,
			Dim2Length:  3,
			Dim3Length:  1,
			Dim4Length:  1,
			MghDataType: MRI_FLOAT,
			RasGoodFlag: 0,
		},
		Data: MghData{MghDataType: MRI_FLOAT, DataMriFloat: []float32{1, 2, 3, 4, 5, 6}},
	}

	nii, err := MghToNifti(mgh)
	if err != nil {
		t.Errorf("MghToNifti failed: %v", err)
	}

	if nii.Header.SformCode != 0 || nii.Header.QformCode != 0 {
		t.Errorf("got sform_code=%d qform_code=%d, wanted 0 and 0 for volume without RAS info", nii.Header.SformCode, nii.Header.QformCode)
	}
	if nii.Header.Pixdim[0] != 1 {
		t.Errorf("got pixdim[0]=%f, wanted 1", nii.Header.Pixdim[0])
	}
}

func TestNiftiToMgh(t *testing.T) {

	nii, err := ReadNifti("testdata/brain.nii")
	if err != nil {
		t.Errorf("ReadNifti failed: %v", err)
	}

	mgh, err := NiftiToMgh(nii)
	if err != nil {
		t.Errorf("NiftiToMgh failed: %v", err)
	}

	// Data and dims.
	if mgh.Header.Dim1Length != 256 || mgh.Header.Dim2Length != 256 || mgh.Header.Dim3Length != 256 || mgh.Header.Dim4Length != 1 {
		t.Errorf("got dims (%d,%d,%d,%d), wanted (256,256,256,1)", mgh.Header.Dim1Length, mgh.Header.Dim2Length, mgh.Header.Dim3Length, mgh.Header.Dim4Length)
	}
	var sum int64 = 0
	for _, v := range mgh.Data.DataMriUchar {
		sum += int64(v)
	}
	if sum != 121035479 {
		t.Errorf("got MGH data sum %d, wanted %d", sum, 121035479)
	}

	// The sform is preferred over the qform, so the MGH vox2ras must equal the sform.
	if mgh.Header.RasGoodFlag != 1 {
		t.Errorf("got ras_good_flag=%d, wanted 1", mgh.Header.RasGoodFlag)
	}
	wantMdc := [9]float32{
		nii.Header.SrowX[0], nii.Header.SrowX[1], nii.Header.SrowX[2],
		nii.Header.SrowY[0], nii.Header.SrowY[1], nii.Header.SrowY[2],
		nii.Header.SrowZ[0], nii.Header.SrowZ[1], nii.Header.SrowZ[2],
	}
	wantPxyz := [3]float32{nii.Header.SrowX[3], nii.Header.SrowY[3], nii.Header.SrowZ[3]}
	if mgh.Header.Mdc != wantMdc {
		t.Errorf("got Mdc %v, wanted %v", mgh.Header.Mdc, wantMdc)
	}
	if mgh.Header.Pxyz_c != wantPxyz {
		t.Errorf("got Pxyz_c %v, wanted %v", mgh.Header.Pxyz_c, wantPxyz)
	}
	if mgh.Header.XSize != 1 || mgh.Header.YSize != 1 || mgh.Header.ZSize != 1 {
		t.Errorf("got voxel sizes (%f,%f,%f), wanted (1,1,1)", mgh.Header.XSize, mgh.Header.YSize, mgh.Header.ZSize)
	}
}

func TestNiftiSformPreferredOverQform(t *testing.T) {

	nii, err := ReadNifti("testdata/brain.nii")
	if err != nil {
		t.Errorf("ReadNifti failed: %v", err)
	}

	// Set the sform to the identity and keep the (different) qform. The sform
	// must win when converting to MGH.
	nii.Header.SformCode = 1
	nii.Header.SrowX = [4]float32{1, 0, 0, 10}
	nii.Header.SrowY = [4]float32{0, 1, 0, 20}
	nii.Header.SrowZ = [4]float32{0, 0, 1, 30}

	mgh, err := NiftiToMgh(nii)
	if err != nil {
		t.Errorf("NiftiToMgh failed: %v", err)
	}

	wantMdc := [9]float32{1, 0, 0, 0, 1, 0, 0, 0, 1}
	wantPxyz := [3]float32{10, 20, 30}
	if mgh.Header.Mdc != wantMdc {
		t.Errorf("got Mdc %v, wanted %v (sform must be preferred over qform)", mgh.Header.Mdc, wantMdc)
	}
	if mgh.Header.Pxyz_c != wantPxyz {
		t.Errorf("got Pxyz_c %v, wanted %v", mgh.Header.Pxyz_c, wantPxyz)
	}
}

func TestNiftiQformOnly(t *testing.T) {

	nii, err := ReadNifti("testdata/brain.nii")
	if err != nil {
		t.Errorf("ReadNifti failed: %v", err)
	}

	// Disable the sform so that the qform must be used.
	nii.Header.SformCode = 0

	mgh, err := NiftiToMgh(nii)
	if err != nil {
		t.Errorf("NiftiToMgh failed: %v", err)
	}

	if mgh.Header.RasGoodFlag != 1 {
		t.Errorf("got ras_good_flag=%d, wanted 1", mgh.Header.RasGoodFlag)
	}

	// The qform of brain.nii reconstructs the rotation [-1,0,0; 0,0,1; 0,-1,0]
	// with qfac=-1 (pixdim[0] < 0), scaled by the (1,1,1) voxel sizes. The
	// tolerance accounts for the float32 precision of the stored quaternion.
	wantMdc := [9]float32{-1, 0, 0, 0, 0, 1, 0, -1, 0}
	wantPxyz := [3]float32{nii.Header.QoffsetX, nii.Header.QoffsetY, nii.Header.QoffsetZ}
	if !approxSlice(mgh.Header.Mdc[:], wantMdc[:], 1e-3) {
		t.Errorf("got Mdc %v, wanted %v (within tolerance)", mgh.Header.Mdc, wantMdc)
	}
	if !approxSlice(mgh.Header.Pxyz_c[:], wantPxyz[:], 1e-3) {
		t.Errorf("got Pxyz_c %v, wanted %v", mgh.Header.Pxyz_c, wantPxyz)
	}
}

func TestMghToNiftiToMghRoundtrip(t *testing.T) {

	mgh, err := ReadFsMgh("testdata/brain.mgh", "no")
	if err != nil {
		t.Errorf("ReadFsMgh failed: %v", err)
	}

	nii, err := MghToNifti(mgh)
	if err != nil {
		t.Errorf("MghToNifti failed: %v", err)
	}

	// Write and reread to exercise the full file I/O path.
	file, err := os.CreateTemp("", "")
	if err != nil {
		t.Errorf("CreateTemp failed: %v", err)
	}
	defer os.Remove(file.Name())
	niiFileName := file.Name()
	file.Close()

	if err := WriteNifti(niiFileName, nii, "no"); err != nil {
		t.Errorf("WriteNifti failed: %v", err)
	}
	niiReread, err := ReadNifti(niiFileName)
	if err != nil {
		t.Errorf("ReadNifti failed: %v", err)
	}
	mghBack, err := NiftiToMgh(niiReread)
	if err != nil {
		t.Errorf("NiftiToMgh failed: %v", err)
	}

	// The data must survive the roundtrip exactly.
	if !reflect.DeepEqual(mgh.Data.DataMriUchar, mghBack.Data.DataMriUchar) {
		t.Errorf("MGH data differs after MGH -> NIfTI -> MGH roundtrip.")
	}

	// The vox2ras matrix (Mdc and Pxyz_c) must survive the roundtrip (within a
	// small tolerance due to the float32 quaternion representation).
	for i := 0; i < 9; i++ {
		if !approxEqual(mgh.Header.Mdc[i], mghBack.Header.Mdc[i], 1e-3) {
			t.Errorf("Mdc[%d] differs after roundtrip: %f vs %f", i, mgh.Header.Mdc[i], mghBack.Header.Mdc[i])
		}
	}
	for i := 0; i < 3; i++ {
		if !approxEqual(mgh.Header.Pxyz_c[i], mghBack.Header.Pxyz_c[i], 1e-3) {
			t.Errorf("Pxyz_c[%d] differs after roundtrip: %f vs %f", i, mgh.Header.Pxyz_c[i], mghBack.Header.Pxyz_c[i])
		}
	}
}

func ExampleMghToNifti() {
	var mghFile string = "testdata/brain.mgh"

	// Read an MGH volume and convert it to a Nifti struct.
	mgh, _ := ReadFsMgh(mghFile, "no")
	nii, _ := MghToNifti(mgh)

	fmt.Printf("Converted MGH volume to NIfTI with dimensions (%d, %d, %d, %d) and qform_code=%d sform_code=%d.\n",
		int(nii.Header.Dim[1]), int(nii.Header.Dim[2]), int(nii.Header.Dim[3]), int(nii.Header.Dim[4]),
		nii.Header.QformCode, nii.Header.SformCode)
	// Output: Converted MGH volume to NIfTI with dimensions (256, 256, 256, 1) and qform_code=1 sform_code=1.
}
