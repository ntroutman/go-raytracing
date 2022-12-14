package hittable

import (
	"math"
	"raytracing/mat"
	"raytracing/ray"
	"raytracing/vec"
)

type Sphere struct {
	Center vec.Point3
	Radius float64
	Mat    mat.Material
}

func (s *Sphere) Hit(r *ray.Ray, tMin, tMax float64) (Hit *HitRecord, IsHit bool) {
	oc := r.Origin.Sub(s.Center)
	a := r.Direction.LengthSquared()
	half_b := vec.Dot(oc, r.Direction)
	c := oc.LengthSquared() - s.Radius*s.Radius
	discriminant := half_b*half_b - a*c
	if discriminant < 0 {
		return nil, false
	}

	sqrtd := math.Sqrt(discriminant)

	// Find the nearest root that lies in the acceptable range.
	var root = (-half_b - sqrtd) / a
	//if root < tMin || root > tMax {
	if root < tMin || tMax < root {
		root = (-half_b + sqrtd) / a
		//if root < tMin || root > tMax {
		if root < tMin || tMax < root {
			// fmt.Printf("root: %v\n", root)
			return nil, false
		}
	}

	t := root
	//n := r.At(t).Sub(s.Center).Norm()
	p := r.At(t)
	n := p.Sub(s.Center).Scale(1.0 / s.Radius)
	return createHit(r, t, n, s.Mat), true
}
