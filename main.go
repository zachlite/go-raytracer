package main

import (
	"goraytracer/camera"
	"goraytracer/material"
	"goraytracer/ppm"
	"goraytracer/ray"
	"goraytracer/sphere"
	"goraytracer/vec3"
	"math/rand"
	"time"
)

func lerp(v0 float64, v1 float64, t float64) float64 {
	return (1.0-t)*v0 + t*v1
}

func rayColor(spheres []sphere.Sphere, ray *ray.Ray, depth int) vec3.Vec3 {
	if depth <= 0 {
		return vec3.Vec3{}
	}

	hitRecord := sphere.FindClosestSphereHit(spheres, ray)

	// scatter and recurse if there's a hit record
	if hitRecord.Hit {
		attenuation, scatteredRay := hitRecord.Material.Scatter(*ray, hitRecord.Point, hitRecord.Normal)
		return vec3.Multiply(attenuation, rayColor(spheres, &scatteredRay, depth-1))
	}

	// if there's no sphere hit, render the sky
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

	maxDepth := 10
	samplesPerPixel := 100
	aspectRatio := 4.0 / 3.0
	imageWidth := 640
	imageHeight := int(float64(imageWidth) / aspectRatio)

	// define our scene
	spheres := make([]sphere.Sphere, 3)

	spheres[0] = sphere.Sphere{
		Center:   vec3.Vec3{X: 0, Y: 0, Z: 1},
		Radius:   .5,
		Material: material.Lambertian{Albedo: vec3.Vec3{X: .7, Y: .7, Z: .7}}}

	spheres[1] = sphere.Sphere{
		Center:   vec3.Vec3{X: 1, Y: 0, Z: 1},
		Radius:   .5,
		Material: material.Lambertian{Albedo: vec3.Vec3{X: 0.8, Y: .1, Z: .2}}}

	spheres[2] = sphere.Sphere{
		Center:   vec3.Vec3{X: 0, Y: -100.5, Z: 0},
		Radius:   100,
		Material: material.Lambertian{Albedo: vec3.Vec3{X: 0.0, Y: .7, Z: .7}}}

	pixels := make([]ppm.Pixel, imageWidth*imageHeight)

	var camera camera.Camera
	camera.Init(aspectRatio)
	camera.GetRay(1.0, 1.0)

	index := 0
	for j := imageHeight - 1; j >= 0; j-- {
		println(j)
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
