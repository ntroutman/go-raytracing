package vec

import (
	"fmt"
	"math"
	"raytracing/rnd"
)

type Color = Vec3
type Point3 = Vec3

type Vec3 struct {
	X float64
	Y float64
	Z float64
}

func Random() Vec3 {
	return Vec3{rnd.Float64(), rnd.Float64(), rnd.Float64()}
}

func RandomRange(min, max float64) Vec3 {
	return Vec3{rnd.Float64Range(min, max), rnd.Float64Range(min, max), rnd.Float64Range(min, max)}
}

func RandomInUnitSphere() Vec3 {
	for {
		p := RandomRange(-1, 1)
		if p.LengthSquared() >= 1 {
			continue
		}
		//fmt.Printf("p: %v, len: %f\n", p, p.LengthSquared())

		return p
	}
}

func RandomUnitVector() Vec3 {
	return RandomInUnitSphere().Norm()
}

func RandomInHemiSphere(normal Vec3) Vec3 {
	dir := RandomInUnitSphere()
	if Dot(normal, dir) < 0 {
		dir = dir.Neg()
	}
	return dir
}

func (v *Vec3) String() string {
	return fmt.Sprintf("V<%s, %s, %s>", v.X, v.Y, v.Z)
}

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
