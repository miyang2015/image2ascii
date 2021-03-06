package convert

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestGetTerminalScreenSize test fetch the terminal screen size
func TestGetTerminalScreenSize(t *testing.T) {
	assertions := assert.New(t)
	_, _, err := getTerminalScreenSize()
	assertions.True(err != nil)
}

// TestOpenImageFile test open different type image file
func TestOpenImageFile(t *testing.T) {
	assertions := assert.New(t)
	jpgFilename := "testdata/jpg_sample_image.jpg"
	openedImage, err := OpenImageFile(jpgFilename)
	assertions.True(err == nil, "jpg image format should be supported")
	assertions.True(openedImage != nil, "opened jpg file should not be nil")

	pngFilename := "testdata/png_sample_image.png"
	openedImage, err = OpenImageFile(pngFilename)
	assertions.True(err == nil, "png image format should be supported")
	assertions.True(openedImage != nil, "opened jpg file should not be nil")

	notSupported := "testdata/not_supported_sample_image"
	openedImage, err = OpenImageFile(notSupported)
	assertions.True(err != nil, "should not open unsupported image")
	assertions.True(openedImage == nil, "not supported image should be nil")
}

// TestScaleImageWithFixedHeight test scale the image by fixed height
func TestScaleImageWithFixedHeight(t *testing.T) {
	assertions := assert.New(t)
	imageFilePath := "testdata/cat_2000x1500.jpg"
	img, err := OpenImageFile(imageFilePath)
	assertions.True(img != nil)
	assertions.True(err == nil)

	options := defaultOptions
	options.Colored = false
	options.ExpectedHeight = 100

	scaledImage := ScaleImage(img, &options)
	sz := scaledImage.Bounds()
	oldSz := img.Bounds()
	assertions.Equal(100, sz.Max.Y, "scaled image height should be 100")
	assertions.Equal(oldSz.Max.X, sz.Max.X, "scaled image width should be changed")
}

// TestScaleImageWithFixedWidth test scale the image by fixed width
func TestScaleImageWithFixedWidth(t *testing.T) {
	assertions := assert.New(t)
	imageFilePath := "testdata/cat_2000x1500.jpg"
	img, err := OpenImageFile(imageFilePath)
	assertions.True(img != nil)
	assertions.True(err == nil)

	options := defaultOptions
	options.Colored = false
	options.ExpectedWidth = 200

	scaledImage := ScaleImage(img, &options)
	sz := scaledImage.Bounds()
	oldSz := img.Bounds()
	assertions.Equal(oldSz.Max.Y, sz.Max.Y, "scaled image height should be changed")
	assertions.Equal(200, sz.Max.X, "scaled image width should be 200")
}

// TestScaleImageWithFixedWidthHeight test scale the image by fixed width
func TestScaleImageWithFixedWidthHeight(t *testing.T) {
	assertions := assert.New(t)
	imageFilePath := "testdata/cat_2000x1500.jpg"
	img, err := OpenImageFile(imageFilePath)
	assertions.True(img != nil)
	assertions.True(err == nil)

	options := defaultOptions
	options.Colored = false
	options.ExpectedWidth = 200
	options.ExpectedHeight = 100

	scaledImage := ScaleImage(img, &options)
	sz := scaledImage.Bounds()
	assertions.Equal(100, sz.Max.Y, "scaled image height should be 100")
	assertions.Equal(200, sz.Max.X, "scaled image width should be 200")
}

// TestScaleImageByRatio test scale image by ratio
func TestScaleImageByRatio(t *testing.T) {
	assertions := assert.New(t)
	imageFilePath := "testdata/cat_2000x1500.jpg"
	img, err := OpenImageFile(imageFilePath)
	assertions.True(img != nil)
	assertions.True(err == nil)

	options := defaultOptions
	options.Colored = false
	options.Ratio = 0.5

	scaledImage := ScaleImage(img, &options)
	sz := scaledImage.Bounds()
	oldSz := img.Bounds()
	expectedHeight := int(float64(oldSz.Max.Y) * 0.5 * charWidth())
	expectedWidth := int(float64(oldSz.Max.X) * 0.5)
	assertions.Equal(expectedHeight, sz.Max.Y, fmt.Sprintf("scaled image height should be %d", expectedHeight))
	assertions.Equal(expectedWidth, sz.Max.X, fmt.Sprintf("scaled image width should be %d", expectedHeight))
}

// TestScaleToFitTerminalSize test scale image to fit the terminal
func TestScaleToFitTerminalSize(t *testing.T) {
	assertions := assert.New(t)
	imageFilePath := "testdata/cat_2000x1500.jpg"
	img, err := OpenImageFile(imageFilePath)
	assertions.True(img != nil)
	assertions.True(err == nil)

	options := defaultOptions
	options.Colored = false
	options.FitScreen = true

	// not terminal
	if !isatty.IsTerminal(os.Stdout.Fd()) &&
		!isatty.IsCygwinTerminal(os.Stdout.Fd()) &&
		os.Getenv("BE_CRASHER") == "1" {
		ScaleImage(img, &options)
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestScaleToFitTerminalSize")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	stdout, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Check that the log fatal message is what we expected
	gotBytes, _ := ioutil.ReadAll(stdout)
	got := string(gotBytes)
	expected := "can not detect the terminal, please disable the '-s fitScreen' option"
	if !strings.HasSuffix(got[:len(got)-1], expected) {
		t.Fatalf("Unexpected log message. Got %s but should contain %s", got[:len(got)-1], expected)
	}

	// Check that the program exited
	err = cmd.Wait()
	if e, ok := err.(*exec.ExitError); !ok || e.Success() {
		t.Fatalf("Process ran with err %v, want exit status 1", err)
	}
}
