package neuro

import (
	"fmt"
	"math"
	"strings"
)

// Mesh is a struct that holds a triangular mesh, with vertices and faces. Faces are stored in vertex index representation.
//
// Fields:
//   - Vertices : the vertices of the mesh, as a slice of float32 values. The vertices are stored as a flat array of 3D coordinates, i.e. [x1, y1, z1, x2, y2, z2, ...]
//   - Faces    : the faces (a.k.a polygons or triangles) of the mesh, as a slice of int32 values. The faces are stored as a flat array of vertex indices, i.e. [v1, v2, v3, v1, v2, v3, ...]
type Mesh struct {
	Vertices []float32
	Faces    []int32
}

// Compute some basic mesh statistics.
//
// Parameters:
//   - mesh : the mesh to compute statistics for
//
// Returns:
//   - map[string]float32 : a map of statistics, with keys: 'numVertices' (number of vertices, interpret as int), 'numFaces' (number of faces, interpret as int), 'maxX', 'maxY', 'maxZ', 'minX', 'minY', 'minZ', 'meanX', 'meanY', 'meanZ', 'numEdges', 'avgEdgeLength', 'avgFaceArea', 'totalArea'.
func MeshStats(mesh Mesh) (map[string]float32, error) {

	if len(mesh.Faces) < 3 {
		return nil, fmt.Errorf("MeshStats: mesh has no faces.")
	}
	if len(mesh.Vertices) < 3 {
		return nil, fmt.Errorf("MeshStats: mesh has no vertices.")
	}

	stats := map[string]float32{"numVertices": float32(NumVertices(mesh)),
		"numFaces": float32(NumFaces(mesh))}

	var max_x float32 = 0.0
	var max_y float32 = 0.0
	var max_z float32 = 0.0

	var min_x float32 = math.MaxFloat32
	var min_y float32 = math.MaxFloat32
	var min_z float32 = math.MaxFloat32

	var mean_x float32 = 0.0
	var mean_y float32 = 0.0
	var mean_z float32 = 0.0

	for i := 0; i < len(mesh.Vertices); i += 3 {
		mean_x += mesh.Vertices[i]
		if mesh.Vertices[i] > max_x {
			max_x = mesh.Vertices[i]
		}
		if mesh.Vertices[i] < min_x {
			min_x = mesh.Vertices[i]
		}
	}
	for i := 1; i < len(mesh.Vertices); i += 3 {
		mean_y += mesh.Vertices[i]
		if mesh.Vertices[i] > max_y {
			max_y = mesh.Vertices[i]
		}
		if mesh.Vertices[i] < min_y {
			min_y = mesh.Vertices[i]
		}
	}
	for i := 2; i < len(mesh.Vertices); i += 3 {
		mean_z += mesh.Vertices[i]
		if mesh.Vertices[i] > max_z {
			max_z = mesh.Vertices[i]
		}
		if mesh.Vertices[i] < min_z {
			min_z = mesh.Vertices[i]
		}
	}
	stats["maxX"] = max_x
	stats["maxY"] = max_y
	stats["maxZ"] = max_z

	if Verbosity >= 2 {
		fmt.Printf("MeshStats: max_x: %f, max_y: %f, max_z: %f\n", max_x, max_y, max_z)
		fmt.Printf("MeshStats: min_x: %f, min_y: %f, min_z: %f\n", min_x, min_y, min_z)
		fmt.Printf("MeshStats: numVertices: %d, numFaces: %d\n", int(stats["numVertices"]), int(stats["numFaces"]))
	}

	stats["minX"] = min_x
	stats["minY"] = min_y
	stats["minZ"] = min_z

	stats["meanX"] = mean_x / float32(len(mesh.Vertices)/3)
	stats["meanY"] = mean_y / float32(len(mesh.Vertices)/3)
	stats["meanZ"] = mean_z / float32(len(mesh.Vertices)/3)

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
	stats["numEdges"] = float32(num_edges)
	stats["avgEdgeLength"] = avg_edge_length / float32(num_edges)
	stats["avgFaceArea"] = avg_face_area / float32(len(mesh.Faces)/3)
	stats["totalArea"] = avg_face_area
	return stats, nil
}

