package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful" // Colors management
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"gopkg.in/ini.v1" // INI parser
)

// biomeList will contain a map of type biome of all the biomes found by getBiomes indexed with the biome name as the key
var biomeList = make(map[string][]biome)

// same as above but this map has keys set as climate type
var biomeListByClimate = make(map[string][]biome)

/*
The biome struct will contain biomes infos.
Vanilla reference for the climateZone:
Biomes are separated into 5 categories.
The snow-covered biomes are marked in blue with temperature <0.2
cold in green  with temperature <0.5,
medium/lush in orange  with temperature <0.95
and dry/warm in red  with temperature <=2.
*/

type biome struct {
	name        string
	temperature float64
	climateZone string
}

// Used late in prepPalette to check if a color has already been used
var usedColors [][]float64

func main() {

	totalBiomes := getBiomes()
	fmt.Println("Biomes found: ", totalBiomes)

	prepPalette()
}

// getBiomes will parse the WorldBiomes folder, any biome found will be appended to the biomeList. Returns the total number of biomes found
func getBiomes() int {
	files, err := ioutil.ReadDir("./WorldBiomes")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		biomeCategorization(f.Name())
	}

	return len(files)
}

// biomeCategorization will parse every biome and create its struct with proper fields
func biomeCategorization(s string) {

	// get a biome file
	input, err := ini.LoadSources(ini.LoadOptions{
		SkipUnrecognizableLines: true,
	}, "./WorldBiomes/"+s)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// Get its temperature value from the key BiomeTemperature
	temp, _ := strconv.ParseFloat(input.Section("").Key("BiomeTemperature").String(), 64)

	// Let's see which climate this biome has
	climate := ""
	switch {
	case temp < 0.2:
		climate = "Snowy"
	case temp < 0.5:
		climate = "Cold"
	case temp < 0.95:
		climate = "Medium"
	case temp <= 2:
		climate = "Dry"

	}

	// Save the found biomes in their maps
	biomeList[s] = append(biomeList[s], biome{
		name:        s,
		temperature: temp,
		climateZone: climate,
	})

	biomeListByClimate[climate] = append(biomeListByClimate[climate], biome{
		name:        s,
		temperature: temp,
		climateZone: climate,
	})

}

/*
prepPalette will generate a colour palette for each climate in biomeListByClimate ensuring that all the generated colors falls under the right range with the isSnowy, isCold etc functions.

*/
func prepPalette() {

	usedColors = [][]float64{}
	snowyColors := len(biomeListByClimate["Snowy"])
	snowyPalette, err := colorful.SoftPaletteEx(snowyColors, colorful.SoftPaletteSettings{isSnowy, 50, true})
	if err != nil {
		fmt.Println(err)
	}

	usedColors = [][]float64{}
	coldColors := len(biomeListByClimate["Cold"])
	coldPalette, err := colorful.SoftPaletteEx(coldColors, colorful.SoftPaletteSettings{isCold, 50, true})
	if err != nil {
		fmt.Println(err)
	}

	usedColors = [][]float64{}
	mediumColors := len(biomeListByClimate["Medium"])
	mediumPalette, err := colorful.SoftPaletteEx(mediumColors, colorful.SoftPaletteSettings{isMedium, 50, true})
	if err != nil {
		fmt.Println(err)
	}

	usedColors = [][]float64{}
	dryColors := len(biomeListByClimate["Dry"])
	dryPalette, err := colorful.SoftPaletteEx(dryColors, colorful.SoftPaletteSettings{isDry, 50, true})
	if err != nil {
		fmt.Println(err)
	}

	// Let's organize the palettes
	palettes := map[string][]colorful.Color{
		"Snowy":  snowyPalette,
		"Cold":   coldPalette,
		"Medium": mediumPalette,
		"Dry":    dryPalette,
	}

	// Draw them!
	for i := range palettes {
		drawPalette(palettes[i], len(biomeListByClimate[i]), i)
	}

}

