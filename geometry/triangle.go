package geometry

import (
	"goraytracer/ray"
	"goraytracer/vec3"
)

type Triangle struct {
	P1     vec3.Vec3
	P2     vec3.Vec3
	P3     vec3.Vec3
	Normal vec3.Vec3
}

func (triangle *Triangle) ComputeNormal() {

}

func (triangle Triangle) Hit(ray *ray.Ray, minDistance float64, maxDistance float64) HitRecord {
	// ray - triangle intersection

	// did this intersection happen within min and max distance?

	return HitRecord{Normal: triangle.Normal}
}
