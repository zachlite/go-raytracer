package accel_test

import (
	"encoding/json"
	"goraytracer/accel"
	"goraytracer/geometry"
	"goraytracer/material"
	"goraytracer/mesh"
	"goraytracer/ray"
	samplemodels "goraytracer/sample_models"
	"goraytracer/vec3"
	"log"
	"os"
	"testing"
)

func assertEqual(t *testing.T, got interface{}, expected interface{}, message string) {
	if got != expected {
		t.Errorf("%s got: %v expected: %v", message, got, expected)
	}
}

func TestNonIntersectedNodesReturnZeroCandidates(t *testing.T) {
	// a non-intersected node returns nil, with our without children
	aabb := geometry.AABB{
		Min: vec3.Vec3{X: 10, Y: 10, Z: 10},
		Max: vec3.Vec3{X: 11, Y: 11, Z: 11},
	}

	cameraRay := ray.New(vec3.Vec3{}, vec3.Vec3{Z: 1})

	nodeWithoutChildren := accel.OctTreeNode{
		Aabb:                aabb,
		Children:            nil,
		IntersectCandidates: nil,
		Depth:               0,
	}

	candidates := make([]accel.IntersectCandidate, 1)

	candidates[0] = accel.IntersectCandidate{
		Geometry: geometry.Triangle{
			P1:     vec3.Vec3{},
			P2:     vec3.Vec3{},
			P3:     vec3.Vec3{},
			Normal: vec3.Vec3{},
		},
		Material: material.Lambertian{Albedo: vec3.Vec3{}},
	}
	nodeWithChildren := accel.OctTreeNode{
		Aabb:                accel.SplitAABB(aabb)[0],
		Children:            nil,
		IntersectCandidates: candidates,
		Depth:               1,
	}
	assertEqual(t, len(nodeWithChildren.Search(&cameraRay)), 0, "A non-intersected node without children returns nil")
	assertEqual(t, len(nodeWithoutChildren.Search(&cameraRay)), 0, "A non-intersected node without children returns nil")
}

// TODO
func TestIntersectedNodeReturnsCandidates(t *testing.T) {

	meshes := make([]mesh.Mesh, 1)
	meshes[0] = mesh.Mesh{
		Geometry: samplemodels.LoadBunny(),
		Material: material.Lambertian{
			Albedo: vec3.Vec3{
				X: 0,
				Y: 1,
				Z: 0,
			},
		},
	}

	f, err := os.Create("tree.json")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	tree := accel.BuildOctTree(meshes)

	_ = json.NewEncoder(f).Encode(
		tree,
	)

}

// TODO
func TestNoDuplicateGeometriesAreReturned(t *testing.T) {

}

func TestOctTreeSearch(t *testing.T) {
	// what makes this correct?
	// no duplicate geometries are returned - even if the same geometry exists in multiple intersected leaf nodes.
	// an intersected octree node with no children return nil
	// a non-intersected node returns nil, with our without children
	// an intersected node with attached geometry returns that geometry.

	//	meshes := make([]mesh.Mesh, 1)
	//	meshes[0] = mesh.Mesh{
	//		Geometry: geometry.Sphere{
	//			Center: vec3.Vec3{0, 0, 0},
	//			Radius: .1,
	//		},
	//		Material: material.Lambertian{Albedo: vec3.Vec3{
	//			X: .5,
	//			Y: .5,
	//			Z: .5,
	//		}},
	//	}
	//
	//	//tree := accel.BuildOctTree(meshes)
	//
}

func TestSplitAABB(t *testing.T) {
	parent := geometry.AABB{Min: vec3.Vec3{
		X: -1,
		Y: -1,
		Z: -1,
	}, Max: vec3.Vec3{
		X: 1,
		Y: 1,
		Z: 1,
	}}

	children := accel.SplitAABB(parent)
	assertEqual(t, len(children), 8, "Correct number of sub AABBs.")

	// child 0
	assertEqual(t, children[0].Min, vec3.Vec3{X: -1, Y: -1, Z: -1}, "child 0 min")
	assertEqual(t, children[0].Max, vec3.Vec3{X: 0, Y: 0, Z: 0}, "child 0 max")

	// child 1
	assertEqual(t, children[1].Min, vec3.Vec3{X: -1, Y: -1, Z: 0}, "child 1 min")
	assertEqual(t, children[1].Max, vec3.Vec3{X: 0, Y: 0, Z: 1}, "child 1 max")

	// child 2
	assertEqual(t, children[2].Min, vec3.Vec3{X: 0, Y: -1, Z: 0}, "child 2 min")
	assertEqual(t, children[2].Max, vec3.Vec3{X: 1, Y: 0, Z: 1}, "child 2 max")

	// child 3
	assertEqual(t, children[3].Min, vec3.Vec3{X: 0, Y: -1, Z: -1}, "child 3 min")
	assertEqual(t, children[3].Max, vec3.Vec3{X: 1, Y: 0, Z: 0}, "child 3 max")

	// child 4
	assertEqual(t, children[4].Min, vec3.Vec3{X: -1, Y: 0, Z: -1}, "child 4 min")
	assertEqual(t, children[4].Max, vec3.Vec3{X: 0, Y: 1, Z: 0}, "child 4 max")

	// child 5
	assertEqual(t, children[5].Min, vec3.Vec3{X: -1, Y: 0, Z: 0}, "child 5 min")
	assertEqual(t, children[5].Max, vec3.Vec3{X: 0, Y: 1, Z: 1}, "child 5 max")

	// child 6
	assertEqual(t, children[6].Min, vec3.Vec3{X: 0, Y: 0, Z: 0}, "child 6 min")
	assertEqual(t, children[6].Max, vec3.Vec3{X: 1, Y: 1, Z: 1}, "child 6 max")

	// child 7
	assertEqual(t, children[7].Min, vec3.Vec3{X: 0, Y: 0, Z: -1}, "child 7 min")
	assertEqual(t, children[7].Max, vec3.Vec3{X: 1, Y: 1, Z: 0}, "child 7 max")
}
