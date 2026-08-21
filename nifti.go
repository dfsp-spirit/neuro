package neuro

// Related software: libfs for C++, see:
// https://github.com/dfsp-spirit/libfs/blob/main/include/libfs.h#L5034 for the NIfTI-1
// data type constants and the NIfTI-1 header structure, and the NIfTI-1 specification
// (nifti1.h / nifti1_io.c) for the qform/sform transforms.

import (
	"fmt"
	"math"
)

// NIfTI-1 data type codes supported by this package.
// See the NIfTI-1 specification (nifti1.h) for the full list of codes.
const (
	// NIFTI_DT_UINT8 is the NIfTI-1 data type code for unsigned 8-bit integers.
	NIFTI_DT_UINT8 int16 = 2
	// NIFTI_DT_INT16 is the NIfTI-1 data type code for signed 16-bit integers.
	NIFTI_DT_INT16 int16 = 4
	// NIFTI_DT_INT32 is the NIfTI-1 data type code for signed 32-bit integers.
	NIFTI_DT_INT32 int16 = 8
	// NIFTI_DT_FLOAT32 is the NIfTI-1 data type code for 32-bit floats.
	NIFTI_DT_FLOAT32 int16 = 16
)

// NIFTI_UNITS_MM is the NIfTI-1 xyzt_units code for millimeters (spatial units).
const NIFTI_UNITS_MM byte = 2

// Nifti1Header models the 348-byte header of a NIfTI-1 file.
//
// IMPORTANT: the fields must be declared in exactly the order they appear in the
// NIfTI-1 on-disk format, because encoding/binary.Read and Write serialize struct
// fields in declaration order without inserting padding. See also the test
// TestNiftiHeaderSizeIs348, which guards the total header size.
type Nifti1Header struct {
	SizeofHdr     int32      // Must be 348.
	DataType      [10]byte   // Unused in the NIfTI-1 standard.
	DbName        [18]byte   // Unused in the NIfTI-1 standard.
	Extents       int32      // Unused in the NIfTI-1 standard.
	SessionError  int16      // Unused in the NIfTI-1 standard.
	Regular       byte       // Unused in the NIfTI-1 standard.
	DimInfo       byte       // MRI slice ordering.
	Dim           [8]int16   // Dim[0]=number of dimensions (1-7), Dim[1..7]=extent along each dimension.
	IntentP1      float32    // Intent parameter 1.
	IntentP2      float32    // Intent parameter 2.
	IntentP3      float32    // Intent parameter 3.
	IntentCode    int16      // NIfTI intent code.
	Datatype      int16      // NIfTI data type code (see NIFTI_DT_* constants).
	Bitpix        int16      // Bits per voxel (8, 16, or 32).
	SliceStart    int16      // First slice index.
	Pixdim        [8]float32 // Pixdim[0]=qfac sign, Pixdim[1..7]=voxel sizes per dimension.
	VoxOffset     float32    // Byte offset to the voxel data from the start of the file.
	SclSlope      float32    // Scaling slope (0 means "no scaling", i.e., slope 1).
	SclInter      float32    // Scaling intercept.
	SliceEnd      int16      // Last slice index.
	SliceCode     byte       // Slice timing code.
	XyztUnits     byte       // Units for pixdim[] dimensions.
	CalMax        float32    // Calibrated max.
	CalMin        float32    // Calibrated min.
	SliceDuration float32    // Slice timing duration.
	Toffset       float32    // Time offset.
	Glmax         int32      // Global max (unused).
	Glmin         int32      // Global min (unused).
	Descrip       [80]byte   // Description.
	AuxFile       [24]byte   // Auxiliary filename.
	QformCode     int16      // Quaternion transform code (>0 = valid).
	SformCode     int16      // Affine transform code (>0 = valid).
	QuaternB      float32    // Quaternion b parameter.
	QuaternC      float32    // Quaternion c parameter.
	QuaternD      float32    // Quaternion d parameter.
	QoffsetX      float32    // Quaternion x shift.
	QoffsetY      float32    // Quaternion y shift.
	QoffsetZ      float32    // Quaternion z shift.
	SrowX         [4]float32 // Affine transform row x (sform).
	SrowY         [4]float32 // Affine transform row y (sform).
	SrowZ         [4]float32 // Affine transform row z (sform).
	IntentName    [16]byte   // Intent name.
	Magic         [4]byte    // "n+1\0" (single-file .nii) or "ni1\0" (header/image pair).
}

