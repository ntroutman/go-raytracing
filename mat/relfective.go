package mat

import (
	"raytracing/ray"
	"raytracing/vec"
)

type Metal struct {
	Albedo    vec.Color
	Fuzziness float64
}

func (mat *Metal) Scatter(
	r *ray.Ray,
	p vec.Point3,
	normal vec.Vec3,
) (Scatter, bool) {
	reflected := vec.Reflect(r.Direction.Norm(), normal)
	if mat.Fuzziness > 0 {
		reflected = reflected.Add(vec.RandomInUnitSphere().Scale(mat.Fuzziness))
	}
	scattered := ray.Ray{p, reflected}

	// Catch degenerate scatter direction
	if scattered.Direction.NearZero() {
		scattered.Direction = normal
	}
	attenuation := mat.Albedo
	didScatter := vec.Dot(scattered.Direction, normal) > 0
	return Scatter{Scattered: scattered, Attenuation: attenuation}, didScatter
}
