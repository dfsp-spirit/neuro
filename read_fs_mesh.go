package neurogo

// Related packages and documentation:
// https://pkg.go.dev/github.com/oschwald/maxminddb-golang#example-Reader.Lookup-Interface
// https://pkg.go.dev/encoding/binary#example-Read-Multi
// maybe https://www.jonathan-petitcolas.com/2014/09/25/parsing-binary-files-in-go.html, but it's old
//
// https://github.com/dfsp-spirit/libfs/blob/main/include/libfs.h#L2023 for the fs surface file format

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

// Verbosity is the verbosity level of the package. 0 = silent, 1 = info, 2 = debug.
var Verbosity int = 1

// Read a newline-terminated string from a bytes.Reader.
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

func _magicByte3ToInt(magic []byte) (int) {
	int1, _ := strconv.Atoi(string(magic[:])) // FIXME: incorrect
	int1 = ((int1 >> 8) & 0xffffff);
	return int1
}

// Mesh is a struct that holds a mesh, with vertices and faces.
type Mesh struct {
	Vertices []float32
	Faces []int32
}


// Compute some basic mesh statistics.
func MeshStats(mesh Mesh) (map[string]float32, error) {

	if len(mesh.Faces) == 0 {
		return nil, fmt.Errorf("MeshStats: mesh has no faces.")
	}
	if len(mesh.Vertices) == 0 {
		return nil, fmt.Errorf("MeshStats: mesh has no vertices.")
	}

	stats := map[string]float32{ "numVertices" : float32(len(mesh.Vertices) / 3),
		 "numFaces" : float32(len(mesh.Faces) / 3)}

	var max_x float32 = 0.0
	var max_y float32 = 0.0
	var max_z float32 = 0.0

	var min_x float32 = math.MaxFloat32
	var min_y float32 = math.MaxFloat32
	var min_z float32 = math.MaxFloat32

	for i := 0; i < len(mesh.Vertices); i += 3 {
		if mesh.Vertices[i] > max_x {
			max_x = mesh.Vertices[i]
		}
		if mesh.Vertices[i] < min_x {
			min_x = mesh.Vertices[i]
		}
	}
	for i := 1; i < len(mesh.Vertices); i += 3 {
		if mesh.Vertices[i] > max_y {
			max_y = mesh.Vertices[i]
		}
		if mesh.Vertices[i] < min_y {
			min_y = mesh.Vertices[i]
		}
	}
	for i := 2; i < len(mesh.Vertices); i += 3 {
		if mesh.Vertices[i] > max_z {
			max_z = mesh.Vertices[i]
		}
		if mesh.Vertices[i] < min_z {
			min_z = mesh.Vertices[i]
		}
	}
	stats["max_x"] = max_x
	stats["max_y"] = max_y
	stats["max_z"] = max_z

	fmt.Printf("MeshStats: max_x: %f, max_y: %f, max_z: %f\n", max_x, max_y, max_z)
	fmt.Printf("MeshStats: min_x: %f, min_y: %f, min_z: %f\n", min_x, min_y, min_z)
	fmt.Printf("MeshStats: numVertices: %d, numFaces: %d\n", int(stats["numVertices"]), int(stats["numFaces"]))

	stats["min_x"] = min_x
	stats["min_y"] = min_y
	stats["min_z"] = min_z

	// Compute average edge length and average face area
	var avg_edge_length float32 = 0.0
	var avg_face_area float32 = 0.0
	var num_edges int = 0
	for i := 0; i < len(mesh.Faces); i += 3 {
		// edge 1
		edge1_x := mesh.Vertices[mesh.Faces[i]*3] - mesh.Vertices[mesh.Faces[i+1]*3]
		edge1_y := mesh.Vertices[mesh.Faces[i]*3+1] - mesh.Vertices[mesh.Faces[i+1]*3+1]
		edge1_z := mesh.Vertices[mesh.Faces[i]*3+2] - mesh.Vertices[mesh.Faces[i+1]*3+2]
		edge1_length := float32(math.Sqrt(float64(edge1_x*edge1_x + edge1_y*edge1_y + edge1_z*edge1_z)))
		avg_edge_length += edge1_length
		num_edges++
		// edge 2
		edge2_x := mesh.Vertices[mesh.Faces[i+1]*3] - mesh.Vertices[mesh.Faces[i+2]*3]
		edge2_y := mesh.Vertices[mesh.Faces[i+1]*3+1] - mesh.Vertices[mesh.Faces[i+2]*3+1]
		edge2_z := mesh.Vertices[mesh.Faces[i+1]*3+2] - mesh.Vertices[mesh.Faces[i+2]*3+2]
		edge2_length := float32(math.Sqrt(float64(edge2_x*edge2_x + edge2_y*edge2_y + edge2_z*edge2_z)))
		avg_edge_length += edge2_length
		num_edges++
		// edge 3
		edge3_x := mesh.Vertices[mesh.Faces[i+2]*3] - mesh.Vertices[mesh.Faces[i]*3]
		edge3_y := mesh.Vertices[mesh.Faces[i+2]*3+1] - mesh.Vertices[mesh.Faces[i]*3+1]
		edge3_z := mesh.Vertices[mesh.Faces[i+2]*3+2] - mesh.Vertices[mesh.Faces[i]*3+2]
		edge3_length := float32(math.Sqrt(float64(edge3_x*edge3_x + edge3_y*edge3_y + edge3_z*edge3_z)))
		avg_edge_length += edge3_length
		num_edges++
		// compute face area
		s := (edge1_length + edge2_length + edge3_length) / 2.0
		face_area := float32(math.Sqrt(float64(s * (s - edge1_length) * (s - edge2_length) * (s - edge3_length))))
		avg_face_area += face_area
	}
	stats["num_edges"] = float32(num_edges)
	stats["avg_edge_length"] = avg_edge_length / float32(num_edges)
	stats["avg_face_area"] = avg_face_area / float32(len(mesh.Faces) / 3)
	return stats, nil
}

// ReadFreesurferMesh reads a FreeSurfer mesh file and returns a Mesh struct.
func ReadFreesurferMesh(filepath string) (Mesh, error) {
	//b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40, 0xff, 0x01, 0x02, 0x03, 0xbe, 0xef}
	//r := bytes.NewReader(b)

	endian := binary.BigEndian
	surface := Mesh{}

	if _, err := os.Stat(filepath); err != nil {
		fmt.Printf("Could not stat file '%s'\n.", filepath)
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
		//Mine [3]byte
	}

	hdr1 := header_part1{}


	if err := binary.Read(r, endian, &hdr1); err != nil {
		fmt.Println("binary.Read failed on first part of fs surface header:", err)
		return surface, err
	}


	if Verbosity > 0 {
		fmt.Println("hdr1.Magic_b1:", hdr1.MagicB1)
		fmt.Println("hdr1.Magic_b2:", hdr1.MagicB2)
		fmt.Println("hdr1.Magic_b3:", hdr1.MagicB3)
	}

	magic := make([]byte, 3)
	magic[0] = hdr1.MagicB1
	magic[1] = hdr1.MagicB2
	magic[2] = hdr1.MagicB3
	int1 := _magicByte3ToInt(magic)
	fmt.Printf("Header magic bytes %d %d %d gives: %d.\n", hdr1.MagicB1, hdr1.MagicB2, hdr1.MagicB3, int1)

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