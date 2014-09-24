package main

import (
	"image"
	"image/png" // register the PNG format with the image package
	"image/color"
	"os"
	"log"
)

/**
 ./convertToGrayscale ./resources/code.png  ./resources/new.png
 */
func main() {

	log.Println("args:",os.Args)

	infile, err := os.Open(os.Args[1])
	if err != nil {
		// replace this with real error handling
		log.Println(err)
	}
	defer infile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, err := png.Decode(infile)
	if err != nil {
		// replace this with real error handling
		log.Println(err)
	}

	// Create a new grayscale image
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	gray := image.NewGray(bounds)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := src.At(x, y)

			grayColor := color.GrayModel.Convert(oldColor)

			gray.Set(x, y, grayColor)
			//r,g,b,a := grayColor.RGBA()
			log.Println("x:",x, " y:", y, " color:", grayColor)
		}
	}

	// Encode the grayscale image to the output file
	outfile, err := os.Create(os.Args[2])
	if err != nil {
		// replace this with real error handling
		panic(err)
	}
	defer outfile.Close()
	png.Encode(outfile, gray)
}
