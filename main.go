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

const PixelSize = 3
const ReadingFactor = 4

var MatrixSize int

type Pixel struct {
	RGB [PixelSize]int
}

type Matrix struct {
	Pixels []Pixel
	Size   int
}

type Result struct {
	File                string
	Percent             float64
	NumberOfEqualPixels int
}

var MainImage Matrix
var Best Result
var Second Result
var Third Result

func (m Matrix) getReadingChunckSize() int {
	return (m.Size * 3) / ReadingFactor
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

	result := Result{File: fname, Percent: float64(numEqual) / float64(MainImage.Size), NumberOfEqualPixels: numEqual}
	//FinalResult = append(FinalResult, result)

	updateResults(result)
}

func parseAndCompareMatrixes(startingPoint int, wg *sync.WaitGroup, filePiece []byte, ch chan int) {
	defer wg.Done()
	startingPos := startingPoint
	numberOfEqualPixels := 0

	for i := 0; i < len(filePiece); i += PixelSize {
		pixel := Pixel{RGB: [3]int{int(filePiece[i]), int(filePiece[i+1]), int(filePiece[i+2])}}
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

	for i := 0; i < len(filePiece); i += PixelSize {
		pixel := Pixel{RGB: [3]int{int(filePiece[i]), int(filePiece[i+1]), int(filePiece[i+2])}}
		MainImage.Pixels[startingPos] = pixel
		startingPos += 1
	}
}

func makeMainImage(data []byte) {
	MatrixSize = len(data) / PixelSize
	MainImage = Matrix{Pixels: make([]Pixel, MatrixSize, MatrixSize), Size: MatrixSize}
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
		go parseMainMatrix(startingPoint, &wg, data[readFrom:readTo])
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

	//FinalResult = []Result{}
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

	fmt.Printf("BEST: %v\nSECOND: %v\nTHIRD: %v\n", Best, Second, Third)

	fmt.Println("TIME:", time.Since(now))
}
