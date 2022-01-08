package mathutils

func Lerp(v0 float64, v1 float64, t float64) float64 {
	return (1.0-t)*v0 + t*v1
}
