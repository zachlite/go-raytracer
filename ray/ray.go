package ray

import (
	"goraytracer/vec3"
	"math"
)

type Ray struct {
	Origin     vec3.Vec3
	Direction  vec3.Vec3
	InverseDir vec3.Vec3
}

func New(origin vec3.Vec3, direction vec3.Vec3) Ray {

	inverseDir := vec3.Vec3{}

	if direction.X != 0 {
		inverseDir.X = 1.0 / direction.X
	} else {
		inverseDir.X = math.Inf(1)
	}

	if direction.Y != 0 {
		inverseDir.Y = 1.0 / direction.Y
	} else {
		inverseDir.Y = math.Inf(1)
	}

	if direction.Z != 0 {
		inverseDir.Z = 1.0 / direction.Z
	} else {
		inverseDir.Z = math.Inf(1)
	}

	return Ray{Origin: origin, Direction: direction, InverseDir: inverseDir}
}

func (r *Ray) At(t float64) vec3.Vec3 {
	return vec3.Add(r.Origin, vec3.MultiplyScalar(r.Direction, t))
}