// Convert a mesh to PLY format.
//
// Parameters:
//   - mesh : the mesh to convert
//
// Returns:
//   - string : the mesh string representation in PLY format
//   - error  : the error if one occured, or nil otherwise
func ToPlyFormat(mesh Mesh) (string, error) {

	if Verbosity >= 2 {
		fmt.Printf("Generating PLY representation for mesh with %d vertices and %d faces.\n", len(mesh.Vertices)/3, len(mesh.Faces)/3)
	}
	var ply strings.Builder
	ply.WriteString("ply\n")
	ply.WriteString("format ascii 1.0\n")
	ply.WriteString("comment neurogo\n")
	ply.WriteString(fmt.Sprintf("element vertex %d\n", len(mesh.Vertices)/3))
	ply.WriteString("property float x\n")
	ply.WriteString("property float y\n")
	ply.WriteString("property float z\n")
	ply.WriteString(fmt.Sprintf("element face %d\n", len(mesh.Faces)/3))
	ply.WriteString("property list uchar int vertex_indices\n")
	ply.WriteString("end_header\n")

	for i := 0; i < len(mesh.Vertices); i += 3 {
		ply.WriteString(fmt.Sprintf("%f %f %f\n", mesh.Vertices[i], mesh.Vertices[i+1], mesh.Vertices[i+2]))
	}

	for i := 0; i < len(mesh.Faces); i += 3 {
		ply.WriteString(fmt.Sprintf("3 %d %d %d\n", mesh.Faces[i], mesh.Faces[i+1], mesh.Faces[i+2]))
	}

	return ply.String(), nil
}

// Convert a mesh to OBJ format.
//
// Parameters:
//   - mesh : the mesh to convert
//
// Returns:
//   - string : the mesh string representation in OBJ format
//   - error  : the error if one occured, or nil otherwise
func ToObjFormat(mesh Mesh) (string, error) {

	if Verbosity >= 2 {
		fmt.Printf("Generating OBJ representation for mesh with %d vertices and %d faces.\n", len(mesh.Vertices)/3, len(mesh.Faces)/3)
	}

	var obj strings.Builder
	obj.WriteString("# neurogo\n")
	for i := 0; i < len(mesh.Vertices); i += 3 {
		obj.WriteString(fmt.Sprintf("v %f %f %f\n", mesh.Vertices[i], mesh.Vertices[i+1], mesh.Vertices[i+2]))
	}

	for i := 0; i < len(mesh.Faces); i += 3 {
		obj.WriteString(fmt.Sprintf("f %d %d %d\n", mesh.Faces[i]+1, mesh.Faces[i+1]+1, mesh.Faces[i+2]+1))
	}

	return obj.String(), nil
}

// Convert a mesh to STL format.
//
// Parameters:
//   - mesh : the mesh to convert
//
// Returns:
//   - string : the mesh string representation in STL format
//   - error  : the error if one occured, or nil otherwise
func ToStlFormat(mesh Mesh) (string, error) {

	if Verbosity >= 2 {
		fmt.Printf("Generating STL representation for mesh with %d vertices and %d faces.\n", len(mesh.Vertices)/3, len(mesh.Faces)/3)
	}

	var stl strings.Builder
	stl.WriteString("solid neurogo\n")
	for i := 0; i < len(mesh.Faces); i += 3 {
		// compute face normal
		// edge 1
		edge1_x := mesh.Vertices[mesh.Faces[i]*3] - mesh.Vertices[mesh.Faces[i+1]*3]
		edge1_y := mesh.Vertices[mesh.Faces[i]*3+1] - mesh.Vertices[mesh.Faces[i+1]*3+1]
		edge1_z := mesh.Vertices[mesh.Faces[i]*3+2] - mesh.Vertices[mesh.Faces[i+1]*3+2]
		// edge 2
		edge2_x := mesh.Vertices[mesh.Faces[i+1]*3] - mesh.Vertices[mesh.Faces[i+2]*3]
		edge2_y := mesh.Vertices[mesh.Faces[i+1]*3+1] - mesh.Vertices[mesh.Faces[i+2]*3+1]
		edge2_z := mesh.Vertices[mesh.Faces[i+1]*3+2] - mesh.Vertices[mesh.Faces[i+2]*3+2]
		// cross product
		norm_x := edge1_y*edge2_z - edge1_z*edge2_y
		norm_y := edge1_z*edge2_x - edge1_x*edge2_z
		norm_z := edge1_x*edge2_y - edge1_y*edge2_x
		// normalize
		norm_length := float32(math.Sqrt(float64(norm_x*norm_x + norm_y*norm_y + norm_z*norm_z)))
		norm_x /= norm_length
		norm_y /= norm_length
		norm_z /= norm_length
		// write face normal
		stl.WriteString(fmt.Sprintf("facet normal %f %f %f\n", norm_x, norm_y, norm_z))
		stl.WriteString("outer loop\n")
		// write face vertices
		stl.WriteString(fmt.Sprintf("vertex %f %f %f\n", mesh.Vertices[mesh.Faces[i]*3], mesh.Vertices[mesh.Faces[i]*3+1], mesh.Vertices[mesh.Faces[i]*3+2]))
		stl.WriteString(fmt.Sprintf("vertex %f %f %f	\n", mesh.Vertices[mesh.Faces[i+1]*3], mesh.Vertices[mesh.Faces[i+1]*3+1], mesh.Vertices[mesh.Faces[i+1]*3+2]))
		stl.WriteString(fmt.Sprintf("vertex %f %f %f\n", mesh.Vertices[mesh.Faces[i+2]*3], mesh.Vertices[mesh.Faces[i+2]*3+1], mesh.Vertices[mesh.Faces[i+2]*3+2]))
		stl.WriteString("endloop\n")
		stl.WriteString("endfacet\n")
	}
	stl.WriteString("endsolid neurogo\n")
	return stl.String(), nil
}

