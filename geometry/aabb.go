package geometry

import (
	"goraytracer/ray"
	"goraytracer/vec3"
)

type AABB struct {
	Min vec3.Vec3
	Max vec3.Vec3
}

func (aabb AABB) IntersectsRay(ray *ray.Ray) bool {
	return false
}

// could make an interface Intersects
