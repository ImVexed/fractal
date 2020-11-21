package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/valyala/fastrand"
)

const (
	// Position and size
	px   = -0.5557506
	py   = -0.55560
	size = 0.000000001
	//px   = -2
	//py   = -1.2
	//size = 2.5

	// Quality
	imgWidth = 1024 * 10
	maxIter  = 1000
	samples  = 50
)

func main() {
	log.Println("Allocating image...")
	img := image.NewRGBA(image.Rect(0, 0, int(imgWidth), int(imgWidth)))

	log.Println("Rendering...")
	start := time.Now()
	render(img)
	end := time.Now()

	log.Println("Done rendering in", end.Sub(start))

	log.Println("Encoding image...")
	f, err := os.Create("result.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
	log.Println("Done!")
}

func fastRand() float64 {
	return float64(int64(fastrand.Uint32())+int64(fastrand.Uint32())) / (1 << 63)
}
func render(img *image.RGBA) {
	jobs := make(chan int)

	for w := 1; w <= runtime.NumCPU(); w++ {
		go worker(img, jobs)
	}

	for y := 0; y < imgWidth; y++ {
		jobs <- y
		fmt.Printf("\r%d/%d (%d%%)", y, imgWidth, int(100*(float64(y)/imgWidth)))
	}
	fmt.Println()
}

func worker(img *image.RGBA, jobs <-chan int) {
	for y := range jobs {
		for x := 0; x < imgWidth; x++ {
			var r, g, b int

			for i := 0; i < samples; i++ {
				nx := size*((float64(x)+fastRand())/float64(imgWidth)) + px
				ny := size*((float64(y)+fastRand())/float64(imgWidth)) + py
				colour := paint(mandelbrotIter(nx, ny, maxIter))

				r += int(colour.R)
				g += int(colour.G)
				b += int(colour.B)
			}

			img.SetRGBA(int(x), int(y), color.RGBA{
				R: uint8(float64(r) / float64(samples)),
				G: uint8(float64(g) / float64(samples)),
				B: uint8(float64(b) / float64(samples)),
				A: 255,
			})
		}
	}
}

func paint(r float64, n int) color.RGBA {
	var insideSet = color.RGBA{R: 255, G: 255, B: 255, A: 255}

	if r > 4 {
		c := hslToRGB(float64(n)/800*r, 1, 0.5)
		return c
	}

	return insideSet
}

func mandelbrotIter(px, py float64, maxIter int) (float64, int) {
	var x, y, xx, yy, xy float64

	for i := 0; i < maxIter; i++ {
		xx, yy, xy = x*x, y*y, x*y
		if xx+yy > 4 {
			return xx + yy, i
		}
		x = xx - yy + px
		y = 2*xy + py
	}

	return xx + yy, maxIter
}
