package main

import (
	"flag"
	"fmt"
	"io/fs"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"time"
)

const PixelSize = 3
const ReadingFactor = 4

var ReadingChunck int

var MatrixSize int

type Pixel struct {
	RGB [PixelSize]int
}

type Matrix struct {
	Pixels   []Pixel
	Size     int
	FileName string
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

func parseFiles(file fs.DirEntry, dir string, ch chan interface{}) {

	fileExtension := filepath.Ext(file.Name())

	if fileExtension != ".raw" {
		ch <- struct{}{}
		return
	}

	fname := filepath.Join(dir, file.Name())

	if fname == MainImage.FileName {
		ch <- struct{}{}
		return
	}

	data, err := os.ReadFile(fname)

	if err != nil {
		fmt.Print(err)
		ch <- struct{}{}
		return
	}

	numEqual := 0

	resultChannel := make(chan int, ReadingFactor)

	for i := 0; i < ReadingFactor; i++ {
		readFrom := i * ReadingChunck
		startingPoint := readFrom / PixelSize
		readTo := readFrom + ReadingChunck
		go parseAndCompareMatrixes(startingPoint, data[readFrom:readTo], resultChannel)
	}

	for i := 0; i < ReadingFactor; i++ {
		numEqual += <-resultChannel
	}

	result := Result{File: fname, Percent: float64(numEqual) / float64(MainImage.Size), NumberOfEqualPixels: numEqual}
	//FinalResult = append(FinalResult, result)

	updateResults(result)
	ch <- struct{}{}
}

//parseAndCompareMatrixes takes the starting point in the Main image (the X coordinate) and a few bytes to read and reads them 3 by 3, comparing the resulting RGB to the Main Image
func parseAndCompareMatrixes(startingPoint int, filePiece []byte, ch chan int) {
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

func parseMainMatrix(startingPoint int, filePiece []byte, ch chan interface{}) {
	startingPos := startingPoint

	for i := 0; i < len(filePiece); i += PixelSize {
		pixel := Pixel{RGB: [3]int{int(filePiece[i]), int(filePiece[i+1]), int(filePiece[i+2])}}
		MainImage.Pixels[startingPos] = pixel
		startingPos += 1
	}

	ch <- struct{}{}
}

func makeMainImage(data []byte, fn string) {
	MatrixSize = len(data) / PixelSize
	MainImage = Matrix{Pixels: make([]Pixel, MatrixSize), Size: MatrixSize, FileName: fn}
	ReadingChunck = (MatrixSize * PixelSize) / ReadingFactor
}

func parseMainImage(filename string) (bool, error) {

	data, err := os.ReadFile(filename)

	if err != nil {
		fmt.Println("Error reading main image: ", err)
		return false, err
	}

	makeMainImage(data, filename)

	readCh := make(chan interface{}, ReadingFactor)

	for i := 0; i < ReadingFactor; i += 1 {
		readFrom := i * ReadingChunck
		startingPoint := readFrom / PixelSize
		readTo := readFrom + ReadingChunck
		go parseMainMatrix(startingPoint, data[readFrom:readTo], readCh)
	}

	for i := 0; i < ReadingFactor; i++ {
		<-readCh
	}

	return true, nil
}

func parseImageFiles(directory string) error {
	imageFiles, err := os.ReadDir(directory)

	if err != nil {
		return err
	}

	parseCh := make(chan interface{})

	for _, el := range imageFiles {

		go parseFiles(el, directory, parseCh)
	}

	for i := 0; i < len(imageFiles); i++ {
		<-parseCh
	}

	return nil
}

func main() {
	directory := flag.String("dir", ".\\Bronze\\", "Directory with all the images")
	mainImage := flag.String("img", "", "Main image for comparasion")

	flag.Parse()

	Best = Result{Percent: 0}
	Second = Result{Percent: 0}
	Third = Result{Percent: 0}

	imagePath := filepath.Join(*directory, *mainImage)

	now := time.Now()

	_, err := parseMainImage(imagePath)

	if err != nil {
		return
	}

	err = parseImageFiles(*directory)

	if err != nil {
		return
	}

	fmt.Printf("BEST: %v\nSECOND: %v\nTHIRD: %v\n", Best, Second, Third)

	fmt.Println("TIME:", time.Since(now))
}
