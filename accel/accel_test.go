package accel_test

import (
	"goraytracer/accel"
	"goraytracer/geometry"
	"goraytracer/vec3"
	"testing"
)

func assertEqual(t *testing.T, got interface{}, expected interface{}, message string) {
	if got != expected {
		t.Errorf("%s got: %v expected: %v", message, got, expected)
	}
}

func TestOctTree(t *testing.T) {

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
