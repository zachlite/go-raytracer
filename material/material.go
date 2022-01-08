package material

import (
	"goraytracer/ray"
	"goraytracer/vec3"
	"math"
	"math/rand"
)

type Material interface {
	// ray, hit point, hit normal -> attenuation, scatter
	Scatter(r ray.Ray, hitPoint vec3.Vec3, hitNormal vec3.Vec3, random *rand.Rand) (vec3.Vec3, ray.Ray)
}

type Lambertian struct {
	Albedo vec3.Vec3
}

type Metal struct {
	Albedo vec3.Vec3
	Fuzz   vec3.Vec3
}

func nearZero(v vec3.Vec3) bool {
	limit := .00000001
	return math.Abs(v.X) < limit && math.Abs(v.Y) < limit && math.Abs(v.Z) < limit
}

func (material Lambertian) Scatter(r ray.Ray, hitPoint vec3.Vec3, hitNormal vec3.Vec3, random *rand.Rand) (vec3.Vec3, ray.Ray) {

	randomDirection := vec3.RandomInUnitSphere(random)
	randomDirection = randomDirection.Normalize()
	scatterDir := vec3.Add(hitNormal, randomDirection)

	if nearZero(scatterDir) {
		scatterDir = hitNormal
	}

	scatterRay := ray.Ray{Origin: hitPoint, Direction: scatterDir}
	return material.Albedo, scatterRay
}
