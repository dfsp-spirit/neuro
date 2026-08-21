package neuro

// Related software: libfs for C++, see:
// https://github.com/dfsp-spirit/libfs/blob/main/include/libfs.h#L5594 for writing
// NIfTI-1 files, and the NIfTI-1 specification for the file format.

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// niftiToByteSlice serializes a Nifti struct into the binary NIfTI-1 format
// (single-file, big-endian, not compressed).
//
// Parameters:
//   - n: the Nifti struct to serialize.
//
// Returns:
//   - []byte: the byte representation of the NIfTI-1 file content.
//   - error: an error if one occurred, e.g., if a dimension exceeds the NIfTI-1
//     int16 limit or the data type is unsupported.
func niftiToByteSlice(n Nifti) ([]byte, error) {

	endian := binary.BigEndian // The NIfTI-1 standard byte order for files on disk.

	// Validate the dimensions and the data type.
	for i := 1; i <= 4; i++ {
		if n.Header.Dim[i] < 0 {
			return nil, fmt.Errorf("niftiToByteSlice: negative dimension dim[%d]=%d.", i, n.Header.Dim[i])
		}
	}
	if _, err := niftiDtypeToMri(n.Header.Datatype); err != nil {
		return nil, err
	}

	// Normalize a few header fields so the written file is always conformant.
	n.Header.SizeofHdr = 348
	n.Header.VoxOffset = 352 // 348-byte header + 4-byte extension indicator.
	n.Header.Bitpix = niftiBitpix(n.Header.Datatype)
	copy(n.Header.Magic[:], "n+1\x00")

	var buf bytes.Buffer

	// Write the header.
	if err := binary.Write(&buf, endian, n.Header); err != nil {
		return nil, fmt.Errorf("niftiToByteSlice: could not write NIfTI-1 header: %s", err)
	}

	// Write the 4-byte extension indicator (0 = no extensions).
	if err := binary.Write(&buf, endian, int32(0)); err != nil {
		return nil, fmt.Errorf("niftiToByteSlice: could not write extension indicator: %s", err)
	}

	// Validate the data size against the header and write the voxel data.
	numValues := n.Header.NumValues()
	switch n.Header.Datatype {
	case NIFTI_DT_UINT8:
		if int64(len(n.Data.DataUint8)) != numValues {
			return nil, fmt.Errorf("niftiToByteSlice: uint8 data has %d values, but the header dimensions declare %d.", len(n.Data.DataUint8), numValues)
		}
		if err := binary.Write(&buf, endian, n.Data.DataUint8); err != nil {
			return nil, fmt.Errorf("niftiToByteSlice: could not write uint8 voxel data: %s", err)
		}
	case NIFTI_DT_INT16:
		if int64(len(n.Data.DataInt16)) != numValues {
			return nil, fmt.Errorf("niftiToByteSlice: int16 data has %d values, but the header dimensions declare %d.", len(n.Data.DataInt16), numValues)
		}
		if err := binary.Write(&buf, endian, n.Data.DataInt16); err != nil {
			return nil, fmt.Errorf("niftiToByteSlice: could not write int16 voxel data: %s", err)
		}
	case NIFTI_DT_INT32:
		if int64(len(n.Data.DataInt32)) != numValues {
			return nil, fmt.Errorf("niftiToByteSlice: int32 data has %d values, but the header dimensions declare %d.", len(n.Data.DataInt32), numValues)
		}
		if err := binary.Write(&buf, endian, n.Data.DataInt32); err != nil {
			return nil, fmt.Errorf("niftiToByteSlice: could not write int32 voxel data: %s", err)
		}
	case NIFTI_DT_FLOAT32:
		if int64(len(n.Data.DataFloat32)) != numValues {
			return nil, fmt.Errorf("niftiToByteSlice: float32 data has %d values, but the header dimensions declare %d.", len(n.Data.DataFloat32), numValues)
		}
		if err := binary.Write(&buf, endian, n.Data.DataFloat32); err != nil {
			return nil, fmt.Errorf("niftiToByteSlice: could not write float32 voxel data: %s", err)
		}
	default:
		return nil, fmt.Errorf("niftiToByteSlice: unsupported NIfTI-1 data type code %d.", n.Header.Datatype)
	}

	return buf.Bytes(), nil
}

// WriteNifti writes a Nifti struct to a file in NIfTI-1 format (single-file, "n+1").
//
// The NIfTI-1 file is written in big endian byte order (the NIfTI-1 standard byte
// order). Gzipped files (file extension .nii.gz) are supported. Note that this
// function always writes standard-conformant NIfTI-1 files; in particular, the
// so-called "FreeSurfer hack" is never produced.
//
// Parameters:
//   - filename: path to the output file, e.g. '<subject>/mri/brain.nii' or
//     '<subject>/mri/brain.nii.gz'.
//   - n: the Nifti struct to write.
//   - isGzipped: whether to write the file gzip-compressed (NIfTI-GZ format). If
//     "auto", the file extension is used to determine whether the file should be
//     gzip-compressed. If not "auto", it has to be "yes"/"nii.gz" or "no"/"nii" to
//     force gzip-compressed or uncompressed format, respectively.
//
// Returns:
//   - error: an error if one occurred, or nil on success
func WriteNifti(filename string, n Nifti, isGzipped string) error {

	bs, err := niftiToByteSlice(n)
	if err != nil {
		return fmt.Errorf("WriteNifti: could not serialize Nifti struct to NIfTI-1 format: %s", err)
	}

	isGzipped = getIsGzippedNifti(isGzipped)
	gzipIt := getIsGzipped(filename, isGzipped)

	if Verbosity >= 1 {
		fmt.Printf("WriteNifti: Writing %d bytes of NIfTI-1 data to file '%s', gzip=%t.\n", len(bs), filename, gzipIt)
	}

	if err := writeByteSliceToFile(bs, filename, gzipIt); err != nil {
		return fmt.Errorf("WriteNifti: could not write NIfTI-1 data to file '%s': %s", filename, err)
	}

	return nil
}
