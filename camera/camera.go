package camera

import (
	"raytracing/image"
	"raytracing/ray"
	"raytracing/vec"
)

type Camera struct {
	aspectRatio     float64
	origin          vec.Point3
	lowerLeftCorner vec.Point3
	horizontal      vec.Vec3
	vertical        vec.Vec3
}

func Create() *Camera {
	aspectRatio := 16.0 / 9.0
	viewportHeight := 2.0
	viewportWidth := aspectRatio * viewportHeight
	focalLength := 1.0

	origin := vec.Point3{0.0, 0.0, 0.0}
	horizontal := vec.Vec3{viewportWidth, 0.0, 0.0}
	vertical := vec.Vec3{0, viewportHeight, 0}
	lowerLeftCorner := origin.Sub(horizontal.Scale(0.5)).Sub(vertical.Scale(0.5)).Sub(vec.Vec3{0, 0, focalLength})

	return &Camera{
		aspectRatio:     aspectRatio,
		origin:          origin,
		lowerLeftCorner: lowerLeftCorner,
		horizontal:      horizontal,
		vertical:        vertical,
	}
}

func (camera *Camera) GetRay(u, v float64) ray.Ray {
	return ray.Ray{
		Origin:    camera.origin,
		Direction: camera.lowerLeftCorner.Add(camera.horizontal.Scale(u)).Add(camera.vertical.Scale(v)).Sub(camera.origin),
	}
}

func (camera *Camera) CreateImage(imageWidth int) *image.Image {
	imageHeight := int(float64(imageWidth) / camera.aspectRatio)
	img := image.New(imageWidth, imageHeight)
	return &img
}
