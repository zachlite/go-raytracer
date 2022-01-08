package vec3

import (
	"math"
	"math/rand"
)

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

func (v Vec3) Normalized() Vec3 {
	length := Dot(v, v)
	if length > 0 {
		length = 1.0 / math.Sqrt(length)
	}

	return Vec3{X: v.X * length, Y: v.Y * length, Z: v.Z * length}
}

func (v Vec3) Length() float64 {
	return math.Sqrt(Dot(v, v))
}

func (v *Vec3) NearZero() bool {
	limit := .00000001
	return math.Abs(v.X) < limit && math.Abs(v.Y) < limit && math.Abs(v.Z) < limit
}

// return a vector with length between 0 and 1
func RandomInUnitSphere(r *rand.Rand) Vec3 {
	var u = r.Float64()
	var v = r.Float64()
	var theta = u * 2.0 * math.Pi
	var phi = math.Acos(2.0*v - 1.0)
	var _r = math.Cbrt(r.Float64())
	var sinTheta = math.Sin(theta)
	var cosTheta = math.Cos(theta)
	var sinPhi = math.Sin(phi)
	var cosPhi = math.Cos(phi)
	var x = _r * sinPhi * cosTheta
	var y = _r * sinPhi * sinTheta
	var z = _r * cosPhi
	return Vec3{x, y, z}
}
