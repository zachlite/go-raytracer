package geometry

import (
	"goraytracer/ray"
	"goraytracer/vec3"
	"math"
)

type Sphere struct {
	Center vec3.Vec3
	Radius float64
}

func (s Sphere) Hit(r *ray.Ray, minDistance float64, maxDistance float64) HitRecord {

	// solve the quadratic equation to see if ray intersects sphere at all
	originToCenter := vec3.Sub(r.Origin, s.Center)
	a := vec3.Dot(r.Direction, r.Direction)
	halfB := vec3.Dot(originToCenter, r.Direction)
	c := vec3.Dot(originToCenter, originToCenter) - (s.Radius * s.Radius)
	discriminant := halfB*halfB - (a * c)
	if discriminant < 0 {
		return HitRecord{Hit: false}
	}

	// assert that intersection lies between min and max distance bounds
	sqrtD := math.Sqrt(discriminant)
	root := (-halfB - sqrtD) / a
	if root < minDistance || maxDistance < root {
		root = (-halfB + sqrtD) / a
		if root < minDistance || maxDistance < root {
			return HitRecord{Hit: false}
		}
	}

	distance := root
	point := r.At(distance)
	outwardNormal := vec3.MultiplyScalar(vec3.Sub(point, s.Center), 1.0/s.Radius).Normalized()
	isFrontFace := vec3.Dot(r.Direction, outwardNormal) < 0

	var normal vec3.Vec3

	if isFrontFace {
		normal = outwardNormal
	} else {
		normal = vec3.MultiplyScalar(outwardNormal, -1.0)
	}

	return HitRecord{Hit: true, Distance: distance, Point: point, Normal: normal}
}

func (s Sphere) AABBIntersections(aabb AABB) []*Sphere {
	intersections := make([]*Sphere, 0)

	if s.IntersectsAABB(aabb) {
		intersections = append(intersections, &s)
	}

	return intersections
}

func (s Sphere) IntersectsAABB(aabb AABB) bool {

	return false

}