// NumValues computes the number of voxels based on the Dim header fields.
// Dimensions with value <= 0 are treated as 1, following the NIfTI-1 convention.
//
// Returns:
//   - int64: the number of voxels.
func (h Nifti1Header) NumValues() int64 {
	return int64(dimOrOne(h.Dim[1])) * int64(dimOrOne(h.Dim[2])) *
		int64(dimOrOne(h.Dim[3])) * int64(dimOrOne(h.Dim[4]))
}

// NiftiData models the data part of a NIfTI-1 file. Only the field identified by
// Datatype is valid.
type NiftiData struct {
	Datatype    int16     // The NIfTI data type code, see NIFTI_DT_* constants.
	DataUint8   []uint8   // The data, if Datatype is NIFTI_DT_UINT8.
	DataInt16   []int16   // The data, if Datatype is NIFTI_DT_INT16.
	DataInt32   []int32   // The data, if Datatype is NIFTI_DT_INT32.
	DataFloat32 []float32 // The data, if Datatype is NIFTI_DT_FLOAT32.
}

// Nifti models a full NIfTI-1 file (single-file "n+1" format), including the
// Nifti1Header and the NiftiData. See the separate documentation of the structs
// for details on accessing the fields.
type Nifti struct {
	Header Nifti1Header
	Data   NiftiData
}

// dimOrOne returns d if d > 0, and 1 otherwise. This mirrors the NIfTI-1 convention
// that dimensions beyond the declared number of dimensions may be 0.
func dimOrOne(d int16) int16 {
	if d > 0 {
		return d
	}
	return 1
}

// niftiDtypeToMri maps a NIfTI-1 data type code to the respective MGH MRI_* constant.
//
// Parameters:
//   - niftiDtype: the NIfTI-1 data type code, e.g., NIFTI_DT_UINT8.
//
// Returns:
//   - int32: the MGH MRI data type code (MRI_UCHAR, MRI_SHORT, MRI_INT, or MRI_FLOAT).
//   - error: an error if the NIfTI data type is unsupported.
func niftiDtypeToMri(niftiDtype int16) (int32, error) {
	switch niftiDtype {
	case NIFTI_DT_UINT8:
		return MRI_UCHAR, nil
	case NIFTI_DT_INT16:
		return MRI_SHORT, nil
	case NIFTI_DT_INT32:
		return MRI_INT, nil
	case NIFTI_DT_FLOAT32:
		return MRI_FLOAT, nil
	default:
		return -1, fmt.Errorf("unsupported NIfTI-1 data type code %d. Supported types: UINT8 (2), INT16 (4), INT32 (8), FLOAT32 (16).", niftiDtype)
	}
}

// mriDtypeToNifti maps an MGH MRI_* data type code to the respective NIfTI-1 data type code.
//
// Parameters:
//   - mriDtype: the MGH MRI data type code (MRI_UCHAR, MRI_SHORT, MRI_INT, or MRI_FLOAT).
//
// Returns:
//   - int16: the NIfTI-1 data type code.
//   - error: an error if the MGH data type is unsupported for NIfTI output.
func mriDtypeToNifti(mriDtype int32) (int16, error) {
	switch mriDtype {
	case MRI_UCHAR:
		return NIFTI_DT_UINT8, nil
	case MRI_SHORT:
		return NIFTI_DT_INT16, nil
	case MRI_INT:
		return NIFTI_DT_INT32, nil
	case MRI_FLOAT:
		return NIFTI_DT_FLOAT32, nil
	default:
		return -1, fmt.Errorf("unsupported MGH data type code %d for NIfTI output. Supported types: MRI_UCHAR, MRI_SHORT, MRI_INT, MRI_FLOAT.", mriDtype)
	}
}

// niftiBitpix returns the number of bits per voxel for a supported NIfTI-1 data type code.
//
// Parameters:
//   - niftiDtype: the NIfTI-1 data type code.
//
// Returns:
//   - int16: the number of bits per voxel.
func niftiBitpix(niftiDtype int16) int16 {
	switch niftiDtype {
	case NIFTI_DT_UINT8:
		return 8
	case NIFTI_DT_INT16:
		return 16
	case NIFTI_DT_INT32, NIFTI_DT_FLOAT32:
		return 32
	default:
		return 0
	}
}

