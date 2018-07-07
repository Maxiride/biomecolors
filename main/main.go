package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// BiomeList will contain a slice of all the biomes found
var BiomeList []string

/*
temperatureGroups will contain biomes grouped by climate zone
Vanilla reference
Biomes are separated into 5 categories. The snow-covered biomes are marked in blue, cold in green,
medium/lush in orange and dry/warm in red. I will use shades of those colors for the palette.
*/
var temperatureGroups = map[string][]string{
	"Snowy":  {""},
	"Cold":   {""},
	"Medium": {""},
	"Dry":    {""},
}

func main() {
	// Prints the number of biomes found, if true it will print also a list of their names
	totalBiomes := getBiomes(false)
	fmt.Println("Biomes found: ", totalBiomes)

	// Assign each biome to a temperature group
	assignBiomes()

	for key := range temperatureGroups{
		fmt.Println(key, len(temperatureGroups[key]))
	}

}

// Get Biome list, append them to the BiomeList and return the number of biomes found
func getBiomes(test bool) int {
	files, err := ioutil.ReadDir("./WorldBiomes")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		// Add the biomes found to the BiomeList slice
		BiomeList = append(BiomeList, f.Name())
	}
	// If test == true, prints out a list of the biomes found.
	if test {
		fmt.Println(BiomeList)
	}

	return len(files)
}

// assignBiomes will parse through the biomes found, read their temperature via biomeTemp and assign them to a TemperatureGroups map key
func assignBiomes() {
	for _, biome := range BiomeList {
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

// biomeTemp will opne each biome config file to retrieve its temperature value
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
