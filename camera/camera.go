package camera

import (
	"goraytracer/ray"
	"goraytracer/vec3"
	"math"
)

type Camera struct {
	origin          vec3.Vec3
	right           vec3.Vec3
	up              vec3.Vec3
	lowerLeftCorner vec3.Vec3
}

func New(eye vec3.Vec3, look vec3.Vec3, verticalFovDegrees float64, aspectRatio float64) Camera {
	theta := verticalFovDegrees * 0.0174533
	h := math.Tan(theta / 2.0)

	viewportHeight := 2.0 * h
	viewportWidth := viewportHeight * aspectRatio

	up := vec3.Vec3{Y: 1}

	w := vec3.Sub(eye, look).Normalized()
	u := vec3.Cross(up, w).Normalized()
	v := vec3.Cross(w, u)

	camera := Camera{}
	camera.origin = eye
	camera.right = vec3.MultiplyScalar(u, viewportWidth)
	camera.up = vec3.MultiplyScalar(v, viewportHeight)

	// lowerLeftCorner = origin - camera.right/2 - camera.up/2 - w
	camera.lowerLeftCorner = camera.origin
	camera.lowerLeftCorner = vec3.Sub(camera.lowerLeftCorner, vec3.MultiplyScalar(camera.right, .5))
	camera.lowerLeftCorner = vec3.Sub(camera.lowerLeftCorner, vec3.MultiplyScalar(camera.up, .5))
	camera.lowerLeftCorner = vec3.Sub(camera.lowerLeftCorner, w)

	return camera
}

func (camera *Camera) GetRay(u float64, v float64) *ray.Ray {

	// direction = lowerLeftCorner + (camera.right * u) + (camera.up * v) - camera.origin
	direction := camera.lowerLeftCorner
	direction = vec3.Add(direction, vec3.MultiplyScalar(camera.right, u))
	direction = vec3.Add(direction, vec3.MultiplyScalar(camera.up, v))
	direction = vec3.Sub(direction, camera.origin)

	return ray.New(camera.origin, direction.Normalized())
}
