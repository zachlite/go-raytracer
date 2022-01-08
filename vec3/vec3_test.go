package vec3_test

import (
	"goraytracer/vec3"
	"math"
	"math/rand"
	"testing"
)

func TestNormalizeVectorLengthIsOne(t *testing.T) {
	var vectors [3]vec3.Vec3

	vectors[0] = vec3.Vec3{.1, .1, .1}
	vectors[1] = vec3.Vec3{-12.0, -432.3, 132.0}
	vectors[2] = vec3.Vec3{1, 0, 0}

	nearOne := func(x float64) bool {
		return math.Abs(x-1.0) < .0000001
	}

	for i, vector := range vectors {
		l := vector.Normalized().Length()
		if !nearOne(l) {
			t.Error(i, l)
		}
	}
}

func TestRandomVectorInUnitSphereHasLength_LTE_One(t *testing.T) {
	randSource := rand.NewSource(0)
	r := rand.New(randSource)

	for i := 0; i < 10000; i++ {

		randomVec := vec3.RandomInUnitSphere(r)
		if randomVec.Length() > 1 {
			t.FailNow()
		}
	}
}
