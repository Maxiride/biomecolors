package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

// biomeList will contain a slice of all the biomes found by getBiomes
var biomeList []string

/*
temperatureGroups will contain biomes grouped by climate zone, I will use shades of those colors for the palette.
Vanilla reference:
Biomes are separated into 5 categories. The snow-covered biomes are marked in blue, cold in green,
medium/lush in orange and dry/warm in red.
*/
var temperatureGroups = map[string][]string{
	"Snowy":  {""},
	"Cold":   {""},
	"Medium": {""},
	"Dry":    {""},
}

var usuedColors [][]float64

func main() {
	// Prints the number of biomes found, if true it will print also a list of their names
	totalBiomes := getBiomes(false)
	fmt.Println("Biomes found: ", totalBiomes)

	// Assign each biome to a temperatureGroups key
	assignBiomes()

	// Just for debug to ensure that the numbers add up
	for key := range temperatureGroups {
		fmt.Println(key, len(temperatureGroups[key]))
	}

	prepPalette()
}

// getBiomes will parse the WorldBiomes folder, any biome found will be appended to the biomeList. Returns the total number of biomes found
func getBiomes(test bool) int {
	files, err := ioutil.ReadDir("./WorldBiomes")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		// Add the biomes found to the BiomeList slice
		biomeList = append(biomeList, f.Name())
	}
	// If test == true, prints out a list of the biomes found.
	if test {
		fmt.Println(biomeList)
	}

	return len(files)
}

// assignBiomes will parse through the biomes found in biomeList, read their temperature via biomeTemp and assign them to a temperatureGroups map key
func assignBiomes() {
	for _, biome := range biomeList {
		temp := biomeTemp(biome)

		switch {
		case temp < 0.2:
			temperatureGroups["Snowy"] = append(temperatureGroups["Snowy"], biome)
		case temp < 0.5:
			temperatureGroups["Cold"] = append(temperatureGroups["Cold"], biome)
		case temp < 0.95:
			temperatureGroups["Medium"] = append(temperatureGroups["Medium"], biome)
		case temp <= 2:
			temperatureGroups["Dry"] = append(temperatureGroups["Dry"], biome)

		}
	}

}

// biomeTemp will open each biome config file to retrieve and return its temperature value
func biomeTemp(s string) float64 {
	input, _ := ioutil.ReadFile("./WorldBiomes/" + s)
	lines := strings.Split(string(input), "\n")
	reTemp := regexp.MustCompile(`\d\.?\d*`)

	var temp float64 // create variable so that it is visible outside the for scope
	for _, line := range lines {
		if strings.HasPrefix(line, "BiomeTemperature:") {
			temp, _ = strconv.ParseFloat(reTemp.FindString(line), 64)

		}
	}
	return temp
}

// prepPalette will generate a colour palette for each group in temperatureGroups by calling the respective color settings in
// isSnowy, isCold etc
func prepPalette() {

	usuedColors = [][]float64{}
	snowyColors := len(temperatureGroups["Snowy"])
	snowyPalette, err := colorful.SoftPaletteEx(snowyColors, colorful.SoftPaletteSettings{isSnowy, 50, true})
	if err != nil {
		fmt.Println(err)
	}

	usuedColors = [][]float64{}
	coldColors := len(temperatureGroups["Cold"])
	coldPalette, err := colorful.SoftPaletteEx(coldColors, colorful.SoftPaletteSettings{isCold, 50, true})
	if err != nil {
		fmt.Println(err)
	}

	usuedColors = [][]float64{}
	mediumColors := len(temperatureGroups["Medium"])
	mediumPalette, err := colorful.SoftPaletteEx(mediumColors, colorful.SoftPaletteSettings{isMedium, 50, true})
	if err != nil {
		fmt.Println(err)
	}

	usuedColors = [][]float64{}
	dryColors := len(temperatureGroups["Dry"])
	dryPalette, err := colorful.SoftPaletteEx(dryColors, colorful.SoftPaletteSettings{isDry, 50, true})
	if err != nil {
		fmt.Println(err)
	}

	palettes := map[string][]colorful.Color{
		"Snowy":  snowyPalette,
		"Cold":   coldPalette,
		"Medium": mediumPalette,
		"Dry":    dryPalette,
	}

	for i := range palettes {
		printPalette(palettes[i], len(temperatureGroups[i]), i)
	}

}

// isUnique gets called by each color setting (isSnowy etc) to check that the color hasn't be already used
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

func isCold(l, a, b float64) bool {
	h, c, L := colorful.LabToHcl(l, a, b)

	// ok is a placeholder for Go visibility and scope rules
	ok := false

	// if the the generated color falls into the criteria, check if it has already been generated
	if 295.0 < h && h < 360.0 && 0.1 < c && c < 0.5 && L < 0.76 {
		if unique := isUnique(h, c, L); unique {
			usuedColors = append(usuedColors, []float64{h, c, L})
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
			usuedColors = append(usuedColors, []float64{h, c, L})
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
			usuedColors = append(usuedColors, []float64{h, c, L})
			ok = true
		}
	}

	return ok
}

// printPalette will take in a palette (created in prepPalette), how many colors it contains and its category name (Snowy, Cold etc from temperatureGroups)
func printPalette(currentPalette []colorful.Color, colors int, group string) {

	// Palette squares side lenght and spacing between each of them
	blockw := 40
	space := 5

	// Create image to fit 10 squares per row and as many rows as neeeded to fit all the colors
	img := image.NewRGBA(image.Rect(0, 0, 10*(blockw+space), colors/10*(blockw+space)))

	// Draw on the image
	i := 0
	for row := 0; row < colors/10; row++ {
		for col := 0; col < 10; col++ {
			draw.Draw(img, image.Rect(col*(blockw+space), row*(blockw+space), col*(blockw+space)+blockw, row*(blockw+space)+blockw), &image.Uniform{currentPalette[i]}, image.ZP, draw.Src)
			i++
		}
	}

	// Make image file
	toimg, err := os.Create("./palettes/" + group + "-palette.png")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer toimg.Close()

	png.Encode(toimg, img)
}