// rotationMatrixToQuaternion converts a 3x3 rotation matrix to a unit quaternion
// (a, b, c, d), where a is the scalar ("w") component. This uses the numerically
// robust "largest component" method. The returned quaternion is normalized and has
// a >= 0, so that a = sqrt(1 - b*b - c*c - d*d) holds, as required by the NIfTI-1
// qform convention.
//
// Parameters:
//   - r: a 3x3 rotation matrix (row-major).
//
// Returns:
//   - a, b, c, d: the quaternion components.
func rotationMatrixToQuaternion(r [3][3]float64) (a float64, b float64, c float64, d float64) {

	trace := r[0][0] + r[1][1] + r[2][2]
	if trace > 0 {
		s := math.Sqrt(trace+1.0) * 2.0
		a = 0.25 * s
		b = (r[2][1] - r[1][2]) / s
		c = (r[0][2] - r[2][0]) / s
		d = (r[1][0] - r[0][1]) / s
	} else if r[0][0] > r[1][1] && r[0][0] > r[2][2] {
		s := math.Sqrt(1.0+r[0][0]-r[1][1]-r[2][2]) * 2.0
		a = (r[2][1] - r[1][2]) / s
		b = 0.25 * s
		c = (r[0][1] + r[1][0]) / s
		d = (r[0][2] + r[2][0]) / s
	} else if r[1][1] > r[2][2] {
		s := math.Sqrt(1.0+r[1][1]-r[0][0]-r[2][2]) * 2.0
		a = (r[0][2] - r[2][0]) / s
		b = (r[0][1] + r[1][0]) / s
		c = 0.25 * s
		d = (r[1][2] + r[2][1]) / s
	} else {
		s := math.Sqrt(1.0+r[2][2]-r[0][0]-r[1][1]) * 2.0
		a = (r[1][0] - r[0][1]) / s
		b = (r[0][2] + r[2][0]) / s
		c = (r[1][2] + r[2][1]) / s
		d = 0.25 * s
	}

	// Negating the whole quaternion represents the same rotation. The NIfTI-1
	// convention requires a >= 0 (a is implied as sqrt(1-b*b-c*c-d*d)).
	if a < 0 {
		a, b, c, d = -a, -b, -c, -d
	}

	// Normalize for safety.
	norm := math.Sqrt(a*a + b*b + c*c + d*d)
	if norm > 0 {
		a, b, c, d = a/norm, b/norm, c/norm, d/norm
	}
	return a, b, c, d
}

// quaternionToRotationMatrix converts a quaternion (a, b, c, d) to a 3x3 rotation
// matrix (row-major). This is the inverse of rotationMatrixToQuaternion.
//
// Parameters:
//   - a, b, c, d: the quaternion components.
//
// Returns:
//   - [3][3]float64: the rotation matrix.
func quaternionToRotationMatrix(a float64, b float64, c float64, d float64) [3][3]float64 {
	return [3][3]float64{
		{a*a + b*b - c*c - d*d, 2 * (b*c - a*d), 2 * (b*d + a*c)},
		{2 * (b*c + a*d), a*a + c*c - b*b - d*d, 2 * (c*d - a*b)},
		{2 * (b*d - a*c), 2 * (c*d + a*b), a*a + d*d - b*b - c*c},
	}
}

// vox2rasToQuaternionAndPixdim decomposes the 3x3 part of the MGH vox2ras matrix
// (the Mdc header field) into the NIfTI-1 qform parameters: the quaternion
// components (b, c, d), the qfac sign, and the voxel sizes. This follows the
// NIfTI-1 standard algorithm (nifti_mat44_to_quatern): the columns of the matrix
// are normalized to obtain the rotation part, the determinant decides the qfac
// sign (a negative determinant means the third column must be flipped), and the
// resulting rotation matrix is converted to a unit quaternion.
//
// Parameters:
//   - mdc: the 9 float values of the MGH Mdc field (row-major 3x3 matrix).
//
// Returns:
//   - b, c, d: the quaternion components to store in the NIfTI-1 header.
//   - qfac: the qfac sign (+1 or -1), which determines the sign of Pixdim[0].
//   - pixdim: the voxel sizes, computed as the column norms of the matrix.
func vox2rasToQuaternionAndPixdim(mdc [9]float32) (b float32, c float32, d float32, qfac float32, pixdim [3]float32) {

	var r [3][3]float64
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r[i][j] = float64(mdc[i*3+j])
		}
	}

	// Compute the lengths of the three columns of the matrix.
	colNorm := func(col int) float64 {
		return math.Sqrt(r[0][col]*r[0][col] + r[1][col]*r[1][col] + r[2][col]*r[2][col])
	}
	xd := colNorm(0)
	yd := colNorm(1)
	zd := colNorm(2)

	// Normalize the columns to obtain the pure rotation part.
	for i := 0; i < 3; i++ {
		r[i][0] /= xd
		r[i][1] /= yd
		r[i][2] /= zd
	}

	// Compute the determinant and flip the third column if it is negative.
	det := r[0][0]*(r[1][1]*r[2][2]-r[1][2]*r[2][1]) -
		r[0][1]*(r[1][0]*r[2][2]-r[1][2]*r[2][0]) +
		r[0][2]*(r[1][0]*r[2][1]-r[1][1]*r[2][0])
	qfac = 1.0
	if det < 0 {
		qfac = -1.0
		r[0][2] = -r[0][2]
		r[1][2] = -r[1][2]
		r[2][2] = -r[2][2]
	}

	// Convert the rotation matrix to a quaternion.
	_, qb, qc, qd := rotationMatrixToQuaternion(r)

	return float32(qb), float32(qc), float32(qd), float32(qfac), [3]float32{float32(xd), float32(yd), float32(zd)}
}

