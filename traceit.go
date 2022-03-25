package main

import (
	"fmt"
	"raytracing/hittable"
	"raytracing/image"
	"raytracing/ray"
	"raytracing/vec"
)

func main() {
	// img := image.Gradient(1280, 720)
	// img.WriteTarga("gradient.tga")

	// Image
	aspect_ratio := 16.0 / 9.0
	image_width := 400
	image_height := int(float64(image_width) / aspect_ratio)
	img := image.New(image_width, image_height)

	// Camera

	viewport_height := 2.0
	viewport_width := aspect_ratio * viewport_height
	focal_length := 1.0

	origin := vec.Point3{0.0, 0.0, 0.0}
	horizontal := vec.Vec3{viewport_width, 0.0, 0.0}
	vertical := vec.Vec3{0, viewport_height, 0}
	lower_left_corner := origin.Sub(horizontal.Scale(0.5)).Sub(vertical.Scale(0.5)).Sub(vec.Vec3{0, 0, focal_length})
	fmt.Printf("lower_left_corner: %v\n", lower_left_corner)

	var objects = new(hittable.HittableList)
	objects.Add(&hittable.Sphere{Center: vec.Point3{0.0, 0.0, -1.0}, Radius: 0.5})
	objects.Add(&hittable.Sphere{Center: vec.Point3{0.0, -100.5, -1.0}, Radius: 100.})

	// Render

	for j := image_height - 1; j >= 0; j-- {
		for i := 0; i < image_width; i++ {
			u := float64(i) / float64(image_width-1)
			v := float64(j) / float64(image_height-1)
			r := ray.Ray{origin, lower_left_corner.Add(horizontal.Scale(u)).Add(vertical.Scale(v)).Sub(origin)}
			// fmt.Printf("r: %v\n", r)
			pixel_color := ray_color(&r, objects)
			//fmt.Printf("pixel_color: %v\n", pixel_color)
			img.SetPixel(i, j, pixel_color)
		}
	}

	img.WriteTarga("output.tga")
}

func ray_color(r *ray.Ray, world hittable.Hittable) vec.Color {
	hit, isHit := world.Hit(r, 0, 5)
	if isHit {
		return vec.Color{hit.Normal.X + 1, hit.Normal.Y + 1, hit.Normal.Z + 1}.Scale(0.5)
	}

	unit_direction := r.Direction.Norm()
	t := 0.5 * (unit_direction.Y + 1.0)
	return vec.Color{1.0, 1.0, 1.0}.Scale(1.0 - t).Add(vec.Color{0.5, 0.7, 1.0}.Scale(t))
}
