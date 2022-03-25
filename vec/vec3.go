package vec

import (
	"fmt"
	"math"
)

type Color = Vec3
type Point3 = Vec3

type Vec3 struct {
	X float64
	Y float64
	Z float64
}

func (v *Vec3) String() string {
	return fmt.Sprintf("V<%s, %s, %s>", v.X, v.Y, v.Z)
}

// func New(x, y, z float32) Vec3 {
// 	return Vec3{[3]float32{x, y, z}}
// }

// func (v *Vec3) x() float32 {
// 	return v.e[0]
// }

// func (v *Vec3) y() float32 {
// 	return v.e[1]
// }

// func (v *Vec3) z() float32 {
// 	return v.e[1]
// }

func (v Vec3) Scale(s float64) Vec3 {
	return Vec3{v.X * s, v.Y * s, v.Z * s}
}

func (v Vec3) Add(w Vec3) Vec3 {
	return Vec3{v.X + w.X, v.Y + w.Y, v.Z + w.Z}
}

func (v Vec3) Sub(w Vec3) Vec3 {
	return Vec3{v.X - w.X, v.Y - w.Y, v.Z - w.Z}
}

func (v Vec3) LengthSquared() float64 {
	//return v.e[0]*v.e[0] + v.e[1]*v.e[1] + v.e[2]*v.e[2]
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vec3) Length() float64 {
	return math.Sqrt(float64(v.LengthSquared()))
}

func (v Vec3) Neg() Vec3 {
	return Vec3{-v.X, -v.Y, -v.Z}
}

func (v Vec3) Norm() Vec3 {
	mag := v.Length()
	return Vec3{v.X / mag, v.Y / mag, v.Z / mag}
}

func Dot(a Vec3, b Vec3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}