// Export exports a mesh to a file in the specified mesh file format.
//
// Parameters:
//   - mesh     : the mesh to export
//   - filepath : the filepath to export the mesh to
//   - format   : the mesh file format to use, one of 'obj' (for Wavefront Object Format), 'ply' (for Stanford PLY format), 'stl' (for StereoLithography format)
//
// Returns
//   - string : the mesh string representation in the requested format
//   - error  : the error if one occured, or nil otherwise
func Export(mesh Mesh, filepath string, format string) (string, error) {
	var mesh_rep string
	var err error
	if format == "stl" || format == "STL" {
		mesh_rep, err = ToStlFormat(mesh)
	} else if format == "obj" || format == "OBJ" {
		mesh_rep, err = ToObjFormat(mesh)
	} else if format == "ply" || format == "PLY" {
		mesh_rep, err = ToPlyFormat(mesh)
	} else {
		err = fmt.Errorf("Invalid mesh export format specified, use one of 'obj', 'ply', 'stl'.")
		return mesh_rep, err
	}
	if err != nil {
		return mesh_rep, err
	}
	err = strToTextFile(mesh_rep, filepath)
	return mesh_rep, err
}

// NumVertices computes the number of vertices of a triangular mesh.
//
// Parameters:
//   - mesh : the mesh to compute the number of vertices for
//
// Returns:
//   - int : the number of vertices
func NumVertices(mesh Mesh) int {
	return len(mesh.Vertices) / 3
}

// NumFaces computes the number of faces (aka polygons, or triangles) of a triangular mesh.
//
// Parameters:
//   - mesh : the mesh to compute the number of faces for
//
// Returns:
//   - int : the number of faces
func NumFaces(mesh Mesh) int {
	return len(mesh.Faces) / 3
}

// GenerateCube creates and returns a Mesh representing a cube.
//
// This is mainly used in the examples and documentation.
//
// Returns:
//   - Mesh : the cube mesh
func GenerateCube() Mesh {

	var mesh Mesh

	mesh.Vertices = []float32{1.0, 1.0, 1.0,
		1.0, 1.0, -1.0,
		1.0, -1.0, 1.0,
		1.0, -1.0, -1.0,
		-1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0,
		-1.0, -1.0, 1.0,
		-1.0, -1.0, -1.0}

	mesh.Faces = []int32{0, 2, 3,
		3, 1, 0,
		4, 6, 7,
		7, 5, 4,
		0, 4, 5,
		5, 1, 0,
		2, 6, 7,
		7, 3, 2,
		0, 4, 6,
		6, 2, 0,
		1, 5, 7,
		7, 3, 1}
	return mesh
}

// GenerateSphere creates and returns a Mesh representing a sphere.
//
// This is mainly used in the examples and documentation.
//
// Parameters:
//   - radius : the radius of the sphere
//   - slices : the number of slices (horizontal divisions)
//   - stacks : the number of stacks (vertical divisions)
//
// Returns:
//   - Mesh : the sphere mesh
func GenerateSphere(radius float32, slices int, stacks int) Mesh {

	var mesh Mesh

	// generate vertices
	for i := 0; i <= stacks; i++ {
		phi := float32(i) * math.Pi / float32(stacks)
		for j := 0; j <= slices; j++ {
			theta := float32(j) * 2 * math.Pi / float32(slices)
			x := radius * float32(math.Cos(float64(theta))) * float32(math.Sin(float64(phi)))
			y := radius * float32(math.Sin(float64(theta))) * float32(math.Sin(float64(phi)))
			z := radius * float32(math.Cos(float64(phi)))
			mesh.Vertices = append(mesh.Vertices, x, y, z)
		}
	}

	// generate faces
	for i := 0; i < stacks; i++ {
		for j := 0; j < slices; j++ {
			mesh.Faces = append(mesh.Faces, int32((i+1)*(slices+1)+j))
			mesh.Faces = append(mesh.Faces, int32(i*(slices+1)+j))
			mesh.Faces = append(mesh.Faces, int32(i*(slices+1)+j+1))
			mesh.Faces = append(mesh.Faces, int32((i+1)*(slices+1)+j))
			mesh.Faces = append(mesh.Faces, int32(i*(slices+1)+j+1))
			mesh.Faces = append(mesh.Faces, int32((i+1)*(slices+1)+j+1))
		}
	}

	return mesh
}
