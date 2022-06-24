package main

import (
	"GO_Pixel_Matrix_Compare/pkg/pixelReader"
	"flag"
	"fmt"
	"time"
)

func main() {
	directory := flag.String("dir", "..\\Bronze\\", "Directory with all the images")
	mainImage := flag.String("img", "", "Main image for comparasion")

	flag.Parse()

	now := time.Now()

	pixelReader.Run(directory, mainImage)

	pixelReader.PrintTopThreeString()

	fmt.Printf("TIME: %v\n", time.Since(now))
}
