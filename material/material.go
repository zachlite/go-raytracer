package material

import (
	"goraytracer/ray"
	"goraytracer/vec3"
	"math/rand"
)

type Attenuation = vec3.Vec3
type ScatterRay = ray.Ray

type Material interface {
	Scatter(r ray.Ray, hitPoint vec3.Vec3, hitNormal vec3.Vec3, random *rand.Rand) (Attenuation, ScatterRay)
}

type Lambertian struct {
	Albedo vec3.Vec3
}

type Metal struct {
	Albedo vec3.Vec3
	Fuzz   vec3.Vec3
}

func (material Lambertian) Scatter(r ray.Ray, hitPoint vec3.Vec3, hitNormal vec3.Vec3, random *rand.Rand) (Attenuation, ScatterRay) {
	scatterDir := vec3.Add(hitNormal, vec3.RandomInUnitSphere(random).Normalized())
	if scatterDir.NearZero() {
		scatterDir = hitNormal
	}
	scatterRay := ray.New(hitPoint, scatterDir)
	return material.Albedo, scatterRay
}

func (material Metal) Scatter(r ray.Ray, hitPoint vec3.Vec3, hitNormal vec3.Vec3, random *rand.Rand) (Attenuation, ScatterRay) {
	return vec3.Vec3{}, ray.Ray{}
}