// MghToNifti converts an Mgh struct (the native in-memory representation used by
// this package) to a Nifti struct.
//
// The NIfTI-1 header is built from the MGH header, including a proper computation
// of the sform and qform transforms from the MGH vox2ras matrix (the Mdc and Pxyz_c
// header fields): the sform is set directly from the vox2ras matrix, and the qform
// is computed by decomposing the vox2ras matrix into voxel sizes and a unit
// quaternion, following the NIfTI-1 standard algorithm. The voxel data is copied
// unchanged (MGH stores real values, so the NIfTI scaling fields are set to
// scl_slope=1 and scl_inter=0).
//
// Parameters:
//   - mgh: the Mgh struct to convert.
//
// Returns:
//   - Nifti: the converted Nifti struct.
//   - error: an error if one occurred, e.g., if a dimension exceeds the NIfTI-1
//     int16 limit (32767) or the MGH data type is unsupported.
func MghToNifti(mgh Mgh) (Nifti, error) {

	var n Nifti
	hdr := &n.Header

	// Dimensions. NIfTI-1 uses int16 for the dim[] fields.
	dims := [4]int32{mgh.Header.Dim1Length, mgh.Header.Dim2Length, mgh.Header.Dim3Length, mgh.Header.Dim4Length}
	for i, d := range dims {
		if d > 32767 {
			return n, fmt.Errorf("MghToNifti: MGH dimension %d (%d) exceeds the NIfTI-1 int16 limit of 32767.", i+1, d)
		}
	}
	hdr.SizeofHdr = 348
	hdr.Dim[0] = 4
	hdr.Dim[1] = int16(dims[0])
	hdr.Dim[2] = int16(dims[1])
	hdr.Dim[3] = int16(dims[2])
	hdr.Dim[4] = int16(dims[3])
	hdr.Dim[5] = 1
	hdr.Dim[6] = 1
	hdr.Dim[7] = 1

	// Data type.
	niftiDtype, err := mriDtypeToNifti(mgh.Header.MghDataType)
	if err != nil {
		return n, err
	}
	hdr.Datatype = niftiDtype
	hdr.Bitpix = niftiBitpix(niftiDtype)

	// Voxel sizes and other header fields.
	hdr.Pixdim[0] = 1.0
	hdr.Pixdim[1] = positiveOrOne(mgh.Header.XSize)
	hdr.Pixdim[2] = positiveOrOne(mgh.Header.YSize)
	hdr.Pixdim[3] = positiveOrOne(mgh.Header.ZSize)
	hdr.Pixdim[4] = 1.0
	hdr.Pixdim[5] = 1.0
	hdr.Pixdim[6] = 1.0
	hdr.Pixdim[7] = 1.0
	hdr.VoxOffset = 352 // 348-byte header + 4-byte extension indicator.
	hdr.SclSlope = 1.0
	hdr.SclInter = 0.0
	hdr.XyztUnits = NIFTI_UNITS_MM
	copy(hdr.Magic[:], "n+1\x00")

	// Copy the data.
	n.Data.Datatype = niftiDtype
	switch mgh.Header.MghDataType {
	case MRI_UCHAR:
		n.Data.DataUint8 = append([]uint8(nil), mgh.Data.DataMriUchar...)
	case MRI_SHORT:
		n.Data.DataInt16 = append([]int16(nil), mgh.Data.DataMriShort...)
	case MRI_INT:
		n.Data.DataInt32 = append([]int32(nil), mgh.Data.DataMriInt...)
	case MRI_FLOAT:
		n.Data.DataFloat32 = append([]float32(nil), mgh.Data.DataMriFloat...)
	}

	// Spatial transform: compute the sform and qform from the MGH vox2ras matrix.
	if mgh.Header.RasGoodFlag == 1 {
		hdr.SformCode = 1 // Scanner Anatomical.
		hdr.QformCode = 1 // Scanner Anatomical.

		// sform: the rows of the vox2ras matrix.
		hdr.SrowX = [4]float32{mgh.Header.Mdc[0], mgh.Header.Mdc[1], mgh.Header.Mdc[2], mgh.Header.Pxyz_c[0]}
		hdr.SrowY = [4]float32{mgh.Header.Mdc[3], mgh.Header.Mdc[4], mgh.Header.Mdc[5], mgh.Header.Pxyz_c[1]}
		hdr.SrowZ = [4]float32{mgh.Header.Mdc[6], mgh.Header.Mdc[7], mgh.Header.Mdc[8], mgh.Header.Pxyz_c[2]}

		// qform: decompose the vox2ras matrix into voxel sizes and a unit quaternion.
		qb, qc, qd, qfac, qformPixdim := vox2rasToQuaternionAndPixdim(mgh.Header.Mdc)
		hdr.QuaternB = qb
		hdr.QuaternC = qc
		hdr.QuaternD = qd
		hdr.QoffsetX = mgh.Header.Pxyz_c[0]
		hdr.QoffsetY = mgh.Header.Pxyz_c[1]
		hdr.QoffsetZ = mgh.Header.Pxyz_c[2]
		hdr.Pixdim[0] = qfac
		hdr.Pixdim[1] = qformPixdim[0]
		hdr.Pixdim[2] = qformPixdim[1]
		hdr.Pixdim[3] = qformPixdim[2]
	}

	return n, nil
}

