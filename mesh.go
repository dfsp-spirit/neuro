package neurogo

import (
	"fmt"
	"math"
)

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
	stats["total_area"] = avg_face_area
	return stats, nil
}
