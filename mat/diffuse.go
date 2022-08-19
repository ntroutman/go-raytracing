package mat

import (
	"raytracing/ray"
	"raytracing/vec"
)

type Lambertian struct {
	Albedo vec.Color
}

func (mat *Lambertian) Scatter(
	r *ray.Ray,
	p vec.Point3,
	normal vec.Vec3,
) (Scatter, bool) {
	scatter := new(Scatter)
	scatter.Scattered = ray.Ray{p, vec.RandomInHemiSphere(normal)}

	// Catch degenerate scatter direction
	if scatter.Scattered.Direction.NearZero() {
		scatter.Scattered.Direction = normal
	}
	scatter.Attenuation = mat.Albedo
	return *scatter, true
}