// NiftiToMgh converts a Nifti struct to an Mgh struct (the native in-memory
// representation used by this package).
//
// The MGH header is built from the NIfTI-1 header. The spatial information is
// extracted preferring the sform (affine) transform over the qform (quaternion)
// transform; if neither is set, only the voxel sizes are stored and the MGH RAS
// good flag is set to 0. The NIfTI scaling fields (scl_slope and scl_inter) are
// applied to the voxel data (a slope of 0 means no scaling, i.e., slope 1), since
// MGH stores real values without a separate scaling.
//
// Parameters:
//   - n: the Nifti struct to convert.
//
// Returns:
//   - Mgh: the converted Mgh struct.
//   - error: an error if one occurred, e.g., if the NIfTI data type is unsupported
//     or the data size does not match the dimensions in the header.
func NiftiToMgh(n Nifti) (Mgh, error) {

	var mgh Mgh
	mgh.Header.MghVersion = 1
	mgh.Header.DoF = 0

	mriDtype, err := niftiDtypeToMri(n.Header.Datatype)
	if err != nil {
		return mgh, err
	}
	mgh.Header.MghDataType = mriDtype
	mgh.Header.Dim1Length = int32(dimOrOne(n.Header.Dim[1]))
	mgh.Header.Dim2Length = int32(dimOrOne(n.Header.Dim[2]))
	mgh.Header.Dim3Length = int32(dimOrOne(n.Header.Dim[3]))
	mgh.Header.Dim4Length = int32(dimOrOne(n.Header.Dim[4]))

	// Extract the spatial information, preferring sform over qform.
	switch {
	case n.Header.SformCode > 0:
		mgh.Header.RasGoodFlag = 1
		mgh.Header.Mdc = [9]float32{
			n.Header.SrowX[0], n.Header.SrowX[1], n.Header.SrowX[2],
			n.Header.SrowY[0], n.Header.SrowY[1], n.Header.SrowY[2],
			n.Header.SrowZ[0], n.Header.SrowZ[1], n.Header.SrowZ[2],
		}
		mgh.Header.Pxyz_c = [3]float32{n.Header.SrowX[3], n.Header.SrowY[3], n.Header.SrowZ[3]}
		mgh.Header.XSize = n.Header.Pixdim[1]
		mgh.Header.YSize = n.Header.Pixdim[2]
		mgh.Header.ZSize = n.Header.Pixdim[3]
	case n.Header.QformCode > 0:
		mgh.Header.RasGoodFlag = 1
		qfac := 1.0
		if n.Header.Pixdim[0] < 0 {
			qfac = -1.0
		}
		b := float64(n.Header.QuaternB)
		c := float64(n.Header.QuaternC)
		d := float64(n.Header.QuaternD)
		a := math.Sqrt(math.Max(0.0, 1.0-b*b-c*c-d*d))
		R := quaternionToRotationMatrix(a, b, c, d)
		sx := float64(n.Header.Pixdim[1])
		sy := float64(n.Header.Pixdim[2])
		sz := float64(n.Header.Pixdim[3]) * qfac
		mgh.Header.Mdc = [9]float32{
			float32(R[0][0] * sx), float32(R[0][1] * sy), float32(R[0][2] * sz),
			float32(R[1][0] * sx), float32(R[1][1] * sy), float32(R[1][2] * sz),
			float32(R[2][0] * sx), float32(R[2][1] * sy), float32(R[2][2] * sz),
		}
		mgh.Header.Pxyz_c = [3]float32{n.Header.QoffsetX, n.Header.QoffsetY, n.Header.QoffsetZ}
		mgh.Header.XSize = n.Header.Pixdim[1]
		mgh.Header.YSize = n.Header.Pixdim[2]
		mgh.Header.ZSize = n.Header.Pixdim[3]
	default:
		mgh.Header.RasGoodFlag = 0
		mgh.Header.XSize = n.Header.Pixdim[1]
		mgh.Header.YSize = n.Header.Pixdim[2]
		mgh.Header.ZSize = n.Header.Pixdim[3]
	}

	// Copy the data, applying the NIfTI scaling (slope/intercept).
	slope := float64(n.Header.SclSlope)
	if slope == 0.0 {
		slope = 1.0
	}
	inter := float64(n.Header.SclInter)
	numValues := n.Header.NumValues()

	switch mriDtype {
	case MRI_UCHAR:
		if int64(len(n.Data.DataUint8)) != numValues {
			return mgh, fmt.Errorf("NiftiToMgh: NIfTI data has %d uint8 values, but the header dimensions declare %d.", len(n.Data.DataUint8), numValues)
		}
		mgh.Data.DataMriUchar = make([]uint8, numValues)
		for i, raw := range n.Data.DataUint8 {
			v := math.Round(float64(raw)*slope + inter)
			v = clampFloat64(v, 0, 255)
			mgh.Data.DataMriUchar[i] = uint8(v)
		}
	case MRI_SHORT:
		if int64(len(n.Data.DataInt16)) != numValues {
			return mgh, fmt.Errorf("NiftiToMgh: NIfTI data has %d int16 values, but the header dimensions declare %d.", len(n.Data.DataInt16), numValues)
		}
		mgh.Data.DataMriShort = make([]int16, numValues)
		for i, raw := range n.Data.DataInt16 {
			v := math.Round(float64(raw)*slope + inter)
			v = clampFloat64(v, -32768, 32767)
			mgh.Data.DataMriShort[i] = int16(v)
		}
	case MRI_INT:
		if int64(len(n.Data.DataInt32)) != numValues {
			return mgh, fmt.Errorf("NiftiToMgh: NIfTI data has %d int32 values, but the header dimensions declare %d.", len(n.Data.DataInt32), numValues)
		}
		mgh.Data.DataMriInt = make([]int32, numValues)
		for i, raw := range n.Data.DataInt32 {
			v := math.Round(float64(raw)*slope + inter)
			mgh.Data.DataMriInt[i] = int32(v)
		}
	case MRI_FLOAT:
		if int64(len(n.Data.DataFloat32)) != numValues {
			return mgh, fmt.Errorf("NiftiToMgh: NIfTI data has %d float32 values, but the header dimensions declare %d.", len(n.Data.DataFloat32), numValues)
		}
		mgh.Data.DataMriFloat = make([]float32, numValues)
		for i, raw := range n.Data.DataFloat32 {
			mgh.Data.DataMriFloat[i] = float32(float64(raw)*slope + inter)
		}
	}
	mgh.Data.MghDataType = mriDtype

	return mgh, nil
}

// positiveOrOne returns v if v > 0, and 1 otherwise. This is used for the voxel
// sizes when converting from MGH, where a voxel size of 0 (or negative) is invalid.
func positiveOrOne(v float32) float32 {
	if v > 0 {
		return v
	}
	return 1.0
}

// clampFloat64 clamps a value to the inclusive range [minVal, maxVal].
func clampFloat64(v float64, minVal float64, maxVal float64) float64 {
	if v < minVal {
		return minVal
	}
	if v > maxVal {
		return maxVal
	}
	return v
}
