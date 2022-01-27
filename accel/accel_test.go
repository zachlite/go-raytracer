package accel_test

import (
	"goraytracer/accel"
	"goraytracer/geometry"
	"goraytracer/material"
	"goraytracer/mesh"
	"goraytracer/ray"
	"goraytracer/vec3"
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
		Material: &material.Lambertian{},
	}
	nodeWithChildren := accel.OctTreeNode{
		Aabb:                accel.SplitAABB(aabb)[0],
		Children:            nil,
		IntersectCandidates: candidates,
		Depth:               1,
	}

	assertEqual(t, len(nodeWithChildren.Search(cameraRay)), 0, "A non-intersected node without children returns nil")
	assertEqual(t, len(nodeWithoutChildren.Search(cameraRay)), 0, "A non-intersected node without children returns nil")
}

func countAllCandidates(node accel.OctTreeNode) int {
	sum := 0

	if node.Depth == accel.MaxDepth {
		sum += len(node.IntersectCandidates)
	} else {
		for _, child := range node.Children {
			sum += countAllCandidates(child)
		}
	}

	return sum
}

func TestIntersectedNodeReturnsUniqueCandidates(t *testing.T) {
	candidateX := geometry.Sphere{
		Id:     0,
		Center: vec3.Vec3{X: .2, Y: .5, Z: .3},
		Radius: .1,
	}

	meshes := []mesh.Mesh{
		{
			Geometry: candidateX,
			Material: &material.Lambertian{},
		},
	}

	tree := accel.BuildOctTree(meshes)

	cameraRay := ray.New(
		vec3.Vec3{X: .2, Z: .3},
		vec3.Vec3{X: 0, Y: 1, Z: 0}.Normalized())

	candidates := tree.Search(cameraRay)

	assertEqual(t, len(candidates), 1, "Intersected nodes return unique candidates")
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
