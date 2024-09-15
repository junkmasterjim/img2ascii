package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/crazy3lf/colorconv"
	"github.com/nfnt/resize"
)

type ConfigOptions struct {
	dither bool
	invert bool
	scale  float64
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

	scale := 0.25 // default scale
	args := flag.Args()
	if len(args) < 1 || len(args) > 2 {
		printUsage()
		return
	}
	imagePath := args[0]

	if len(args) == 2 {
		if s, err := strconv.ParseFloat(args[1], 64); err == nil {
			scale = s
		} else {
			fmt.Println("Invalid scale value. Using default scale:", scale)
		}
	}

	config := ConfigOptions{
		dither: *d,
		invert: *inv,
		scale:  scale,
	}

	img := openImage(imagePath)
	scaledImg := resize.Resize(uint(float64(img.Bounds().Max.X)*config.scale), 0, img, resize.Lanczos3)
	lightnessGrid := getLightnessGrid(scaledImg)

	var outPath string
	if config.dither && config.invert {
		outPath = fmt.Sprintf("ascii_%s_inverted_dithered.txt", strings.Split(imagePath, ".")[0])
	} else if config.dither {
		outPath = fmt.Sprintf("ascii_%s_dithered.txt", strings.Split(imagePath, ".")[0])
	} else if config.invert {
		outPath = fmt.Sprintf("ascii_%s_inverted.txt", strings.Split(imagePath, ".")[0])
	} else {
		outPath = fmt.Sprintf("ascii_%s.txt", strings.Split(imagePath, ".")[0])
	}

	out, err := os.Create(outPath)
	if err != nil {
		fmt.Println("error creating file:", err)
		os.Exit(1)
	}
	defer out.Close()

	width, height := scaledImg.Bounds().Max.X, scaledImg.Bounds().Max.Y

	for y := 0; y < height; y++ {
		line := ""
		for x := 0; x < width; x++ {
			var char string
			if config.dither {
				char = getDitheredAscii(lightnessGrid[x][y], config.invert)
			} else {
				char = getAscii(lightnessGrid[x][y], config.invert)
			}
			line += char
		}
		out.WriteString(line + "\n")
	}

	fmt.Printf("Saved ASCII art to ./%s\n", outPath)
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
	fmt.Println("usage: img2ascii [-d] [-i] <path/to/image> [scale]")
	fmt.Println("  -d apply dithering to image")
	fmt.Println("  -i invert colors")
	fmt.Println("  scale: optional scale factor (default: 1.0)")
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
