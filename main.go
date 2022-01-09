package main

import (
	"flag"
	"fmt"
	"goraytracer/camera"
	"goraytracer/geometry"
	"goraytracer/material"
	"goraytracer/mathutils"
	"goraytracer/ppm"
	"goraytracer/ray"
	"goraytracer/vec3"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
	"time"
)

type Mesh struct {
	Geometry geometry.Geometry
	Material material.Material
}

func findClosestMeshHit(meshes []Mesh, ray *ray.Ray) (geometry.HitRecord, material.Material) {
	minDistance := .001
	maxDistance := math.Inf(1)

	closetHit := geometry.HitRecord{Hit: false}
	var material material.Material

	for _, mesh := range meshes {
		hitRecord := mesh.Geometry.Hit(ray, minDistance, maxDistance)
		if hitRecord.Hit {
			maxDistance = hitRecord.Distance
			closetHit = hitRecord
			material = mesh.Material
		}

	}

	return closetHit, material
}

func rayColor(meshes []Mesh, ray *ray.Ray, depth int, random *rand.Rand) vec3.Vec3 {
	if depth <= 0 {
		return vec3.Vec3{}
	}

	hitRecord, material := findClosestMeshHit(meshes, ray)

	// scatter and recurse if there's a hit record
	if hitRecord.Hit {
		attenuation, scatteredRay := material.Scatter(*ray, hitRecord.Point, hitRecord.Normal, random)
		return vec3.Multiply(attenuation, rayColor(meshes, &scatteredRay, depth-1, random))
	}

	// if there's no sphere hit, render the sky
	unitDirection := ray.Direction.Normalized()
	t := .5 * (unitDirection.Y + 1.0)
	return vec3.Add(
		vec3.MultiplyScalar(vec3.Vec3{X: 1.0, Y: 1.0, Z: 1.0}, 1.0-t),
		vec3.MultiplyScalar(vec3.Vec3{X: .5, Y: .7, Z: 1.0}, t),
	)
}

func samplePixel(i int, j int, imageWidth int, imageHeight int, camera camera.Camera, meshes []Mesh) vec3.Vec3 {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	const samplesPerPixel = 100
	const maxDepth = 10
	pixelColor := vec3.Vec3{}

	for sample := 0; sample < samplesPerPixel; sample++ {
		u := (float64(i) + r.Float64()) / (float64(imageWidth) - 1)
		v := (float64(j) + r.Float64()) / (float64(imageHeight) - 1)
		ray := camera.GetRay(u, v)
		pixelColor = vec3.Add(pixelColor, rayColor(meshes, &ray, maxDepth, r))
	}

	// average and gamma correct
	scale := 1.0 / float64(samplesPerPixel)
	pixelColor = vec3.MultiplyScalar(pixelColor, scale)
	pixelColor.X = math.Sqrt(pixelColor.X)
	pixelColor.Y = math.Sqrt(pixelColor.Y)
	pixelColor.Z = math.Sqrt(pixelColor.Z)

	return pixelColor
}

func main() {

	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	const aspectRatio = 4.0 / 3.0
	const imageWidth = 640
	const imageHeight = int(float64(imageWidth) / aspectRatio)

	// define our scene
	meshes := make([]Mesh, 4)

	meshes[0] = Mesh{
		Geometry: geometry.Sphere{
			Center: vec3.Vec3{X: 0, Y: 0, Z: 1},
			Radius: .5,
		},
		Material: material.Lambertian{
			Albedo: vec3.Vec3{X: .7, Y: .7, Z: .7},
		},
	}

	meshes[1] = Mesh{
		Geometry: geometry.Sphere{
			Center: vec3.Vec3{X: 1, Y: 0, Z: 1},
			Radius: .5,
		},
		Material: material.Lambertian{
			Albedo: vec3.Vec3{X: .8, Y: .1, Z: .2},
		},
	}

	meshes[2] = Mesh{
		Geometry: geometry.Sphere{
			Center: vec3.Vec3{X: 0, Y: -100.5, Z: 0},
			Radius: 100,
		},
		Material: material.Lambertian{
			Albedo: vec3.Vec3{X: 0.0, Y: .7, Z: .7},
		},
	}

	p1 := vec3.Vec3{X: -1.25, Y: -.5, Z: 1}
	p2 := vec3.Vec3{X: -1.0, Y: .4, Z: 1}
	p3 := vec3.Vec3{X: -.75, Y: 0.0, Z: 1}

	meshes[3] = Mesh{
		Geometry: geometry.NewTriangle(p1, p2, p3),
		Material: material.Lambertian{
			Albedo: vec3.Vec3{X: .8, Y: .1, Z: .2},
		},
	}

	frameBuffer := make([]ppm.Pixel, imageWidth*imageHeight)

	var camera camera.Camera
	camera.Init(aspectRatio)
	camera.GetRay(1.0, 1.0)
	split := 8

	wg := sync.WaitGroup{}

	for xRegion := 0; xRegion < split; xRegion++ {
		for yRegion := 0; yRegion < split; yRegion++ {
			wg.Add(1)

			iMin := imageWidth / split * xRegion    // 0 - 560
			iMax := iMin + (imageWidth / split) - 1 // 79 - 639

			jMin := imageHeight / split * yRegion
			jMax := jMin + (imageHeight / split) - 1

			go func(iMin int, iMax int, jMin int, jMax int) {
				defer func() {
					fmt.Println("region done: ", iMin, iMax, jMin, jMax)
					wg.Done()
				}()

				for i := iMin; i <= iMax; i++ {
					for j := jMin; j <= jMax; j++ {

						color := samplePixel(i, j, imageWidth, imageHeight, camera, meshes)

						index := (imageHeight-1-j)*imageWidth + i

						// map to range [0-255]
						frameBuffer[index] = ppm.Pixel{
							R: int(mathutils.Lerp(0, 255, color.X)),
							G: int(mathutils.Lerp(0, 255, color.Y)),
							B: int(mathutils.Lerp(0, 255, color.Z)),
						}

					}
				}

			}(iMin, iMax, jMin, jMax)
		}
	}

	wg.Wait()
	// done

	ppm.Write("image.ppm", ppm.Build(imageWidth, imageHeight, frameBuffer))
}
