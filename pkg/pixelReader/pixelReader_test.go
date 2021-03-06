package pixelReader

import (
	"os"
	"testing"
)

type TopThree struct {
	FirstRes  Result
	SecondRes Result
	ThirdRes  Result
}

/*func makeTopThree() TopThree {
	best := Result{Percent: 0.99}
	second := Result{Percent: 0.9}
	third := Result{Percent: 0.8}
	topThree := TopThree{FirstRes: best, SecondRes: second, ThirdRes: third}
	return topThree
}

func remakeTopThree(tt TopThree) {
	Best = tt.FirstRes
	Second = tt.SecondRes
	Third = tt.ThirdRes
}

/*func TestUpdateResults(t *testing.T) {

	//Arrange
	oldTopThree := makeTopThree()
	remakeTopThree(oldTopThree)

	t.Run("Update First Place Test", func(t *testing.T) {
		newBest := Result{Percent: 1}

		//Act
		updateResults(1024*1024, "")

		//Assert
		if Best != newBest || Second != oldTopThree.FirstRes || Third != oldTopThree.SecondRes {
			t.Errorf("failed to update Best Result")
		}

		remakeTopThree(oldTopThree)
	})

	t.Run("Update Third Place", func(t *testing.T) {
		newThird := 850000

		//Act
		updateResults(newThird, "")

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
}*/

func TestMakeMainImage(t *testing.T) {
	//Arrange
	fn := "../../Bronze/main.raw"
	testFile, err := os.ReadFile(fn)

	if err != nil {
		t.Fatalf("Failed to read testing file for image")
	}

	type ExpectedValues struct {
		matrixSize int
		matrix     Matrix
	}

	sizeTF := len(testFile)
	expected := ExpectedValues{matrixSize: sizeTF, matrix: Matrix{Pixels: make([]byte, MatrixSize), NumberPixels: sizeTF / PixelSize}}
	//Act
	makeMainImage(&testFile, &fn)

	//Assert
	if MatrixSize != expected.matrixSize || MainImage.NumberPixels != expected.matrix.NumberPixels {
		t.Error("Failed to make main image")
	}
}

func TestSetUp(t *testing.T) {
	//Arrange
	best := Result{Percent: 0.0}
	second := Result{Percent: 0.0}
	third := Result{Percent: 0.0}
	//Act
	setUp()
	//Assert

	if best != Best || second != Second || third != Third {
		t.Error("Failed to setup the 3 places")
	}
}

/*func TestComparePixels(t *testing.T) {
	t.Run("Successfull compare", func(t *testing.T) {
		//Arrange
		r := byte(255)
		g := byte(255)
		b := byte(255)
		pixel := [3]byte{255, 255, 255}
		//Act
		equal := comparePixels(&r, &g, &b, &pixel)
		//Assert
		if !equal {
			t.Error("Failed Compare Pixels")
		}
	})

	t.Run("Failure Compare", func(t *testing.T) {
		//Arrange
		//Arrange
		fn := "../../Bronze/main.raw"
		testFile, err := os.ReadFile(fn)

		if err != nil {
			t.Fatalf("Failed to read testing file for image")
		}
		makeMainImage(&testFile, &fn)
		r := byte(255)
		g := byte(255)
		b := byte(0)
		pixel := [3]byte{255, 255, 255}
		//Act
		equal := comparePixels(&r, &g, &b, &pixel)
		//Assert
		if equal {
			t.Error("Failed Compare Pixels")
		}
	})
}
*/

func TestParseMainImage(t *testing.T) {
	t.Run("fail: Wrong File Name", func(t *testing.T) {
		//Arrange
		fname := "../../Bronze/fail.raw"
		//Act
		err := parseMainImage(&fname)
		//Assert
		if err == nil {
			t.Errorf("Parsed non existing Image")
		}
	})

	t.Run("pass: Parsed existing file", func(t *testing.T) {
		//Arrange
		fname := "../../Bronze/main.raw"
		//Act
		err := parseMainImage(&fname)
		//Assert
		if err != nil {
			t.Errorf("Didn't parse existing Image")
		}
	})
}

func TestParseImageFiles(t *testing.T) {
	t.Run("Non Existent Directory", func(t *testing.T) {
		//Arrange
		nonExistentDir := "./NA"

		//Act
		err := parseImageFiles(&nonExistentDir)

		//Arrange
		if err == nil {
			t.Error("Parsed a non existint directory")
		}
	})

	t.Run("Existint Directory", func(t *testing.T) {
		//Arrange
		existingDir := "../../Bronze"

		//Act
		err := parseImageFiles(&existingDir)

		//Arrange
		if err != nil {
			t.Error("Failed to parse existing directory")
		}
	})

	t.Run("1 Reading Factor", func(t *testing.T) {
		//Arrange
		existingDir := "../../Bronze"
		oldRF := ReadingFactor
		ReadingFactor = 1
		defer func() {
			ReadingFactor = oldRF
		}()

		//Act
		err := parseImageFiles(&existingDir)

		//Arrange
		if err != nil {
			t.Error("Failed to parse existing directory")
		}
	})
}

func BenchmarkParseMainImage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//Arrange
		fname := "../../Bronze/main.raw"
		//Act
		err := parseMainImage(&fname)
		//Assert
		if err != nil {
			b.Errorf("Didn't parse existing Image")
		}
	}
}

func BenchmarkParseImageFiles(b *testing.B) {
	//Arrange
	fname := "../../Gold/main.raw"
	err := parseMainImage(&fname)
	if err != nil {
		b.Fatalf("%v", err)
	}

	existingDir := "../../Gold"

	for i := 0; i < b.N; i++ {
		//Act
		err := parseImageFiles(&existingDir)

		//Assert
		if err != nil {
			b.Error("Failed to parse existing directory")
		}
	}
}
