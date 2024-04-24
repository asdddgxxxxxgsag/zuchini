package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"path"
	"strings"
)

const STRIP_HEIGHT = 3
const NUM_STRIPS = 6

func main() {
	args := os.Args[1:]
	filePath := args[0]

	if _, err := os.Stat(filePath); err == nil {

	} else if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("path does not exist")
		os.Exit(1)
	} else {
		fmt.Printf("unknown error: %v", err)
		os.Exit(1)
	}

	var inputFile *os.File
	var err error
	if inputFile, err = os.Open(filePath); err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}
	defer inputFile.Close()

	var config image.Config
	var format string
	if config, format, err = image.DecodeConfig(inputFile); err != nil {
		fmt.Printf("error decoding image: %v", err)
		os.Exit(1)
	}

	fmt.Println("Width:", config.Width, "Height:", config.Height, "Format:", format)

	_, err = inputFile.Seek(0, 0)
	if err != nil {
		fmt.Printf("error seeking image: %v", err)
		os.Exit(1)
	}

	original, _, err := image.Decode(inputFile)
	if err != nil {
		fmt.Printf("error opening image: %v", err)
		os.Exit(1)
	}

	bounds := original.Bounds()
	newHeight := bounds.Dy() + 2*STRIP_HEIGHT*NUM_STRIPS
	newRect := image.Rect(0, 0, bounds.Dx(), newHeight)
	newImage := image.NewRGBA(newRect)

	colors := make(map[int]color.RGBA)
	colors[0] = color.RGBA{255, 0, 0, 255}
	colors[1] = color.RGBA{0, 255, 0, 255}
	colors[2] = color.RGBA{0, 0, 255, 255}
	colors[3] = color.RGBA{255, 0, 255, 255}
	colors[4] = color.RGBA{255, 255, 0, 255}
	colors[5] = color.RGBA{0, 255, 255, 255}

	randomInt := rand.Intn(5)

	for i := 0; i < NUM_STRIPS; i++ {
		rectTop := image.Rect(0, i*STRIP_HEIGHT, config.Width, (i+1)*STRIP_HEIGHT)
		rectBot := image.Rect(0, config.Height+NUM_STRIPS*STRIP_HEIGHT+i*STRIP_HEIGHT, config.Width, config.Height+NUM_STRIPS*STRIP_HEIGHT+(i+1)*STRIP_HEIGHT)

		color := colors[(i+randomInt)%6]
		draw.Draw(newImage, rectTop, &image.Uniform{color}, image.Point{}, draw.Src)
		draw.Draw(newImage, rectBot, &image.Uniform{color}, image.Point{}, draw.Src)
	}

	offset := image.Point{0, STRIP_HEIGHT * NUM_STRIPS}
	draw.Draw(newImage, bounds.Add(offset), original, bounds.Min, draw.Over)

	fileDir, fileName := path.Split(filePath)
	ext := path.Ext(fileName)
	fileName = strings.TrimSuffix(fileName, ext)
	fileName = fmt.Sprintf("%s_%s%s", fileName, fmt.Sprint(newHeight), ext)

	finalPath := path.Join(fileDir, fileName)

	var newFile *os.File
	if newFile, err = os.Create(finalPath); err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}
	defer newFile.Close()

	switch format {
	case "jpeg":
		jpeg.Encode(newFile, newImage, &jpeg.Options{Quality: 90})
	case "png":
		png.Encode(newFile, newImage)
	case "gif":
		gif.Encode(newFile, newImage, &gif.Options{})
	}
}
