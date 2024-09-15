package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"slices"
	"strings"

	"github.com/crazy3lf/colorconv"
)

type ConfigOptions struct {
	dither bool
	invert bool
}

type pixel struct {
	r, g, b, a uint8
	x, y       int
}

const CHARS string = " .:-=+#%█"           // 9 chars
const DITHER_CHARS = " .,\"`:-=+^~*;#%▒▓█" // 18 chars

func main() {
	d := flag.Bool("d", false, "Apply dithering to image")
	inv := flag.Bool("i", false, "Invert colors")
	flag.Parse()

	if flag.NArg() != 1 {
		printUsage()
		return
	}

	img := openImage(flag.Args()[0])
	lightnessGrid := getLightnessGrid(img)

	var outPath string
	if *d == true && *inv == true {
		outPath = "ascii_" + strings.Split(flag.Arg(0), ".")[0] + "_inverted_dithered.txt"
	}
	if *d == true && *inv != true {
		outPath = "ascii_" + strings.Split(flag.Arg(0), ".")[0] + "_dithered.txt"
	}
	if *inv == true && *d != true {
		outPath = "ascii_" + strings.Split(flag.Arg(0), ".")[0] + "_inverted.txt"
	} else if *inv != true && *d != true {
		outPath = "ascii_" + strings.Split(flag.Arg(0), ".")[0] + ".txt"
	}

	out, err := os.Create(outPath)
	if err != nil {
		fmt.Println("error creating file:", err)
		os.Exit(1)
	}

	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	asciiGrid := make([][]string, width)
	for x := range asciiGrid {
		asciiGrid[x] = make([]string, height)
	}

	for y := 0; y < height; y++ {
		line := ""
		for x := 0; x < width; x++ {
			var char string
			if *d == true {
				char = getDitheredAscii(lightnessGrid[x][y], *inv)
			} else {
				char = getAscii(lightnessGrid[x][y], *inv)
			}
			line += char
		}
		out.WriteString(line + "\n")
	}

	fmt.Printf("saved ascii art to ./%s\n", outPath)
	defer out.Close()
}

// ------------------------
// helper functions
// ------------------------

func getDitheredAscii(lightness float64, invert bool) string {
	c := strings.Split(DITHER_CHARS, "")
	if invert != true {
		slices.Reverse(c)
	}
	if lightness <= 0.055 {
		return c[0]
	}
	if lightness <= 0.111 {
		return c[1]
	}
	if lightness <= 0.166 {
		return c[2]
	}
	if lightness <= 0.222 {
		return c[3]
	}
	if lightness <= 0.277 {
		return c[4]
	}
	if lightness <= 0.333 {
		return c[5]
	}
	if lightness <= 0.388 {
		return c[6]
	}
	if lightness <= 0.444 {
		return c[7]
	}
	if lightness <= 0.499 {
		return c[8]
	}
	if lightness <= 0.555 {
		return c[9]
	}
	if lightness <= 0.611 {
		return c[10]
	}
	if lightness <= 0.666 {
		return c[11]
	}
	if lightness <= 0.722 {
		return c[12]
	}
	if lightness <= 0.777 {
		return c[13]
	}
	if lightness <= 0.833 {
		return c[14]
	}
	if lightness <= 0.888 {
		return c[15]
	}
	if lightness <= 0.944 {
		return c[17]
	} else {
		return c[17]
	}
}

func getAscii(lightness float64, invert bool) string {
	c := strings.Split(CHARS, "")
	if invert != true {
		slices.Reverse(c)
	}
	if lightness <= 0.111 {
		return c[0]
	}
	if lightness <= 0.222 {
		return c[1]
	}
	if lightness <= 0.333 {
		return c[2]
	}
	if lightness <= 0.444 {
		return c[3]
	}
	if lightness <= 0.555 {
		return c[4]
	}
	if lightness <= 0.666 {
		return c[5]
	}
	if lightness <= 0.777 {
		return c[6]
	}
	if lightness <= 0.888 {
		return c[7]
	} else {
		return c[8]
	}
}

func printUsage() {
	fmt.Println("usage: img2ascii [-d] [-i] <path/to/image>")
	fmt.Println("  -d apply dithering to image")
	fmt.Println("  -i invert colors")
}

func openImage(path string) image.Image {
	f, fErr := os.Open(path)
	if fErr != nil {
		fmt.Println("err: file could not be opened")
		fmt.Println(fErr)
		os.Exit(1)
	}
	defer f.Close()

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)

	img, _, imgErr := image.Decode(f)
	if imgErr != nil {
		fmt.Println("err: file could not be decoded")
		fmt.Println(imgErr)
		os.Exit(1)
	}
	return img
}

func getLightnessGrid(img image.Image) [][]float64 {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	p := make([][]float64, width)
	for x := range p {
		p[x] = make([]float64, height)
		for y := 0; y < height; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			p[x][y] = getLightness(
				uint8(r>>8),
				uint8(g>>8),
				uint8(b>>8),
			)
		}
	}
	return p
}

func getLightness(r, g, b uint8) float64 {
	_, _, l := colorconv.RGBToHSL(r, g, b)
	return l
}
