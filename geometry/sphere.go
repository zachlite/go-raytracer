package geometry

import (
	"goraytracer/ray"
	"goraytracer/vec3"
	"math"
)

type Sphere struct {
	Id     uint32
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

func pointInSphere(s Sphere, point vec3.Vec3) bool {
	return math.Sqrt(
		((point.X-s.Center.X)*(point.X-s.Center.X))+
			((point.Y-s.Center.Y)*(point.Y-s.Center.Y))+
			((point.Z-s.Center.Z)*(point.Y-s.Center.Z))) < s.Radius
}

func (s Sphere) AABBIntersections(aabb AABB) []Geometry {
	intersections := make([]Geometry, 0)

	if s.IntersectsAABB(aabb) {
		intersections = append(intersections, &s)
	}

	return intersections
}

func (s Sphere) IntersectsAABB(aabb AABB) bool {

	x := math.Max(aabb.Min.X, math.Min(s.Center.X, aabb.Max.X))
	y := math.Max(aabb.Min.Y, math.Min(s.Center.Y, aabb.Max.Y))
	z := math.Max(aabb.Min.Z, math.Min(s.Center.Z, aabb.Max.Z))

	distanceSquared := ((x - s.Center.X) * (x - s.Center.X)) +
		((y - s.Center.Y) * (y - s.Center.Y)) +
		((z - s.Center.Z) * (z - s.Center.Z))

	return distanceSquared <= s.Radius*s.Radius

	//rr := s.Radius * s.Radius
	//dmin := 0.0
	//
	//if s.Center.X < aabb.Min.X {
	//	dmin += math.Sqrt(s.Center.X - aabb.Min.X)
	//} else if s.Center.X > aabb.Max.X {
	//	dmin += math.Sqrt(s.Center.X - aabb.Max.X)
	//}
	//
	//if s.Center.Y < aabb.Min.Y {
	//	dmin += math.Sqrt(s.Center.Y - aabb.Min.Y)
	//} else if s.Center.Y > aabb.Max.Y {
	//	dmin += math.Sqrt(s.Center.Y - aabb.Max.Y)
	//}
	//
	//if s.Center.Z < aabb.Min.Z {
	//	dmin += math.Sqrt(s.Center.Z - aabb.Min.Z)
	//} else if s.Center.Z > aabb.Max.Z {
	//	dmin += math.Sqrt(s.Center.Z - aabb.Max.Z)
	//}
	//
	//return dmin <= rr
}

func (s Sphere) GetId() uint32 {
	return s.Id
}
