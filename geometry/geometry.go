package geometry

import (
	"goraytracer/ray"
	"goraytracer/vec3"
)

type HitRecord struct {
	Hit      bool
	Distance float64 // distance from ray origin
	Point    vec3.Vec3
	Normal   vec3.Vec3
}

type Geometry interface {
	// TODO: rename this to IntersectsRay
	Hit(r *ray.Ray, minDistance float64, maxDistance float64) HitRecord

	// Returns a reference to self or sub geometries that intersect with aabb
	AABBIntersections(aabb AABB) []*Geometry

	IntersectsAABB(aabb AABB) bool
}

// intersections:
// ray - geometry
// ray - aabb
// aabb - geometry
