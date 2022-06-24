package pixelReader

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const PixelSize = 3

const NotEnoughEquals = -1

var ReadingFactor = 4

var ReadingChunck int

var MatrixSize int

type Pixel struct {
	RGB [PixelSize]byte
}

type Matrix struct {
	Pixels   [][3]byte
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

func updateResults(numEqual int, fname *string) {
	percent := float64(numEqual) / float64(MainImage.Size)
	res := Result{File: *fname, Percent: percent, NumberOfEqualPixels: numEqual}
	if percent > Best.Percent {
		Third = Second
		Second = Best
		Best = res
		return
	}
	if percent > Second.Percent {
		Third = Second
		Second = res
		return
	}
	if percent > Third.Percent {
		Third = res
	}
}

func parseFiles(file fs.DirEntry, dir *string, ch chan struct{}) {

	fn := file.Name()
	fileExtension := filepath.Ext(fn)

	if fileExtension != ".raw" {
		ch <- struct{}{}
		return
	}

	fname := filepath.Join(*dir, fn)

	if fname == MainImage.FileName {
		ch <- struct{}{}
		return
	}

	data, err := os.ReadFile(fname)

	if err != nil {
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

		//TODO
		if numEqual <= NotEnoughEquals && ReadingFactor == 1 {
			ch <- struct{}{}
			return
		}
	}

	updateResults(numEqual, &fname)
	ch <- struct{}{}
}

func checkIfStillPossible(numberOfEqualPixels, startingPos int) bool {
	availablePixels := MatrixSize - startingPos
	possibleMaxNPixels := availablePixels + numberOfEqualPixels
	possibleMaxPercent := float64(possibleMaxNPixels) / float64(MatrixSize)

	if possibleMaxPercent > Third.Percent {
		return true
	}

	return false
}

func comparePixels(r, g, b *byte, p *[3]byte) bool {
	if *r != p[0] {
		return false
	}
	if *g != p[1] {
		return false
	}
	if *b != p[2] {
		return false
	}
	return true
}

//parseAndCompareMatrixes takes the starting point in the Main image (the X coordinate) and a few bytes to read and reads them 3 by 3, comparing the resulting RGB to the Main Image
func parseAndCompareMatrixes(startingPoint int, filePiece []byte, ch chan int) {
	startingPos := startingPoint
	numberOfEqualPixels := 0

	for i := 0; i < len(filePiece); i += PixelSize {
		if equals := comparePixels(&filePiece[i], &filePiece[i+1], &filePiece[i+2], &MainImage.Pixels[startingPos]); equals {
			numberOfEqualPixels += 1
		}
		startingPos += 1

		//Early stoppage
		if ReadingFactor == 1 {
			if possible := checkIfStillPossible(numberOfEqualPixels, startingPos); !possible {
				ch <- NotEnoughEquals
				return
			}
		}
	}
	ch <- numberOfEqualPixels
}

func parseMainMatrix(startingPoint int, filePiece []byte, ch chan struct{}) {
	startingPos := startingPoint
	for i := 0; i < len(filePiece); i += PixelSize {
		pixel := [3]byte{filePiece[i], filePiece[i+1], filePiece[i+2]}
		MainImage.Pixels[startingPos] = pixel
		startingPos += 1
	}

	ch <- struct{}{}
}

func makeMainImage(data *[]byte, fn *string) {
	MatrixSize = len(*data) / PixelSize
	MainImage = Matrix{Pixels: make([][3]byte, MatrixSize), Size: MatrixSize, FileName: *fn}
	ReadingChunck = (MatrixSize * PixelSize) / ReadingFactor
}

func parseMainImage(filename *string) error {

	data, err := os.ReadFile(*filename)

	if err != nil {
		fmt.Println("Error reading main image: ", err)
		return err
	}

	makeMainImage(&data, filename)

	readCh := make(chan struct{}, ReadingFactor)

	for i := 0; i < ReadingFactor; i++ {
		readFrom := i * ReadingChunck
		startingPoint := readFrom / PixelSize
		readTo := readFrom + ReadingChunck
		go parseMainMatrix(startingPoint, data[readFrom:readTo], readCh)
	}

	for i := 0; i < ReadingFactor; i++ {
		<-readCh
	}

	return nil
}

func parseImageFiles(directory *string) error {
	imageFiles, err := os.ReadDir(*directory)

	if err != nil {
		return err
	}

	numberOfImages := len(imageFiles)
	parseCh := make(chan struct{}, numberOfImages)

	for _, el := range imageFiles {
		go parseFiles(el, directory, parseCh)
	}

	for i := 0; i < numberOfImages; i++ {
		<-parseCh
	}

	return nil
}

func PrintTopThreeString() {
	fmt.Printf("BEST: %v\nSECOND: %v\nTHIRD: %v\n", Best, Second, Third)
}

func setUp() {
	Best = Result{Percent: 0.0}
	Second = Result{Percent: 0.0}
	Third = Result{Percent: 0.0}
}

func Run(directory, mainImage *string) {

	setUp()

	imagePath := filepath.Join(*directory, *mainImage)

	err := parseMainImage(&imagePath)

	if err != nil {
		return
	}

	err = parseImageFiles(directory)

	if err != nil {
		return
	}
}
