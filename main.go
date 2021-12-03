package main

import (
	"fmt"
	"goraytracer/camera"
	"goraytracer/ray"
	"goraytracer/vec3"
	"math/rand"
	"os"
	"strings"
	"time"
)

type pixel struct {
	r int
	g int
	b int
}

func buildPPM(imageWidth int, imageHeight int, pixels []pixel) string {
	ppmHeader := fmt.Sprintf("P3\n%d %d\n255\n", imageWidth, imageHeight)
	var ppmBody strings.Builder

	for _, pixel := range pixels {
		ppmBody.WriteString(fmt.Sprintf("%d %d %d\n", pixel.r, pixel.g, pixel.b))
	}

	ppmFile := ppmHeader + ppmBody.String()
	return ppmFile
}

func writePPM(s string) {
	os.WriteFile("image.ppm", []byte(s), 0644)
}

func rayColor(ray *ray.Ray, depth int) vec3.Vec3 {
	if depth <= 0 {
		return vec3.Vec3{}
	}

	unitDirection := ray.Direction.Normalize()
	t := .5 * (unitDirection.Y + 1.0)
	return vec3.Add(
		vec3.MultiplyScalar(vec3.Vec3{X: 1.0, Y: 1.0, Z: 1.0}, 1.0-t),
		vec3.MultiplyScalar(vec3.Vec3{X: .5, Y: .7, Z: 1.0}, t),
	)
}

func lerp(v0 float64, v1 float64, t float64) float64 {
	return (1.0-t)*v0 + t*v1
}

func main() {
	randSeed := rand.NewSource(time.Now().UnixMicro())
	r1 := rand.New(randSeed)

	maxDepth := 50
	samplesPerPixel := 50
	aspectRatio := 4.0 / 3.0
	imageWidth := 640
	imageHeight := int(float64(imageWidth) / aspectRatio)

	pixels := make([]pixel, imageWidth*imageHeight)

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
				pixelColor = vec3.Add(pixelColor, rayColor(&ray, maxDepth))
			}

			// average pixel color
			pixelColor = vec3.MultiplyScalar(pixelColor, 1.0/float64(samplesPerPixel))

			// map to range [0-255]
			pixel := pixel{
				r: int(lerp(0, 255, pixelColor.X)),
				g: int(lerp(0, 255, pixelColor.Y)),
				b: int(lerp(0, 255, pixelColor.Z)),
			}

			// save
			pixels[index] = pixel
			index++
		}
	}

	writePPM(buildPPM(imageWidth, imageHeight, pixels))
}
