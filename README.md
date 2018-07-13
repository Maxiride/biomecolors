# What's this about

This program is born from the necessity to manage a lot of biomes and their associated colours either when assembling a new World Preset or when making a FromImage Map. It will scan all your biomes and assign them unique colours, cheat sheet reference palettes will also be generated.

## How does it work

In this early stage of the code the program will:
  1. Scan the WorldBiomes folder to fetch all the present biomes
  2. Categorize all the biomes found in 4 climate zones (Snowy, Cold, Medium, Dry) according to the [vanilla standards](https://minecraft.gamepedia.com/Biome#Biome_types).
  3. Generate [4 different palettes](https://github.com/Maxiride/biomecolors/tree/develop/main/palettes) with colours ranging from "coldish" to "dryish" gradients with each colour being enough "distant" from all the others in order to be easily perceived by a human as a different colour according to the CIELab high quality standards.
  4. [Soon to be implemented] Write the generated colours in the biomes *.bc configuration files.
  
  ### Future roadmap
 
 Ideas that need further investigation:
  - Generate ready to use .ase preset for Photoshop with the generated palettes
  - Check if biomes are Isle or Border and make a new palette where to group them for quick reference. This will insanely help in drawing smooth FromImage maps with technical biomes as the border ones.
  - As of now the program will overwrite any previously set BiomeColor, enhancing the code to take in consideration "legacy" biomes to be added to the palettes without being overwritten with a new colour.
  
  
  ## Contributing
   - Submit ideas and I will consider them =)
   - Found a bug? Report it ASAP!
   - enjoy =)
