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
	Hit(r *ray.Ray, minDistance float64, maxDistance float64) HitRecord
}
