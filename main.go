package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const MatrixSize = 1024
const MatrixDivider = 4
const PixelSize = 3
const ReadingFactor = 4

type Pixel struct {
	RGB [PixelSize]int
}

type Matrix struct {
	Pixels []Pixel
}

var MainImage Matrix

func (m Matrix) getReadingChunckSize() int {
	return (len(m.Pixels) * 3) / ReadingFactor
}

func parseMatrix(startingPoint, interval int, file []byte) {
	fmt.Println(startingPoint)
}

func parseFiles(file fs.FileInfo, n int, wg *sync.WaitGroup) {
	defer wg.Done()

	fileExtension := filepath.Ext(file.Name())

	if fileExtension != ".raw" {
		return
	}

	_, err := os.ReadFile("./Bronze/" + file.Name())

	if err != nil {
		fmt.Print(err)
		return
	}

	return
}

func checkAndPutPixel(pixelCount int, pixel Pixel) {
	if pixelCount < 3 {
		return
	}

}

func parseMainMatrix(startingPoint int, wg *sync.WaitGroup, filePiece []byte) {
	defer wg.Done()
	startingPos := startingPoint
	numberOfPixelsToRead := MainImage.getReadingChunckSize()
	rgbCount := 0
	pixel := Pixel{RGB: [3]int{0, 0, 0}}

	for n, el := range filePiece {
		if n == numberOfPixelsToRead {
			break
		}

		pixel.RGB[rgbCount] = int(el)
		rgbCount += 1
		if rgbCount != 3 { //if we read 3 bytes then we read a pixel
			continue
		}

		rgbCount = 0
		MainImage.Pixels[startingPos] = pixel
		startingPos += 1
	}
}

func makeMainImage(data []byte) {
	mainImageMatrixSize := len(data) / PixelSize
	MainImage = Matrix{Pixels: make([]Pixel, mainImageMatrixSize, mainImageMatrixSize)}
}

func parseMainImage(filename string) {

	data, err := os.ReadFile(filename)

	if err != nil {
		fmt.Println("Error reading main image: ", err)
		return
	}

	makeMainImage(data)

	var wg sync.WaitGroup

	readingChunck := MainImage.getReadingChunckSize()

	for i := 0; i < ReadingFactor; i += 1 {
		wg.Add(1)
		readFrom := i * readingChunck
		startingPoint := readFrom / PixelSize
		readTo := readFrom + readingChunck
		fmt.Printf("READ from: %v to %v\n", readFrom, readTo)
		parseMainMatrix(startingPoint, &wg, data[readFrom:readTo])
		//break
	}
	wg.Wait()

	fmt.Println(MainImage.Pixels)
}

func main() {

	directory := flag.String("dir", "./Bronze/", "Directory with all the images")
	mainImage := flag.String("img", "", "Main image for comparasion")

	flag.Parse()

	imagePath := filepath.Join(*directory, *mainImage)

	now := time.Now()
	parseMainImage(imagePath)
	fmt.Println("TIME:", time.Since(now))

	/*bronzeFiles, err := ioutil.ReadDir(*directory)

	if err != nil {
		return
	}

	var wg sync.WaitGroup

	for n, el := range bronzeFiles {

		wg.Add(1)

		go parseFiles(el, n, &wg)
	}
	wg.Wait()*/
}
