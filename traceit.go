package main

import (
	"fmt"
	"math"
	"math/rand"
	"raytracing/camera"
	"raytracing/hittable"
	"raytracing/img"
	"raytracing/ray"
	"raytracing/rnd"
	"raytracing/vec"
)

const SHOW_NORMALS = false

const SAMPLES_PER_PIXEL = 100
const MAX_BOUNCES = 50

func main() {
	rnd.Init()
	// img := image.Gradient(1280, 720)
	// img.WriteTarga("gradient.tga")

	var renderer Renderer
	renderer.cam = camera.Create()
	renderer.img = renderer.cam.CreateImage(400)

	var objects = new(hittable.HittableList)
	objects.Add(&hittable.Sphere{Center: vec.Point3{0.0, 0.0, -1.0}, Radius: 0.5})
	objects.Add(&hittable.Sphere{Center: vec.Point3{0.0, -100.5, -1.0}, Radius: 100.})
	renderer.world = objects

	// Render
	renderer.RenderSingleThread()

}

type Renderer struct {
	cam   *camera.Camera
	img   *img.Image
	world *hittable.HittableList
}

func (r *Renderer) RenderSingleThread() {
	for y := r.img.Height - 1; y >= 0; y-- {
		fmt.Println("Lines Remaining:", y)
		for x := 0; x < r.img.Width; x++ {
			color := r.renderPixel(x, y)
			r.img.SetPixel(x, y, color)
		}
	}

	r.img.WriteTarga("output.tga")
}

func (r *Renderer) RenderMultiThread() {
	lineChannel := make(chan *renderedLine, 50)
	for y := r.img.Height - 1; y >= 0; y-- {
		go r.renderLine(y, lineChannel)
	}

	for y := r.img.Height - 1; y >= 0; y-- {
		line := <-lineChannel
		fmt.Println("Lines Remaining:", y)
		for x := 0; x < r.img.Width; x++ {
			//r.img.SetPixel(x, line.y, line.pixels[x])
			copy(r.img.Pixels[y*3:(y+1)*3], line.pixels)

		}
	}

	r.img.WriteTarga("output.tga")
}

func (r *Renderer) renderLine(y int, lineChannel chan *renderedLine) {
	out := new(renderedLine)
	out.y = y
	out.pixels = make([]byte, r.img.Width*3)
	for x := 0; x < r.img.Width; x++ {
		pixelColor := r.renderPixel(x, y)
		out.pixels[x*3+2] = byte(256 * pixelColor.X)
		out.pixels[x*3+1] = byte(256 * pixelColor.Y)
		out.pixels[x*3+0] = byte(256 * pixelColor.Z)
	}
	lineChannel <- out
}

type renderedLine struct {
	y      int
	pixels []byte
}

func (r *Renderer) renderPixel(x, y int) vec.Color {
	// u := float64(i) / float64(img.Width-1)
	// v := float64(j) / float64(img.Height-1)
	// r := cam.GetRay(u, v)
	pixelColor := vec.Color{0.0, 0.0, 0.0}
	for s := 0; s < SAMPLES_PER_PIXEL; s++ {
		u := (float64(x) + rand.Float64()) / float64(r.img.Width-1)
		v := (float64(y) + rand.Float64()) / float64(r.img.Height-1)
		ray := r.cam.GetRay(u, v)
		sampleContribution := calcRayColor(&ray, r.world, MAX_BOUNCES)
		pixelColor = pixelColor.Add(sampleContribution)
	}

	return colorCorrect(pixelColor, SAMPLES_PER_PIXEL)
}

func calcRayColor(r *ray.Ray, world hittable.Hittable, bouncesLeft int) vec.Color {
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

			target := hit.P.Add(vec.RandomInHemiSphere(hit.Normal))
			return calcRayColor(&ray.Ray{hit.P, target.Sub(hit.P)}, world, bouncesLeft-1).Scale(.5)

			// dir := vec.RandomUnitVector()
			// return calcRayColor(&ray.Ray{hit.P, dir}, world, bouncesLeft-1).Scale(0.5)
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
