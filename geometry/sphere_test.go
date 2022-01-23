package geometry

import (
	"fmt"
	"goraytracer/vec3"
	"testing"
)

// TODO: test cartesian values above and below xz plane
func TestSphere_GetUV(t *testing.T) {
	sphere := Sphere{
		Id:     0,
		Center: vec3.Vec3{},
		Radius: 1,
	}
	type args struct {
		point vec3.Vec3
	}
	tests := []struct {
		name  string
		args  args
		wantU float64
		wantV float64
	}{
		{
			name:  "cartesian coordinates 0,0,1",
			args:  args{point: vec3.Vec3{0, 0, 1}.Normalized()},
			wantU: 1.0,
			wantV: 0.5,
		},
		{
			name:  "cartesian coordinates 1,0,0",
			args:  args{point: vec3.Vec3{1, 0, 0}.Normalized()},
			wantU: 0.25,
			wantV: 0.5,
		},
		{
			name:  "cartesian coordinates 0,0,-1",
			args:  args{point: vec3.Vec3{0, 0, -1}.Normalized()},
			wantU: 0.5,
			wantV: 0.5,
		},
		{
			name:  "cartesian coordinates -1,0,0",
			args:  args{point: vec3.Vec3{-1, 0, 0}.Normalized()},
			wantU: 0.75,
			wantV: 0.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotU, gotV := sphere.GetUV(tt.args.point)
			fmt.Println(gotU, gotV)
			if gotU != tt.wantU {
				t.Errorf("GetUV() gotU = %v, want %v", gotU, tt.wantU)
			}
			if gotV != tt.wantV {
				t.Errorf("GetUV() gotV = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}
