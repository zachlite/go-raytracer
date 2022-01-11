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
