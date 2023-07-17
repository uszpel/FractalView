package mandelbrot

import (
	"image"
	"image/color"
	"math/cmplx"
	"sync"
)

type ImgConfig struct {
	CenterX float64
	CenterY float64
	ZoomX   float64
	ZoomY   float64
}

var width, height int = 1920, 1080
var aspectRatio = float64(width) / float64(height)

func NewConfig() *ImgConfig {
	zoom := 0.5
	return &ImgConfig{
		CenterX: -0.8,
		CenterY: 0.0,
		ZoomX:   zoom,
		ZoomY:   zoom / aspectRatio,
	}
}

func (imgConfig *ImgConfig) Move(direction string) {
	if direction == "left" {
		imgConfig.CenterX = imgConfig.CenterX - 0.03*imgConfig.ZoomX/0.5
	} else if direction == "right" {
		imgConfig.CenterX = imgConfig.CenterX + 0.03*imgConfig.ZoomX/0.5
	} else if direction == "up" {
		imgConfig.CenterY = imgConfig.CenterY - 0.03*imgConfig.ZoomY/(0.5/aspectRatio)
	} else if direction == "down" {
		imgConfig.CenterY = imgConfig.CenterY + 0.03*imgConfig.ZoomY/(0.5/aspectRatio)
	}
}

func (imgConfig *ImgConfig) Scale(direction string) {
	if direction == "plus" {
		imgConfig.ZoomX = imgConfig.ZoomX / 1.5
	} else if direction == "minus" {
		imgConfig.ZoomX = imgConfig.ZoomX * 1.5
	}
	imgConfig.ZoomY = imgConfig.ZoomX / aspectRatio
}

func CreateImage(config *ImgConfig) *image.RGBA {
	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: width, Y: height},
	})

	float_width := float64(width)
	float_height := float64(height)
	var wg sync.WaitGroup
	wg.Add(height)
	for y := 0; y < height; y++ {
		go func(y int) {
			float_y := float64(y)
			for x := 0; x < width; x++ {
				result, i := check(toCoord(float64(x), float_y, float_width, float_height, config))
				img.Set(x, y, chooseColor(i, result))
			}
			wg.Done()
		}(y)
	}
	wg.Wait()
	return img
}

func check(c complex128) (bool, int) {
	var z = c
	for i := 0; i < 100; i++ {
		z = cmplx.Pow(z, 2) + c
		rZ := real(z)
		iZ := imag(z)
		if rZ*rZ+iZ*iZ > 4 {
			return false, i
		}
	}
	return true, 0
}

func toCoord(x float64, y float64, width float64, height float64, config *ImgConfig) complex128 {
	return complex(((x*config.ZoomX/width)-config.ZoomX/2)+(config.CenterX*aspectRatio),
		((y*config.ZoomY/height)-config.ZoomY/2)+(config.CenterY))
}

func chooseColor(i int, result bool) color.RGBA {
	if result {
		return color.RGBA{0, 0, 0, 255}
	} else {
		var set = i * 0xFFFFFF / 100
		return color.RGBA{uint8((set & 0xFF0000) >> 16), uint8((set & 0x00FF00) >> 8), uint8(set & 0x0000FF), 255}
	}
}
