package main

import (
	"goraytracer/camera"
	"goraytracer/ppm"
	"goraytracer/ray"
	"goraytracer/vec3"
	"math"
	"math/rand"
	"time"
)

// figure out different material types
// type Material struct {

// }

type HitRecord struct {
	Hit      bool
	Distance float64
	Point    vec3.Vec3
	Normal   vec3.Vec3
}

type Sphere struct {
	Center   vec3.Vec3
	Radius   float64
	Material int
}

func Hit(s *Sphere, r *ray.Ray, minDistance float64, maxDistance float64) HitRecord {

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
	outwardNormal := vec3.MultiplyScalar(vec3.Sub(point, s.Center), 1.0/s.Radius)
	outwardNormal = outwardNormal.Normalize()
	isFrontFace := vec3.Dot(r.Direction, outwardNormal) < 0

	var normal vec3.Vec3

	if isFrontFace {
		normal = outwardNormal
	} else {
		normal = vec3.MultiplyScalar(outwardNormal, -1.0)
	}

	return HitRecord{Hit: true, Distance: distance, Point: point, Normal: normal}
}

func lerp(v0 float64, v1 float64, t float64) float64 {
	return (1.0-t)*v0 + t*v1
}

func findClosestSphereHit(spheres []Sphere, ray *ray.Ray) HitRecord {
	minDistance := .001
	maxDistance := math.Inf(1)

	closetHit := HitRecord{Hit: false}

	for _, sphere := range spheres {
		hitRecord := Hit(&sphere, ray, minDistance, maxDistance)
		if hitRecord.Hit {
			maxDistance = hitRecord.Distance
			closetHit = hitRecord
		}
	}

	return closetHit
}

func rayColor(spheres []Sphere, ray *ray.Ray, depth int) vec3.Vec3 {
	if depth <= 0 {
		return vec3.Vec3{}
	}

	hitRecord := findClosestSphereHit(spheres, ray)

	if hitRecord.Hit {
		// scatter and recurse
		return vec3.Vec3{X: 1.0, Y: 0.0, Z: 0.0}
	}

	unitDirection := ray.Direction.Normalize()
	t := .5 * (unitDirection.Y + 1.0)
	return vec3.Add(
		vec3.MultiplyScalar(vec3.Vec3{X: 1.0, Y: 1.0, Z: 1.0}, 1.0-t),
		vec3.MultiplyScalar(vec3.Vec3{X: .5, Y: .7, Z: 1.0}, t),
	)
}

func main() {
	randSeed := rand.NewSource(time.Now().UnixMicro())
	r1 := rand.New(randSeed)

	maxDepth := 50
	samplesPerPixel := 50
	aspectRatio := 4.0 / 3.0
	imageWidth := 640
	imageHeight := int(float64(imageWidth) / aspectRatio)

	// define our scene
	spheres := make([]Sphere, 1)
	spheres[0] = Sphere{Center: vec3.Vec3{X: 0, Y: 0, Z: 1}, Radius: .5, Material: 1}

	pixels := make([]ppm.Pixel, imageWidth*imageHeight)

	var camera camera.Camera
	camera.Init(aspectRatio)
	camera.GetRay(1.0, 1.0)

	index := 0
	for j := imageHeight - 1; j >= 0; j-- {
		for i := 0; i < imageWidth; i++ {

			pixelColor := vec3.Vec3{}

			// sample
			for sample := 0; sample < samplesPerPixel; sample++ {
				u := (float64(i) + r1.Float64()) / (float64(imageWidth) - 1)
				v := (float64(j) + r1.Float64()) / (float64(imageHeight) - 1)
				ray := camera.GetRay(u, v)
				pixelColor = vec3.Add(pixelColor, rayColor(spheres, &ray, maxDepth))
			}

			// average pixel color
			pixelColor = vec3.MultiplyScalar(pixelColor, 1.0/float64(samplesPerPixel))

			// map to range [0-255]
			pixels[index] = ppm.Pixel{
				R: int(lerp(0, 255, pixelColor.X)),
				G: int(lerp(0, 255, pixelColor.Y)),
				B: int(lerp(0, 255, pixelColor.Z)),
			}

			index++
		}
	}

	ppm.Write("image.ppm", ppm.Build(imageWidth, imageHeight, pixels))
}
