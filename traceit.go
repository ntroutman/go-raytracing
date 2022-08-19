package main

import (
	"fmt"
	"math"
	"math/rand"
	"raytracing/camera"
	"raytracing/hittable"
	"raytracing/img"
	"raytracing/mat"
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
	renderer.img = renderer.cam.CreateImage(100)

	var objects = new(hittable.HittableList)

	materialGround := &mat.Lambertian{vec.Color{0.8, 0.8, 0.0}}
	materialCenter := &mat.Lambertian{vec.Color{0.7, 0.3, 0.3}}
	materialLeft := &mat.Metal{vec.Color{0.8, 0.8, 0.8}, 0.3}
	materialRight := &mat.Metal{vec.Color{0.8, 0.6, 0.2}, 1.0}

	objects.Add(&hittable.Sphere{Center: vec.Point3{0.0, -100.5, -1.0}, Radius: 100., Mat: materialGround})
	objects.Add(&hittable.Sphere{Center: vec.Point3{0.0, 0.0, -1.0}, Radius: 0.5, Mat: materialCenter})
	objects.Add(&hittable.Sphere{Center: vec.Point3{-1.0, 0.0, -1.0}, Radius: 0.5, Mat: materialLeft})
	objects.Add(&hittable.Sphere{Center: vec.Point3{1.0, 0.0, -1.0}, Radius: 0.5, Mat: materialRight})

	renderer.world = objects

	// Render
	//renderer.RenderSingleThread()
	renderer.RenderMultiThreadYBatch()
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
	lineChannel := make(chan *renderedLine, 500)
	for y := r.img.Height - 1; y >= 0; y-- {
		go r.renderSingleLine(y, lineChannel)
	}

	for y := r.img.Height - 1; y >= 0; y-- {
		line := <-lineChannel
		fmt.Println("Lines Remaining:", y)
		copy(r.img.Pixels[y*3:(y+1)*3], line.pixels)
	}

	r.img.WriteTarga("output.tga")
}

func (r *Renderer) renderSingleLine(y int, lineChannel chan *renderedLine) {
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

func (r *Renderer) RenderMultiThreadYBatch() {
	lineChannel := make(chan *renderedLine, 500)
	const BATCH_SIZE = 25
	for y := r.img.Height - 1; y >= 0; y -= BATCH_SIZE {
		yEnd := y - BATCH_SIZE
		if yEnd < 0 {
			yEnd = 0
		}
		fmt.Println("Batch:", y, yEnd)
		go r.renderLineBatch(y, yEnd, lineChannel)
	}

	for y := r.img.Height - 1; y >= 0; y -= BATCH_SIZE {
		line := <-lineChannel
		fmt.Println("Lines Remaining:", y)
		copy(r.img.Pixels[y*3:(y+BATCH_SIZE)*3], line.pixels)
	}

	r.img.WriteTarga("output.tga")
}

func (r *Renderer) renderLineBatch(yStart, yEnd int, lineChannel chan *renderedLine) {
	width := r.img.Width

	out := new(renderedLine)
	out.y = yStart
	out.pixels = make([]byte, width*3*(yStart-yEnd+1))
	idx := 0
	for y := yStart; y >= yEnd; y-- {
		fmt.Println("Line:", y, yStart, yEnd)
		for x := 0; x < width; x++ {
			pixelColor := r.renderPixel(x, y)

			out.pixels[idx] = byte(256 * pixelColor.Z)
			idx += 1
			out.pixels[idx] = byte(256 * pixelColor.X)
			idx += 1
			out.pixels[idx] = byte(256 * pixelColor.Y)
			idx += 1
		}
	}
	lineChannel <- out
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
			// target := hit.P.Add(vec.RandomInHemiSphere(hit.Normal))
			// return calcRayColor(&ray.Ray{hit.P, target.Sub(hit.P)}, world, bouncesLeft-1).Scale(.5)
			scatter, didScatter := hit.Scatter(r)
			if didScatter {
				return calcRayColor(&scatter.Scattered, world, bouncesLeft-1).Mul(scatter.Attenuation)
			} else {
				return vec.Color{0, 0, 0}
			}
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
