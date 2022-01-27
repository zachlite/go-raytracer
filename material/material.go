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
	Scatter(hitRecord geometry.HitRecord, random *rand.Rand) (Attenuation, *ScatterRay)
}

type Texture struct {
}

// Until PBR model is implemented, assume materials are Lambertian.
type MaterialProps struct {
	Albedo           vec3.Vec3
	BaseColorTexture *Texture
	EmittanceColor   vec3.Vec3
}

type Lambertian struct {
	Properties MaterialProps
}

// Attenuation doesn't seem so appropriate now that materials can emit light.
func (material *Lambertian) Scatter(hitRecord geometry.HitRecord, random *rand.Rand) (Attenuation, *ScatterRay) {
	black := vec3.Vec3{}
	if material.Properties.EmittanceColor != black {
		// emit light and do not scatter.
		// umm... how can light and baseColorTexture be blended together?
		// does the blended texture effect the light being emitted? It should.

		//if material.Properties.BaseColorTexture != nil {
		blendFactor := .2 // 10% image, 90% emitted light color
		return vec3.Lerp(material.Properties.Albedo, material.Properties.EmittanceColor, blendFactor), nil
		//}
	} else {

		scatterDir := vec3.Add(hitRecord.Normal, vec3.RandomInUnitSphere(random).Normalized())
		if scatterDir.NearZero() {
			scatterDir = hitRecord.Normal
		}
		scatterRay := ray.New(hitRecord.Point, scatterDir)

		//color := vec3.Vec3{
		//	X: hitRecord.U,
		//	Y: 0.0,
		//	Z: hitRecord.V,
		//}
		return material.Properties.Albedo, scatterRay
	}
}
