package material

import (
	"goraytracer/geometry"
	"goraytracer/ray"
	"goraytracer/vec3"
	"math/rand"
)

type Attenuation = vec3.Vec3
type ScatterRay = ray.Ray

type Material interface {
	Scatter(r ray.Ray, hitRecord geometry.HitRecord, random *rand.Rand) (Attenuation, *ScatterRay)
}

type Lambertian struct {
	Albedo vec3.Vec3
}

type Metal struct {
	Albedo vec3.Vec3
	Fuzz   vec3.Vec3
}

func (material Lambertian) Scatter(r ray.Ray, hitRecord geometry.HitRecord, random *rand.Rand) (Attenuation, *ScatterRay) {
	scatterDir := vec3.Add(hitRecord.Normal, vec3.RandomInUnitSphere(random).Normalized())
	if scatterDir.NearZero() {
		scatterDir = hitRecord.Normal
	}
	scatterRay := ray.New(hitRecord.Point, scatterDir)

	color := vec3.Vec3{
		X: hitRecord.U,
		Y: 0.0,
		Z: hitRecord.V,
	}
	return color, scatterRay
}

func (material Metal) Scatter(r ray.Ray, hitRecord geometry.HitRecord, random *rand.Rand) (Attenuation, *ScatterRay) {
	return vec3.Vec3{}, ray.New(vec3.Vec3{}, vec3.Vec3{})
}
