package mat

import (
	"raytracing/ray"
	"raytracing/vec"
)

type Scatter struct {
	Attenuation vec.Color
	Scattered   ray.Ray
}

type Material interface {
	Scatter(ray *ray.Ray, p vec.Point3, normal vec.Vec3) (scatter Scatter, didScatter bool)
}