// isUnique gets called by each color setting (isSnowy etc) to check that the color hasn't be already used
func isUnique(h, c, L float64) bool {
	// Cycle through the usuedColors slice, if no match is found, returns true and it will be added to the slice
	for _, color := range usedColors {
		if color[0] == h && color[1] == c && color[2] == L {
			//fmt.Println("This color already exist!")
			//fmt.Println(usedColors)
			return false
		}
	}
	return true
}

// All the color checking functions, I guess all of them could be grouped as a interface? And then prepPalette just call check.isSnowy, check.isCold mmm
func isSnowy(l, a, b float64) bool {
	h, c, L := colorful.LabToHcl(l, a, b)

	// ok is a placeholder for Go visibility and scope rules, do I really need this?? Must be a smarter way
	ok := false

	// if the the generated color falls into the criteria, check if it has already been generated
	if 180.0 < h && h < 280.0 && 0.1 < c && c < 0.5 && L < 0.76 {
		if unique := isUnique(h, c, L); unique {
			usedColors = append(usedColors, []float64{h, c, L})
			ok = true
		}
	}

	return ok
}

func isCold(l, a, b float64) bool {
	h, c, L := colorful.LabToHcl(l, a, b)

	// ok is a placeholder for Go visibility and scope rules
	ok := false

	// if the the generated color falls into the criteria, check if it has already been generated
	if 295.0 < h && h < 360.0 && 0.1 < c && c < 0.5 && L < 0.76 {
		if unique := isUnique(h, c, L); unique {
			usedColors = append(usedColors, []float64{h, c, L})
			ok = true
		}
	}

	return ok
}

func isMedium(l, a, b float64) bool {
	h, c, L := colorful.LabToHcl(l, a, b)

	// ok is a placeholder for Go visibility and scope rules
	ok := false

	// if the the generated color falls into the criteria, check if it has already been generated
	if 120.0 < h && h < 160.0 && 0.1 < c && c < 0.5 && L < 0.76 {
		if unique := isUnique(h, c, L); unique {
			usedColors = append(usedColors, []float64{h, c, L})
			ok = true
		}
	}

	return ok
}

func isDry(l, a, b float64) bool {
	h, c, L := colorful.LabToHcl(l, a, b)

	// ok is a placeholder for Go visibility and scope rules
	ok := false

	// if the the generated color falls into the criteria, check if it has already been generated
	if 65.0 < h && h < 110.0 && 0.1 < c && c < 0.5 && L < 0.76 {
		if unique := isUnique(h, c, L); unique {
			usedColors = append(usedColors, []float64{h, c, L})
			ok = true
		}
	}

	return ok
}

// drawPalette will take in a palette (created in prepPalette), how many colors it contains and its climate name (Snowy, Cold etc)
func drawPalette(currentPalette []colorful.Color, colors int, group string) {
	/*
		The image layout is meant to be on two columns "colored square" "HEX Value" "Biome Name"
	*/

	// Palette squares side length and spacing between each of them
	blockw := 40
	space := 5

	// Create an image with width 800px and height as much as needed to fit everything
	img := image.NewRGBA(image.Rect(0, 0, 800, colors/2*(blockw+space)))

	// White background
	draw.Draw(img, image.Rect(0, 0, 800, colors/2*(blockw+space)), image.NewUniform(color.RGBA{255, 255, 255, 255}), image.ZP, draw.Src)

	// Draw on the image each color of the current palette
	i := 0
	for row := 0; row < colors/2; row++ {
		for j := 0; j <= 1; j++ {
			draw.Draw(img, image.Rect(400*j, row*(blockw+space), 400*j+2*blockw, row*(blockw+space)+blockw), &image.Uniform{currentPalette[i]}, image.ZP, draw.Src)
			addLabel(img, 12+400*j, 25+(row*(blockw+space)), strings.ToUpper(currentPalette[i].Hex()))
			addLabel(img, 120+400*j, 25+(row*(blockw+space)), biomeListByClimate[group][i].name)
			i++
		}

	}

	os.Mkdir("./palettes", os.ModePerm)
	// Make image file
	toimg, err := os.Create("./palettes/" + group + "-palette.png")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer toimg.Close()

	png.Encode(toimg, img)
}

// Write down on the image the HEX color and the biome name
func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{127, 255, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
