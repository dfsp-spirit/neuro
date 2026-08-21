package neuro

// Related software: libfs for C++, see:
// https://github.com/dfsp-spirit/libfs/blob/main/include/libfs.h#L5413 for reading
// NIfTI-1 files, and the NIfTI-1 specification for the file format.

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// getIsGzippedNifti translates special is_gzipped values "nii" and "nii.gz" into
// "no" and "yes", respectively.
//
// Parameters:
//   - isGzipped: whether to treat the file as gzip-compressed. If it is one of "nii"
//     or "nii.gz", it will be translated to "no" and "yes", respectively.
//
// Returns:
//   - string: "yes" if isGzipped is "nii.gz", "no" if isGzipped is "nii", otherwise isGzipped.
func getIsGzippedNifti(isGzipped string) string {
	if isGzipped == "nii" {
		return "no"
	}
	if isGzipped == "nii.gz" {
		return "yes"
	}
	return isGzipped
}

// ReadNifti reads a NIfTI-1 file and returns it as a Nifti struct.
//
// Only standard-conformant single-file NIfTI-1 files (magic "n+1", i.e., .nii files)
// are supported. The header data is stored as-is, and the voxel data is read without
// applying the scaling (scl_slope/scl_inter), so that the Nifti struct is a faithful
// representation of the file on disk. To obtain an Mgh struct with the scaling
// applied and the spatial information converted to the MGH representation, use the
// NiftiToMgh function.
//
// Note: the so-called "FreeSurfer hack" for non-conformant NIfTI-1 files (used to
// store per-vertex data in a NIfTI file) is not supported; only standard-conformant
// NIfTI-1 files are read.
//
// Parameters:
//   - filepath: path to the NIfTI file, e.g. '<subject>/mri/brain.nii'. Gzipped
//     files (file extension .nii.gz or .gz) are also supported.
//
// Returns:
//   - Nifti: a Nifti struct containing the header and data.
//   - error: an error if one occurred
func ReadNifti(filepath string) (Nifti, error) {

	var n Nifti

	isGzipped := getIsGzippedNifti("auto")
	treatGzipped := getIsGzipped(filepath, isGzipped)
	bs, err := readFileIntoByteSlice(filepath, treatGzipped)
	if err != nil {
		return n, fmt.Errorf("ReadNifti: could not read file '%s' into byte slice: %s", filepath, err)
	}
	if len(bs) < 348 {
		return n, fmt.Errorf("ReadNifti: file '%s' is too small (%d bytes) to contain a NIfTI-1 header (348 bytes).", filepath, len(bs))
	}

	// Detect the byte order from the sizeof_hdr field (must be 348).
	var sizeofHdr int32
	if err := binary.Read(bytes.NewReader(bs[0:4]), binary.LittleEndian, &sizeofHdr); err != nil {
		return n, fmt.Errorf("ReadNifti: could not read sizeof_hdr from file '%s': %s", filepath, err)
	}
	var order binary.ByteOrder
	if sizeofHdr == 348 {
		order = binary.LittleEndian
	} else {
		if err := binary.Read(bytes.NewReader(bs[0:4]), binary.BigEndian, &sizeofHdr); err != nil {
			return n, fmt.Errorf("ReadNifti: could not read sizeof_hdr from file '%s': %s", filepath, err)
		}
		if sizeofHdr == 348 {
			order = binary.BigEndian
		} else {
			return n, fmt.Errorf("ReadNifti: file '%s' has invalid sizeof_hdr (%d), expected 348. This is not a valid NIfTI-1 file.", filepath, sizeofHdr)
		}
	}

	// Read the header.
	r := bytes.NewReader(bs)
	if err := binary.Read(r, order, &n.Header); err != nil {
		return n, fmt.Errorf("ReadNifti: could not read NIfTI-1 header from file '%s': %s", filepath, err)
	}

	// Validate the magic string: only single-file NIfTI ("n+1") is supported.
	if !(n.Header.Magic == [4]byte{'n', '+', '1', 0}) {
		return n, fmt.Errorf("ReadNifti: file '%s' has invalid magic string '%q'. Only single-file NIfTI-1 files (magic 'n+1') are supported.", filepath, n.Header.Magic)
	}

	// Validate the number of dimensions.
	if n.Header.Dim[0] < 1 || n.Header.Dim[0] > 7 {
		return n, fmt.Errorf("ReadNifti: file '%s' has invalid number of dimensions dim[0]=%d (must be between 1 and 7).", filepath, n.Header.Dim[0])
	}
	if n.Header.Dim[1] <= 0 {
		return n, fmt.Errorf("ReadNifti: file '%s' has invalid dim[1]=%d. This is not a standard-conformant NIfTI-1 file.", filepath, n.Header.Dim[1])
	}

	// Validate the data type and the bits per voxel.
	if _, err := niftiDtypeToMri(n.Header.Datatype); err != nil {
		return n, fmt.Errorf("ReadNifti: file '%s': %s", filepath, err)
	}
	if n.Header.Bitpix != niftiBitpix(n.Header.Datatype) {
		return n, fmt.Errorf("ReadNifti: file '%s' has inconsistent bitpix=%d for data type %d (expected %d).", filepath, n.Header.Bitpix, n.Header.Datatype, niftiBitpix(n.Header.Datatype))
	}

	// Validate the voxel data offset and the data size.
	if n.Header.VoxOffset < 348 {
		return n, fmt.Errorf("ReadNifti: file '%s' has invalid vox_offset=%f (< 348).", filepath, n.Header.VoxOffset)
	}
	numValues := n.Header.NumValues()
	bytesPerElement := int64(niftiBitpix(n.Header.Datatype) / 8)
	payloadBytes := numValues * bytesPerElement
	dataStart := int64(n.Header.VoxOffset)
	if dataStart+payloadBytes > int64(len(bs)) {
		return n, fmt.Errorf("ReadNifti: file '%s' is truncated: the dimensions require %d bytes of voxel data at offset %d, but the file only has %d bytes.", filepath, payloadBytes, dataStart, len(bs))
	}

	if Verbosity >= 1 {
		fmt.Printf("ReadNifti: file '%s' has dimensions (%d, %d, %d, %d), data type %d, byte order %s, %d voxels.\n",
			filepath, n.Header.Dim[1], n.Header.Dim[2], n.Header.Dim[3], n.Header.Dim[4], n.Header.Datatype, orderString(order), numValues)
	}

	// Read the voxel data (without applying scaling).
	n.Data.Datatype = n.Header.Datatype
	dataReader := bytes.NewReader(bs[dataStart : dataStart+payloadBytes])
	switch n.Header.Datatype {
	case NIFTI_DT_UINT8:
		n.Data.DataUint8 = make([]uint8, numValues)
		if err := binary.Read(dataReader, order, &n.Data.DataUint8); err != nil {
			return n, fmt.Errorf("ReadNifti: could not read uint8 voxel data from file '%s': %s", filepath, err)
		}
	case NIFTI_DT_INT16:
		n.Data.DataInt16 = make([]int16, numValues)
		if err := binary.Read(dataReader, order, &n.Data.DataInt16); err != nil {
			return n, fmt.Errorf("ReadNifti: could not read int16 voxel data from file '%s': %s", filepath, err)
		}
	case NIFTI_DT_INT32:
		n.Data.DataInt32 = make([]int32, numValues)
		if err := binary.Read(dataReader, order, &n.Data.DataInt32); err != nil {
			return n, fmt.Errorf("ReadNifti: could not read int32 voxel data from file '%s': %s", filepath, err)
		}
	case NIFTI_DT_FLOAT32:
		n.Data.DataFloat32 = make([]float32, numValues)
		if err := binary.Read(dataReader, order, &n.Data.DataFloat32); err != nil {
			return n, fmt.Errorf("ReadNifti: could not read float32 voxel data from file '%s': %s", filepath, err)
		}
	}

	return n, nil
}

// orderString returns a human-readable name for a binary.ByteOrder.
func orderString(order binary.ByteOrder) string {
	if order == binary.LittleEndian {
		return "little-endian"
	}
	return "big-endian"
}
