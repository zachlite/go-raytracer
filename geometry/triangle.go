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

func edges(p1 vec3.Vec3, p2 vec3.Vec3, p3 vec3.Vec3) (vec3.Vec3, vec3.Vec3) {
	return vec3.Sub(p2, p1), vec3.Sub(p3, p1)
}

func NewTriangle(p1 vec3.Vec3, p2 vec3.Vec3, p3 vec3.Vec3) Triangle {

	// TODO: confirm that winding order is correct
	// TODO: be consistent when calculating edges for ray-triangle insersection

	edge1, edge2 := edges(p1, p2, p3)
	normal := vec3.Cross(edge1, edge2).Normalized()
	return Triangle{
		P1:     p1,
		P2:     p2,
		P3:     p3,
		Normal: normal,
	}
}

func (triangle Triangle) Hit(ray *ray.Ray, minDistance float64, maxDistance float64) HitRecord {
	// ray - triangle intersection
	// https://en.wikipedia.org/wiki/M%C3%B6ller%E2%80%93Trumbore_intersection_algorithm

	const EPSILON = 0.0000001

	edge1, edge2 := edges(triangle.P1, triangle.P2, triangle.P3)

	h := vec3.Cross(ray.Direction, edge2)
	a := vec3.Dot(edge1, h)

	if a > -EPSILON && a < EPSILON {
		return HitRecord{Hit: false} // ray is parallel to triangle
	}

	f := 1.0 / a
	s := vec3.Sub(ray.Origin, triangle.P1)
	u := f * vec3.Dot(s, h)

	if u < 0.0 || u > 1.0 {
		return HitRecord{Hit: false}
	}

	q := vec3.Cross(s, edge1)
	v := f * vec3.Dot(ray.Direction, q)

	if v < 0.0 || u+v > 1.0 {
		return HitRecord{Hit: false}
	}

	t := f * vec3.Dot(edge2, q)
	if t > EPSILON && t >= minDistance && t <= maxDistance {
		return HitRecord{
			Hit:      true,
			Distance: t,
			Point:    ray.At(t),
			Normal:   triangle.Normal,
		}
	}

	return HitRecord{Hit: false}
}

func pointInAABB(point vec3.Vec3, aabb AABB) bool {
	return point.X >= aabb.Min.X && point.X <= aabb.Max.X &&
		point.Y >= aabb.Min.Y && point.Y <= aabb.Max.Y &&
		point.Z >= aabb.Min.Z && point.Z <= aabb.Max.Z
}

func (triangle Triangle) AABBIntersections(aabb AABB) []Geometry {
	if triangle.IntersectsAABB(aabb) {
		intersections := make([]Geometry, 1)
		intersections[0] = triangle
		return intersections
	}

	return nil
}

func (triangle Triangle) IntersectsAABB(aabb AABB) bool {
	return pointInAABB(triangle.P1, aabb) &&
		pointInAABB(triangle.P2, aabb) &&
		pointInAABB(triangle.P3, aabb)
}
