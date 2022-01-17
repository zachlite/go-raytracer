package geometry

import (
	"goraytracer/ray"
	"goraytracer/vec3"
	"math"
)

type AABB struct {
	Min vec3.Vec3
	Max vec3.Vec3
}

// highly optimized ray aabb intersection test. no division.
// https://tavianator.com/cgit/dimension.git/tree/libdimension/bvh/bvh.c#n196
func (aabb AABB) IntersectsRay(ray *ray.Ray) bool {

	tx1 := (aabb.Min.X - ray.Origin.X) * ray.InverseDir.X
	tx2 := (aabb.Max.X - ray.Origin.X) * ray.InverseDir.X

	tmin := math.Min(tx1, tx2)
	tmax := math.Max(tx1, tx2)

	ty1 := (aabb.Min.Y - ray.Origin.Y) * ray.InverseDir.Y
	ty2 := (aabb.Max.Y - ray.Origin.Y) * ray.InverseDir.Y

	tmin = math.Max(tmin, math.Min(ty1, ty2))
	tmax = math.Min(tmax, math.Max(ty1, ty2))

	tz1 := (aabb.Min.Z - ray.Origin.Z) * ray.InverseDir.Z
	tz2 := (aabb.Max.Z - ray.Origin.Z) * ray.InverseDir.Z

	tmin = math.Max(tmin, math.Min(tz1, tz2))
	tmax = math.Min(tmax, math.Max(tz1, tz2))

	return tmax >= math.Max(0.0, tmin)
}
