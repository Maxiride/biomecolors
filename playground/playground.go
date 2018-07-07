package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"github.com/lucasb-eyer/go-colorful"
	//"math"
)

var usuedColors [][]float64

func main() {
	colors := 200
	snowyPalette, err := colorful.SoftPaletteEx(colors, colorful.SoftPaletteSettings{isSnowy, 50, true})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(len(snowyPalette))

	blockw := 40
	space := 5
	img := image.NewRGBA(image.Rect(0, 0, 10*(blockw+space), colors/10*(blockw+space)))

	i := 0
	for row := 0; row < colors/10; row++ {
		for col := 0; col < 10; col++ {
			draw.Draw(img, image.Rect(col*(blockw+space), row*(blockw+space), col*(blockw+space)+blockw, row*(blockw+space)+blockw), &image.Uniform{snowyPalette[i]}, image.ZP, draw.Src)
			i++
		}
	}

	toimg, err := os.Create("palettegens.png")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer toimg.Close()

	png.Encode(toimg, img)
}

func isUnique(h, c, L float64) bool {
	// Cycle through the usuedColors slice, if no match is found, returns true and it will be added to the slice
	for _, color := range usuedColors {
		if color[0] == h && color[1] == c && color[2] == L {
			//fmt.Println("This color already exist!")
			//fmt.Println(usuedColors)
			return false
		}
	}
	return true
}

func isSnowy(l, a, b float64) bool {
	h, c, L := colorful.LabToHcl(l, a, b)

	// ok is a placeholder for Go visibility and scope rules
	ok := false

	// if the the generated color falls into the criteria, check if it has already been generated
	if 180.0 < h && h < 280.0 && 0.1 < c && c < 0.5 && L < 0.76 {
		if unique := isUnique(h, c, L); unique {
			usuedColors = append(usuedColors, []float64{h, c, L})
			ok = true
		}
	}

	return ok
}
