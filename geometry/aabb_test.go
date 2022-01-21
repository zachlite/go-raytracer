package geometry

import (
	"goraytracer/ray"
	"goraytracer/vec3"
	"testing"
)

func TestAABB_IntersectsRay(t *testing.T) {

	testAABB := AABB{
		Min: vec3.Vec3{X: -1, Y: -1, Z: -1},
		Max: vec3.Vec3{X: 1, Y: 1, Z: 1},
	}

	type args struct {
		ray *ray.Ray
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Intersection: origin outside",
			args: args{ray: ray.New(vec3.Vec3{Z: -2}, vec3.Vec3{Z: 1})},
			want: true,
		},
		{
			name: "Intersection: origin inside",
			args: args{ray: ray.New(vec3.Vec3{}, vec3.Vec3{Z: 1})},
			want: true,
		},
		{
			name: "No intersection",
			args: args{ray: ray.New(vec3.Vec3{Z: -2}, vec3.Vec3{Y: 1})},
			want: false,
		},
		{
			name: "No intersection along planes",
			args: args{ray: ray.New(vec3.Vec3{X: 1, Z: -2}, vec3.Vec3{Z: 1})},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := testAABB.IntersectsRay(tt.args.ray); got != tt.want {
				t.Errorf("IntersectsRay() = %v, want %v", got, tt.want)
			}
		})
	}
}
