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

const MatrixDivider = 4
const PixelSize = 3
const ReadingFactor = 4

var MatrixSize int

type Pixel struct {
	RGB [PixelSize]int
}

type Matrix struct {
	Pixels []Pixel
}

type Result struct {
	File                string
	Percent             float64
	NumberOfEqualPixels int
}

var MainImage Matrix
var FinalResult []Result
var Best Result
var Second Result
var Third Result

func (m Matrix) getReadingChunckSize() int {
	return (m.getSize() * 3) / ReadingFactor
}

func (m Matrix) getSize() int {
	return len(m.Pixels)
}

func parseMatrix(startingPoint, interval int, file []byte) {
	fmt.Println(startingPoint)
}

func updateResults(res Result) {
	if res.Percent > Best.Percent {
		Third = Second
		Second = Best
		Best = res
		return
	}
	if res.Percent > Second.Percent {
		Third = Second
		Second = res
		return
	}
	if res.Percent > Third.Percent {
		Third = res
	}
}

func parseFiles(file fs.DirEntry, wg *sync.WaitGroup, dir string) {
	defer wg.Done()

	fileExtension := filepath.Ext(file.Name())

	if fileExtension != ".raw" {
		return
	}

	fname := filepath.Join(dir, file.Name())
	data, err := os.ReadFile(fname)

	if err != nil {
		fmt.Print(err)
		return
	}

	var wgLocal sync.WaitGroup
	readingChunck := MainImage.getReadingChunckSize()

	numEqual := 0

	for i := 0; i < ReadingFactor; i++ {
		wgLocal.Add(1)
		readFrom := i * readingChunck
		startingPoint := readFrom / PixelSize
		readTo := readFrom + readingChunck
		resultChannel := make(chan int)
		go parseAndCompareMatrixes(startingPoint, &wgLocal, data[readFrom:readTo], resultChannel)
		numEqual += <-resultChannel
	}

	wgLocal.Wait()

	result := Result{File: fname, Percent: float64(numEqual) / float64(MainImage.getSize()), NumberOfEqualPixels: numEqual}
	FinalResult = append(FinalResult, result)

	updateResults(result)
}

func checkAndPutPixel(pixelCount int, pixel Pixel) {
	if pixelCount < 3 {
		return
	}

}

func parseAndCompareMatrixes(startingPoint int, wg *sync.WaitGroup, filePiece []byte, ch chan int) {
	defer wg.Done()
	startingPos := startingPoint
	numberOfPixelsToRead := MainImage.getReadingChunckSize()
	rgbCount := 0
	pixel := Pixel{RGB: [3]int{0, 0, 0}}
	numberOfEqualPixels := 0

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
		if pixel == MainImage.Pixels[startingPos] {
			numberOfEqualPixels += 1
		}
		startingPos += 1
	}
	ch <- numberOfEqualPixels
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
	MatrixSize = len(data) / PixelSize
	MainImage = Matrix{Pixels: make([]Pixel, MatrixSize, MatrixSize)}
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
		parseMainMatrix(startingPoint, &wg, data[readFrom:readTo])
		//break
	}
	wg.Wait()
}

func main() {

	directory := flag.String("dir", "./Bronze/", "Directory with all the images")
	mainImage := flag.String("img", "", "Main image for comparasion")

	flag.Parse()

	imagePath := filepath.Join(*directory, *mainImage)

	now := time.Now()
	parseMainImage(imagePath)

	FinalResult = []Result{}
	Best = Result{Percent: 0}
	Second = Result{Percent: 0}
	Third = Result{Percent: 0}

	bronzeFiles, err := os.ReadDir(*directory)

	if err != nil {
		return
	}

	var wg sync.WaitGroup

	for _, el := range bronzeFiles {

		wg.Add(1)

		go parseFiles(el, &wg, *directory)
	}
	wg.Wait()

	fmt.Printf("%v\n%v\n%v\n", Best, Second, Third)

	fmt.Println("TIME:", time.Since(now))
}
