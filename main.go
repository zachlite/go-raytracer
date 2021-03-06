package main

import (
	"flag"
	"fmt"
	"goraytracer/accel"
	"goraytracer/camera"
	"goraytracer/geometry"
	"goraytracer/material"
	"goraytracer/mathutils"
	"goraytracer/mesh"
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

func findClosestMeshHit(candidates []accel.IntersectCandidate, ray *ray.Ray) (geometry.HitRecord, material.Material) {
	minDistance := .001
	maxDistance := math.Inf(1)

	closetHit := geometry.HitRecord{Hit: false}
	var material material.Material

	for _, candidate := range candidates {
		hitRecord := candidate.Geometry.Hit(ray, minDistance, maxDistance)
		if hitRecord.Hit {
			maxDistance = hitRecord.Distance
			closetHit = hitRecord
			material = candidate.Material
		}

	}

	return closetHit, material
}

func rayColor(tree *accel.OctTree, ray *ray.Ray, depth int, random *rand.Rand) vec3.Vec3 {
	if depth <= 0 {
		return vec3.Vec3{}
	}

	candidates := tree.Search(ray)
	hitRecord, material := findClosestMeshHit(candidates, ray)

	// scatter and recurse if there's a hit record
	if hitRecord.Hit {
		attenuation, scatteredRay := material.Scatter(hitRecord, random)

		if scatteredRay != nil {
			return vec3.Multiply(attenuation, rayColor(tree, scatteredRay, depth-1, random))
		} else {
			return attenuation
		}
	}

	return vec3.Vec3{}
	// if there's no sphere hit, render the sky
	//unitDirection := ray.Direction.Normalized()
	//t := .5 * (unitDirection.Y + 1.0)
	//return vec3.Add(
	//	vec3.MultiplyScalar(vec3.Vec3{X: 1.0, Y: 1.0, Z: 1.0}, 1.0-t),
	//	vec3.MultiplyScalar(vec3.Vec3{X: .5, Y: .7, Z: 1.0}, t),
	//)
}

func samplePixel(i int, j int, imageWidth int, imageHeight int, camera camera.Camera, tree *accel.OctTree) vec3.Vec3 {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	const samplesPerPixel = 50
	const maxDepth = 10
	pixelColor := vec3.Vec3{}

	for sample := 0; sample < samplesPerPixel; sample++ {
		u := (float64(i) + r.Float64()) / (float64(imageWidth) - 1)
		v := (float64(j) + r.Float64()) / (float64(imageHeight) - 1)
		ray := camera.GetRay(u, v)
		pixelColor = vec3.Add(pixelColor, rayColor(tree, ray, maxDepth, r))
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
	const imageWidth = 320
	const imageHeight = int(float64(imageWidth) / aspectRatio)

	// define our scene
	meshes := []mesh.Mesh{
		{
			Geometry: geometry.Sphere{
				Id:     0,
				Center: vec3.Vec3{X: -51, Y: 0, Z: 0},
				Radius: 50,
			},
			Material: &material.Lambertian{Properties: material.MaterialProps{
				Albedo:           vec3.Vec3{1, 0, 0},
				BaseColorTexture: nil,
				EmittanceColor:   vec3.Vec3{1, 1, 1},
			}},
		},
		{
			Geometry: geometry.Sphere{
				Id:     1,
				Center: vec3.Vec3{20, 0, -5},
				Radius: 2,
			},
			Material: &material.Lambertian{Properties: material.MaterialProps{
				Albedo:           vec3.Vec3{.5, .5, .5},
				BaseColorTexture: nil,
				EmittanceColor:   vec3.Vec3{0, 0, 0},
			}},
		},
	}

	startTime := time.Now().UnixMicro()
	tree := accel.BuildOctTree(meshes)
	endTime := time.Now().UnixMicro()
	fmt.Printf("OctTree built in %f seconds\n", float64(endTime-startTime)/1e6)

	frameBuffer := make([]ppm.Pixel, imageWidth*imageHeight)

	cam := camera.New(vec3.Vec3{X: 0, Y: 0, Z: 100}, vec3.Vec3{}, 45, aspectRatio)

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

						color := samplePixel(i, j, imageWidth, imageHeight, cam, &tree)

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
