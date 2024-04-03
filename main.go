package main

import (
	"errors"
	"image"
	"image/png"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/nfnt/resize"
)

func main() {
	args, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	file, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// if we get a second argument, use it as the width to resize the image
	width := 0
	if len(args) == 2 {
		width, err = strconv.Atoi(args[1])
		if err != nil {
			log.Fatal(err)
		}

	}

	pixels, err := getPixels(file, width)
	if err != nil {
		log.Fatal(err)
	}

	ascii, err := generateText(pixels)
	if err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile("ascii.txt", []byte(ascii), 0644); err != nil {
		log.Fatal(err)
	}
}

func parseArgs() ([]string, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return nil, errors.New("no arguments given. usage: imagetoascii image.png [width]")
	}
	if len(args) > 1 {
		_, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, errors.New("second argument must be an integer")
		}
	}

	return args, nil
}

type pixel struct {
	R int
	G int
	B int
	A int
}

func getPixels(imgFile io.Reader, width int) ([][]pixel, error) {
	img, _, err := image.Decode(imgFile)

	if err != nil {
		return nil, err
	}

	if width > 0 {
		img = resize.Resize(uint(width), 0, img, resize.Bilinear)
	}

	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	pixels := make([][]pixel, height)
	for y := 0; y < height; y++ {
		row := make([]pixel, width)
		for x := 0; x < width; x++ {
			row[x] = rgbaToPixel(img.At(x, y).RGBA())
		}
		pixels[y] = row
	}

	return pixels, nil
}

func rgbaToPixel(r, g, b, a uint32) pixel {
	return pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func generateText(pixels [][]pixel) (string, error) {
	if len(pixels) == 0 {
		return "", errors.New("empty image")
	}

	text := ""

	// iterate through columns and rows
	for _, row := range pixels {
		for _, pixel := range row {
			text += getChar(pixel)
		}
		text += "\n"
	}

	return text, nil
}

func getChar(pixel pixel) string {
	avg := average(pixel)

	// we are going to use 16 symbols in the ascii art:
	// " ", ".", ":", ";", "=", "+", "*", "!", "?", "^", "&", "#", "$", "%", "@", "█"
	// the breakpoints are: 0, 16, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 224, 240

	if avg < 16 {
		return " "
	} else if avg < 32 {
		return "."
	} else if avg < 48 {
		return ":"
	} else if avg < 64 {
		return ";"
	} else if avg < 80 {
		return "="
	} else if avg < 96 {
		return "+"
	} else if avg < 112 {
		return "*"
	} else if avg < 128 {
		return "!"
	} else if avg < 144 {
		return "?"
	} else if avg < 160 {
		return "^"
	} else if avg < 176 {
		return "&"
	} else if avg < 192 {
		return "#"
	} else if avg < 208 {
		return "$"
	} else if avg < 224 {
		return "%"
	} else if avg < 240 {
		return "@"
	}
	return "█"
}

func average(pixel pixel) int {
	return (pixel.R + pixel.G + pixel.B) / 3
}
