package ppm

import (
	"fmt"
	"os"
	"strings"
)

type Pixel struct {
	R int
	G int
	B int
}

func Build(imageWidth int, imageHeight int, pixels []Pixel) string {
	ppmHeader := fmt.Sprintf("P3\n%d %d\n255\n", imageWidth, imageHeight)
	var ppmBody strings.Builder

	for _, pixel := range pixels {
		ppmBody.WriteString(fmt.Sprintf("%d %d %d\n", pixel.R, pixel.G, pixel.B))
	}

	return ppmHeader + ppmBody.String()
}

func Write(filename string, s string) {
	os.WriteFile(filename, []byte(s), 0644)
}
