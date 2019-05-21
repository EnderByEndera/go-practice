package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

// Putpixel describes a function expected to draw a point on a bitmap at (x, y) coordinates.
type Putpixel func(x, y int)

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

func DrawDot(x, y int, brush Putpixel) {
	for i := x - 2; i <= x+2; i++ {
		for j := y - 2; j <= y+2; j++ {
			brush(i, j)
		}
	}
}

func DrawLine(x0, y0, x1, y1 int, brush Putpixel) { // draw a line in a picture
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx, sy := 1, 1
	if x0 >= x1 {
		sx = -1
	}
	if y0 >= y1 {
		sy = -1
	}
	err := dx - dy

	for {
		brush(x0, y0)
		if x0 == x1 && y0 == y1 {
			return
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func (resources *priores) DrawPic() {
	// Create image of the path
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	// start drawing a picture
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			img.Set(i, j, color.RGBA{A: 240})
		}
	}
	for i := 0; i < resources.N-1; i++ {
		DrawDot(resources.resource[resources.seq.seq[i]].X, resources.resource[resources.seq.seq[i]].Y, func(x, y int) {
			switch resources.resource[resources.seq.seq[i]].Priority {
			case 0:
				img.Set(x, y, color.RGBA{R: 255, A: 255})
			case 1:
				img.Set(x, y, color.RGBA{B: 255, A: 255})
			case 2:
				img.Set(x, y, color.RGBA{G: 255, A: 255})
			case 3:
				img.Set(x, y, color.RGBA{R: 255, B: 255, A: 255})
			}
		})
		DrawLine(resources.resource[resources.seq.seq[i]].X, resources.resource[resources.seq.seq[i]].Y,
			resources.resource[resources.seq.seq[i+1]].X, resources.resource[resources.seq.seq[i+1]].Y, func(x, y int) {
				img.Set(x, y, color.RGBA{})
			})
	}
	DrawDot(resources.resource[resources.seq.seq[0]].X, resources.resource[resources.seq.seq[0]].Y, func(x, y int) {
		img.Set(x, y, color.RGBA{R: 255, G: 255, A: 255})
	})
	if resources.N != 0 {
		DrawDot(resources.resource[resources.seq.seq[resources.N-1]].X, resources.resource[resources.seq.seq[resources.N-1]].Y, func(x, y int) {
			switch resources.resource[resources.seq.seq[resources.N-1]].Priority {
			case 0:
				img.Set(x, y, color.RGBA{R: 255, A: 255})
			case 1:
				img.Set(x, y, color.RGBA{B: 255, A: 255})
			case 2:
				img.Set(x, y, color.RGBA{G: 255, A: 255})
			case 3:
				img.Set(x, y, color.RGBA{R: 255, B: 255, A: 255})
			}
		})
	}

	imgFile, _ := os.Create(fmt.Sprintf("output.png"))
	defer imgFile.Close()
	// Save file using .png format
	err := png.Encode(imgFile, img)
	if err != nil {
		log.Fatal(err)
	}
}
