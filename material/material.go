package material

import (
	"goraytracer/ray"
	"goraytracer/vec3"
	"math"
	"math/rand"
)

type Material interface {
	// ray, hit point, hit normal -> attenuation, scatter
	Scatter(r ray.Ray, hitPoint vec3.Vec3, hitNormal vec3.Vec3) (vec3.Vec3, ray.Ray)
}

type Lambertian struct {
	Albedo vec3.Vec3
}

type Metal struct {
	Albedo vec3.Vec3
	Fuzz   vec3.Vec3
}

func randomScatter() vec3.Vec3 {
	scale := 1.0
	r := rand.Float64()
	z := rand.Float64()
	zScale := math.Sqrt(1.0-z*z) * scale
	return vec3.Vec3{
		X: math.Cos(r) * zScale,
		Y: math.Sin(r) * zScale,
		Z: z * scale,
	}
}

func nearZero(v vec3.Vec3) bool {
	limit := .00000001
	return math.Abs(v.X) < limit && math.Abs(v.Y) < limit && math.Abs(v.Z) < limit
}

func (material Lambertian) Scatter(r ray.Ray, hitPoint vec3.Vec3, hitNormal vec3.Vec3) (vec3.Vec3, ray.Ray) {

	scatterDir := vec3.Add(hitNormal, randomScatter())

	if nearZero(scatterDir) {
		scatterDir = hitNormal
	}

	scatterRay := ray.Ray{Origin: hitPoint, Direction: scatterDir}
	return material.Albedo, scatterRay
}
