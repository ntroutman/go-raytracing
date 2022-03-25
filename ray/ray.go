package ray

import (
	"fmt"
	"raytracing/vec"
)

type Ray struct {
	Origin    vec.Point3
	Direction vec.Point3
}

func (r *Ray) At(t float64) vec.Point3 {
	return r.Origin.Add(r.Direction.Scale(t))
}

func (r *Ray) String() string {
	return fmt.Sprintf("R{%s->%s}", &r.Origin, &r.Direction)
}
