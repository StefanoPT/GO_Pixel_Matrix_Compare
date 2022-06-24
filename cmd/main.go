package main

import (
	"flag"
	"fmt"
	"pixelReader"
	"time"
)

func main() {
	directory := flag.String("dir", "..\\Bronze\\", "Directory with all the images")
	mainImage := flag.String("img", "", "Main image for comparasion")

	flag.Parse()

	now := time.Now()

	pixelReader.Run(directory, mainImage)

	topThree := pixelReader.GetTopThreeString()

	fmt.Println(topThree)

	fmt.Println("TIME:", time.Since(now))
}
