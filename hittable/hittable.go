package hittable

import (
	"raytracing/ray"
	"raytracing/vec"
)

type HitRecord struct {
	P           vec.Point3
	Normal      vec.Vec3
	T           float64
	IsFrontFace bool
}

type Hittable interface {
	Hit(r *ray.Ray, tMin, tMax float64) (Hit *HitRecord, IsHit bool)
}

func createHit(r *ray.Ray, t float64, outwardNormal vec.Vec3) (Hit *HitRecord) {
	Hit = new(HitRecord)
	Hit.P = r.At(t)
	Hit.T = t
	Hit.IsFrontFace = vec.Dot(r.Direction, outwardNormal) < 0
	if Hit.IsFrontFace {
		Hit.Normal = outwardNormal
	} else {
		Hit.Normal = outwardNormal.Neg()
	}
	return
}

type HittableList struct {
	objects []Hittable
}

func (list *HittableList) Hit(r *ray.Ray, tMin, tMax float64) (Hit *HitRecord, IsHit bool) {
	var closestHitDistance = tMax
	var closestObj *HitRecord = nil
	for _, obj := range list.objects {
		hit, isHit := obj.Hit(r, tMin, closestHitDistance)
		if isHit {
			closestObj = hit
			closestHitDistance = hit.T
		}
	}
	return closestObj, closestObj != nil
}

func (list *HittableList) Add(obj Hittable) {
	list.objects = append(list.objects, obj)
}
