package geometry

import "goraytracer/ray"

type Polygon struct {
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

// https://stackoverflow.com/questions/63264153/go-interface-method-returning-a-pointer-to-itself
// I think I need to use pointer receivers
func (polygon Polygon) AABBIntersections(aabb AABB) []Geometry {
	intersections := make([]Triangle, 0)
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
