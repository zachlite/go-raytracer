package vec3

import "math"

type Vec3 struct {
	X, Y, Z float64
}

func Dot(a Vec3, b Vec3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func Add(a Vec3, b Vec3) Vec3 {
	return Vec3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// func AddScalar(a Vec3, s float64) Vec3 {
// 	return Vec3{a.X + s, a.Y + s, a.Z + s}
// }

func Sub(a Vec3, b Vec3) Vec3 {
	return Vec3{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func Multiply(a Vec3, b Vec3) Vec3 {
	return Vec3{
		X: a.X * b.X,
		Y: a.Y * b.Y,
		Z: a.Z * b.Z,
	}
}

func MultiplyScalar(v Vec3, s float64) Vec3 {
	return Vec3{v.X * s, v.Y * s, v.Z * s}
}

func (v *Vec3) Normalize() Vec3 {
	length := Dot(*v, *v)
	if length > 0 {
		length = 1.0 / math.Sqrt(length)
	}

	return Vec3{X: v.X * length, Y: v.Y * length, Z: v.Z * length}
}
