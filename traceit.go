package main

import (
	"fmt"
	"math"
	"math/rand"
	"raytracing/camera"
	"raytracing/hittable"
	"raytracing/ray"
	"raytracing/rnd"
	"raytracing/vec"
)

const SHOW_NORMALS = false

const SAMPLES_PER_PIXEL = 50
const MAX_BOUNCES = 25

func main() {
	rnd.Init()
	// img := image.Gradient(1280, 720)
	// img.WriteTarga("gradient.tga")

	// Camera
	cam := camera.Create()

	// Image
	img := cam.CreateImage(400)

	var objects = new(hittable.HittableList)
	objects.Add(&hittable.Sphere{Center: vec.Point3{0.0, 0.0, -1.0}, Radius: 0.5})
	objects.Add(&hittable.Sphere{Center: vec.Point3{0.0, -100.5, -1.0}, Radius: 100.})

	// Render

	for j := img.Height - 1; j >= 0; j-- {
		fmt.Println("Lines Remaining:", j)
		for i := 0; i < img.Width; i++ {
			// u := float64(i) / float64(img.Width-1)
			// v := float64(j) / float64(img.Height-1)
			// r := cam.GetRay(u, v)
			pixelColor := vec.Color{0.0, 0.0, 0.0}
			for s := 0; s < SAMPLES_PER_PIXEL; s++ {
				u := (float64(i) + rand.Float64()) / float64(img.Width-1)
				v := (float64(j) + rand.Float64()) / float64(img.Height-1)
				r := cam.GetRay(u, v)
				sampleContribution := ray_color(&r, objects, MAX_BOUNCES)
				pixelColor = pixelColor.Add(sampleContribution)
			}

			img.SetPixel(i, j, colorCorrect(pixelColor, SAMPLES_PER_PIXEL))
		}
	}

	img.WriteTarga("output.tga")
}

func ray_color(r *ray.Ray, world hittable.Hittable, bouncesLeft int) vec.Color {
	if bouncesLeft <= 0 {
		return vec.Color{0, 0, 0}
	}

	hit, isHit := world.Hit(r, 0.001, 10000000.0)
	if isHit {
		if SHOW_NORMALS {
			return vec.Color{hit.Normal.X + 1, hit.Normal.Y + 1, hit.Normal.Z + 1}.Scale(0.5)
		} else {
			// // target := hit.P.Add(hit.Normal).Add(vec.RandomUnitSphere())
			// // dir := target.Sub(hit.P)
			dir := vec.RandomUnitSphere()
			if vec.Dot(hit.Normal, dir) < 0 {
				dir = dir.Neg()
			}
			return ray_color(&ray.Ray{hit.P, dir}, world, bouncesLeft-1).Scale(0.5)
			// // return vec.Color{1.0, 1.0, 1.0}.Add(dir).Scale(0.5)
		}
	}

	unit_direction := r.Direction.Norm()
	t := 0.5 * (unit_direction.Y + 1.0)
	return vec.Color{1.0, 1.0, 1.0}.Scale(1.0 - t).Add(vec.Color{0.5, 0.7, 1.0}.Scale(t))
}

func colorCorrect(color vec.Color, samplesPerPixel int) vec.Color {
	// Sample Scaling
	scale := 1.0 / float64(samplesPerPixel)
	color = color.Scale(scale)

	// Gamma 2.0 Coorrection
	color.X = math.Sqrt(color.X)
	color.Y = math.Sqrt(color.Y)
	color.Z = math.Sqrt(color.Z)

	// Color Space Clamping
	color.X = clamp(color.X, 0.0, 0.999999)
	color.Y = clamp(color.Y, 0.0, 0.999999)
	color.Z = clamp(color.Z, 0.0, 0.999999)

	return color
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	} else if v > max {
		return max
	} else {
		return v
	}
}
