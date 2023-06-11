package neuro

// Related software: libfs for C++, see:
// https://github.com/dfsp-spirit/libfs/blob/main/include/libfs.h#L2023 for the fs surface file format

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// Read a newline-terminated string from a bytes.Reader.
//
// Parameters:
//  - r: a bytes.Reader
//  - endian: the byte order, e.g. binary.BigEndian
//  - do_strip_newline: if true, strip the newline character from the end of the string
//
// Returns:
//  - string: the string
//  - error: an error if one occurred
func readNewlineTerminatedString(r *bytes.Reader, endian binary.ByteOrder, do_strip_newline bool) (string, error) {

	//endian = binary.BigEndian // TODO: make this a parameter

	//return "not implemented yet", nil
	var line string = ""
	char := make([]byte, 1)
	var char_num int = 0
	for string(char[:]) != "\n" {
		if err := binary.Read(r, endian, &char); err != nil {
			fmt.Printf("binary.Read failed on character %d of newline-terminated string: %s\n", char_num, err)
			return "", err
		} else {
			if !(do_strip_newline && string(char[:]) == "\n") {
				line = line + string(char[:])
			}
		}
		char_num++
	}
	return line, nil
}


// ReadFsSurface reads a FreeSurfer surface file and returns a Mesh struct.
//
// A surface file is a binary file containing the reconstructed surface of a brain hemisphere.
//
// Parameters:
//  - filepath: path to the FreeSurfer mesh file, e.g. '<subject>/surf/lh.white'
//
// Returns:
//  - Mesh: a Mesh struct containing the mesh data
//  - error: an error if one occurred
func ReadFsSurface(filepath string) (Mesh, error) {

	endian := binary.BigEndian
	surface := Mesh{}

	if _, err := os.Stat(filepath); err != nil {
		err := fmt.Errorf("Could not stat surface file '%f': %s\n", filepath, err)
		return surface, err
	}

	file, err := os.Open(filepath)
	if err != nil {
   		panic(err)
	}
	defer file.Close()

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
	   fmt.Println(err)
	   return surface, err
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
	   fmt.Println(err)
	   return surface, err
	}

	// Read the byte slice
	r := bytes.NewReader(bs)

	type header_part1 struct {
		MagicB1 uint8
		MagicB2 uint8
		MagicB3 uint8
	}

	hdr1 := header_part1{}


	if err := binary.Read(r, endian, &hdr1); err != nil {
		fmt.Println("binary.Read failed on first part of fs surface header:", err)
		return surface, err
	}


	if ! (hdr1.MagicB1 == 255 && hdr1.MagicB2 == 255 && hdr1.MagicB3 == 254) {
		fmt.Println("Error: surface magic bytes are not 255 255 254, this is not a FreeSurfer surface file. Provide a recon-all output file like '<subject>/surf/lh.white'.")
		return surface, err
	}


	if Verbosity > 0 {
		fmt.Printf("Surface header magic bytes: %d %d %d.\n", hdr1.MagicB1, hdr1.MagicB2, hdr1.MagicB3)
	}

	createdLine, err := readNewlineTerminatedString(r, endian, true);
    commentLine, err := readNewlineTerminatedString(r, endian, true);

	if Verbosity > 0 {
		fmt.Printf("createdLine: '%s'\n", createdLine)
		fmt.Printf("commentLine: '%s'\n", commentLine)
	}

	type header_part2 struct {
		NumVerts int32
		NumFaces int32
	}

	hdr2 := header_part2{}

	if err := binary.Read(r, endian, &hdr2); err != nil {
		fmt.Println("binary.Read failed on second part of fs surface header:", err)
		return surface, err
	}

	if Verbosity > 0 {
		fmt.Println("NumVerts:", hdr2.NumVerts)
		fmt.Println("NumFaces:", hdr2.NumFaces)
	}

	// read mesh data
	surface.Vertices = make([]float32, hdr2.NumVerts * 3) // x,y,z coordinates for each vertex
	surface.Faces = make([]int32, hdr2.NumFaces * 3)  // vertex 1, 2, 3 for each face

	// read vertices
	if err := binary.Read(r, endian, &surface.Vertices); err != nil {
		fmt.Println("binary.Read failed on mesh vertices array:", err)
		return surface, err
	}

	// read faces
	if err := binary.Read(r, endian, &surface.Faces); err != nil {
		fmt.Println("binary.Read failed on mesh faces array:", err)
		return surface, err
	}

	if Verbosity >= 2 {
		var numToPrint int = 5
		if hdr2.NumVerts >= int32(numToPrint) {
			// print first 5 vertices
			for i := 0; i < numToPrint; i++ {
				for j := 0; j < 3; j++ {
					fmt.Println("surface.Vertices[", i, "][", j, "]:", surface.Vertices[i*3+j])
				}
			}
		}
		// print first 5 faces
		if hdr2.NumFaces >= int32(numToPrint) {
			for i := 0; i < numToPrint; i++ {
				for j := 0; j < 3; j++ {
					fmt.Println("surface.Faces[", i, "][", j, "]:", surface.Faces[i*3+j])
				}
			}
		}
	}

	return surface, nil
}