package geometry

import "goraytracer/ray"

type Polygon struct {
	Id        uint32
	Triangles []Triangle
}

func (polygon Polygon) Hit(r *ray.Ray, minDistance float64, maxDistance float64) HitRecord {

	for _, triangle := range polygon.Triangles {
		hit := triangle.Hit(r, minDistance, maxDistance)
		if hit.Hit {
			return hit
		}
	}

	return HitRecord{Hit: false}
}

func (polygon Polygon) AABBIntersections(aabb AABB) []Geometry {
	intersections := make([]Geometry, 0)
	for _, triangle := range polygon.Triangles {
		if triangle.IntersectsAABB(aabb) {
			intersections = append(intersections, triangle)
		}

	}
	return intersections
}

func (polygon Polygon) IntersectsAABB(aabb AABB) bool {
	for _, triangle := range polygon.Triangles {
		if triangle.IntersectsAABB(aabb) {
			return true
		}
	}
	return false
}

func (polygon Polygon) GetId() uint32 {
	return polygon.Id
}
