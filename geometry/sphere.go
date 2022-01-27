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

/**
Sphere uv thoughts:

on sphere hit:
hit euclidean coordinates -> polar coordinages -> uv coordinates
save uv coordinates in hit record

in Material.Scatter, use uv from hit record to look up image data.

Image data:
- width
- height

can create images manually for testing
or procedurally with things like perlin noise

Multiple images can allow for cool material properties
- albedo map
- glossiness map
- roughness map
- normal map

*/

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
	u, v := s.GetUV(point)
	outwardNormal := vec3.MultiplyScalar(vec3.Sub(point, s.Center), 1.0/s.Radius).Normalized()
	isFrontFace := vec3.Dot(r.Direction, outwardNormal) < 0

	var normal vec3.Vec3

	if isFrontFace {
		normal = outwardNormal
	} else {
		normal = vec3.MultiplyScalar(outwardNormal, -1.0)
	}

	return HitRecord{Hit: true, Distance: distance, Point: point, Normal: normal, U: u, V: v}
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
}

func (s Sphere) GetUV(point vec3.Vec3) (u float64, v float64) {
	P := vec3.Sub(s.Center, point).Normalized()
	u = 0.5 + (math.Atan2(P.X, P.Z) / (2.0 * math.Pi))
	v = 0.5 - (math.Asin(P.Y) / math.Pi)
	return u, v
}

func (s Sphere) GetId() uint32 {
	return s.Id
}

//https://thebookofshaders.com/13/
