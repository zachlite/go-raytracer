package mesh

import (
	"goraytracer/geometry"
	"goraytracer/material"
)

type Mesh struct {
	Geometry geometry.Geometry
	Material material.Material
}
