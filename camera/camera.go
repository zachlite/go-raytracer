package camera

import (
	"goraytracer/ray"
	"goraytracer/vec3"
)

type Camera struct {
	origin          vec3.Vec3
	right           vec3.Vec3
	up              vec3.Vec3
	lowerLeftCorner vec3.Vec3
}

// TODO: define camera in world space coordinates
func (camera *Camera) Init(aspectRatio float64) {
	viewportHeight := 2.0
	viewportWidth := viewportHeight * aspectRatio
	focalLength := 1.0

	camera.origin = vec3.Vec3{}
	camera.right = vec3.Vec3{X: viewportWidth, Y: 0, Z: 0}
	camera.up = vec3.Vec3{X: 0, Y: viewportHeight, Z: 0}
	camera.lowerLeftCorner = vec3.Sub(
		vec3.Sub(camera.origin, vec3.MultiplyScalar(camera.right, 0.5)),
		vec3.Sub(vec3.MultiplyScalar(camera.up, 0.5), vec3.Vec3{X: 0, Y: 0, Z: focalLength}))
}

func (camera *Camera) GetRay(u float64, v float64) *ray.Ray {
	direction := vec3.Add(
		camera.lowerLeftCorner,
		vec3.Add(
			vec3.MultiplyScalar(camera.right, u),
			vec3.MultiplyScalar(camera.up, v))).Normalized()
	return ray.New(camera.origin, direction)
}
