package main

import (
	"os"
	"testing"
)

func getMatrixForTests() Matrix {
	matrixSize := MatrixSize * MatrixSize / PixelSize
	matrix := Matrix{Pixels: make([]Pixel, matrixSize, matrixSize), Size: matrixSize}
	return matrix
}

type TopThree struct {
	FirstRes  Result
	SecondRes Result
	ThirdRes  Result
}

func makeTopThree() TopThree {
	best := Result{Percent: 0.99}
	second := Result{Percent: 0.9}
	third := Result{Percent: 0.8}
	topThree := TopThree{FirstRes: best, SecondRes: second, ThirdRes: third}
	return topThree
}

func TestGetReadingChunckSize(t *testing.T) {

	//Arrange
	matrix := getMatrixForTests()
	expected := (matrix.Size * PixelSize) / ReadingFactor

	//Act
	received := matrix.getReadingChunckSize()

	//Assert
	if expected != received {
		t.Errorf("Expected (%d) != received (%d)", expected, received)
	}
}

func remakeTopThree(tt TopThree) {
	Best = tt.FirstRes
	Second = tt.SecondRes
	Third = tt.ThirdRes
}

func TestUpdateResults(t *testing.T) {

	//Arrange
	oldTopThree := makeTopThree()
	remakeTopThree(oldTopThree)

	t.Run("Update First Place Test", func(t *testing.T) {
		newBest := Result{Percent: 1}

		//Act
		updateResults(newBest)

		//Assert
		if Best != newBest || Second != oldTopThree.FirstRes || Third != oldTopThree.SecondRes {
			t.Errorf("failed to update Best Result")
		}

		remakeTopThree(oldTopThree)
	})

	t.Run("Update Third Place", func(t *testing.T) {
		newThird := Result{Percent: 0.81}

		//Act
		updateResults(newThird)

		//Assert
		if Best != oldTopThree.FirstRes || Second != oldTopThree.SecondRes || Third != newThird {
			t.Errorf("failed to update Third Best Result")
		}

		remakeTopThree(oldTopThree)
	})

	t.Run("Don't update Top3", func(t *testing.T) {
		worsePercent := Result{Percent: 0.79}

		//Act
		updateResults(worsePercent)

		//Assert
		if Best != oldTopThree.FirstRes || Second != oldTopThree.SecondRes || Third != oldTopThree.ThirdRes {
			t.Errorf("updated when it shouldn't")
		}

		remakeTopThree(oldTopThree)

	})
}

func TestMakeMainImage(t *testing.T) {
	//Arrange
	testFile, err := os.ReadFile("./Bronze/main.raw")

	if err != nil {
		t.Fatalf("Failed to read testing file for image")
	}

	type ExpectedValues struct {
		matrixSize int
		matrix     Matrix
	}

	expected := ExpectedValues{matrixSize: len(testFile) / PixelSize, matrix: Matrix{Pixels: make([]Pixel, len(testFile), len(testFile)), Size: len(testFile) / PixelSize}}
	//Act
	makeMainImage(testFile)

	//Assert
	if MatrixSize != expected.matrixSize || MainImage.Size != expected.matrix.Size {
		t.Error("Failed to make main image")
	}
}

func TestParseMainImage(t *testing.T) {
	t.Run("fail: Wrong File Name", func(t *testing.T) {
		//Arrange
		fname := "./Bronze/fail.raw"
		//Act
		parsed, err := parseMainImage(fname)
		//Assert
		if parsed != false || err == nil {
			t.Errorf("Parsed non existing Image")
		}
	})

	t.Run("pass: Parsed existing file", func(t *testing.T) {
		//Arrange
		fname := "./Bronze/main.raw"
		//Act
		parsed, err := parseMainImage(fname)
		//Assert
		if parsed != true || err != nil {
			t.Errorf("Didn't parse existing Image")
		}
	})
}

func TestParseImageFiles(t *testing.T) {
	t.Run("Non Existent Directory", func(t *testing.T) {
		//Arrange
		nonExistentDir := "./NA"

		//Act
		err := parseImageFiles(nonExistentDir)

		//Arrange
		if err == nil {
			t.Error("Parsed a non existint directory")
		}
	})

	t.Run("Existint Directory", func(t *testing.T) {
		//Arrange
		existingDir := "./Bronze"

		//Act
		err := parseImageFiles(existingDir)

		//Arrange
		if err != nil {
			t.Error("Failed to parse existing directory")
		}
	})
}

func BenchmarkParseMainImage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//Arrange
		fname := "./Bronze/main.raw"
		//Act
		parsed, err := parseMainImage(fname)
		//Assert
		if parsed != true || err != nil {
			b.Errorf("Didn't parse existing Image")
		}
	}
}

func BenchmarkParseImageFiles(b *testing.B) {
	//Arrange
	fname := "./Bronze/main.raw"
	_, err := parseMainImage(fname)
	if err != nil {
		b.Fatalf("%v", err)
	}

	existingDir := "./Gold"

	for i := 0; i < b.N; i++ {
		//Act
		err := parseImageFiles(existingDir)

		//Assert
		if err != nil {
			b.Error("Failed to parse existing directory")
		}
	}
}
