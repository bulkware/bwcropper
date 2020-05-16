package main

/*
	bwCropper
	---------

	A tool to rotate and crop scanned images automatically.
*/

// Define imports
import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	"golang.org/x/image/tiff"
)

// Declare variables
var (
	angle         (float64) = 1.00 // Angle to rotate to
	bestAngle     (float64)        // Best angle stored in a separate variable
	exceedCounter (int)            // Measure when threshold is exceeded
	fileExists    (bool)           // Boolean to check is a file exists or not
	filePath      (string)         // Input file path
	luminance     (float64)        // Float to measure luminance
	lumPercent    (int)            // Luminance percentage
	lumThreshold  (int)     = 75   // Threshold for luminance
	output        (string)         // Function output for printing
	pos1, pos2    (int)            // Positions for x and y
	perThreshold  (int)     = 33   // Threshold for luminance threshold
	total         (int)            // Integer to hold total percent
)

// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// Functions
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func calculatePosition(image image.Image, opts map[string]int) (int, int) {

	pos1 = opts["start1"]
	currentPixel := image.At(0, 0)

	for {

		exceedCounter = 0
		luminance = 0
		lumPercent = 0

		// Loop through pixels and calculate how many times pixel luminance exceeds
		// given threshold value.
		for pos2 = opts["start2"]; pos2 <= opts["stop2"]; pos2 += opts["step2"] {

			if opts["xy"] == 1 {
				currentPixel = image.At(pos1, pos2)
			} else {
				currentPixel = image.At(pos2, pos1)
			}

			r, g, b, _ := currentPixel.RGBA()

			luminance = 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
			luminance = math.Round(luminance / 256)

			if int(luminance) >= opts["luminance_threshold"] {
				exceedCounter++
			}
		}

		// Calculate column/row percentage
		lumPercent = exceedCounter * 100 / opts["stop2"]

		// If last column/row luminance percentage exceeds threshold, break
		if lumPercent >= opts["percentage_threshold"] {
			break
		}

		// Break according xy and step1 exceeds given value
		if opts["stop1"] == 0 {
			if pos1 <= opts["stop1"] {
				break
			}
		} else {
			if pos1 >= opts["stop1"] {
				break
			}
		}

		// Increase/decrease iterator
		pos1 = pos1 + opts["step1"]
	}

	return pos1, lumPercent
}

// Check if file exists
func checkFileExists(filePath string) (bool, string) {
	fileInfo, err := os.Stat(filePath)
	switch {
	case err != nil:
		return false, "Unable to open file: " + filePath
	case os.IsNotExist(err):
		return false, "Path does not exist: " + filePath
	case fileInfo.IsDir():
		return false, "Given path is a directory: " + filePath
	default:
		return true, "File exists."
	}
}

// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// Main
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

func main() {

	// Take input file from command-line arguments
	filePath = os.Args[1]

	// Check input file
	fileExists, output = checkFileExists(filePath)
	if fileExists == false {
		fmt.Println(output)
		os.Exit(1)
	}

	// Read image from file
	existingfilePath, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Unable to open image file: " + filePath)
		os.Exit(1)
	}
	defer existingfilePath.Close()

	// Decode image
	loadedImage, err := tiff.Decode(existingfilePath)
	if err != nil {
		fmt.Println("Unable to decode image file: " + filePath)
		os.Exit(1)
	}

	// Retrieve filename without extension
	baseName := filepath.Base(filePath)
	extlessName := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	outputFile := extlessName + "-cropped.tiff"

	// Check output file
	fileExists, output = checkFileExists(outputFile)
	if fileExists {
		fmt.Println(output)
		os.Exit(1)
	}

	// Declare variables
	blackColor := color.RGBA{uint8(0), uint8(0), uint8(0), uint8(0)}
	loadedImage2 := loadedImage
	opts := make(map[string]int)
	opts["luminance_threshold"] = lumThreshold
	opts["percentage_threshold"] = perThreshold
	results := make(map[string]int)
	width, height := loadedImage.Bounds().Max.X, loadedImage2.Bounds().Max.Y

	// Print results table header
	fmt.Println("   Angle |    Total |      Top |    Right |   Bottom |    Left ")
	fmt.Println("---------------------------------------------------------------")

	for {

		// If angle is over maximum values, correct it
		if angle < 0.00 {
			angle = 359.99
		} else if angle > 359.99 {
			angle = 0.00
		}

		// Rotate image
		loadedImage2 = imaging.Rotate(loadedImage, angle, blackColor)

		// Retrieve image dimensions
		width, height = loadedImage2.Bounds().Max.X, loadedImage2.Bounds().Max.Y

		// Options for detecting the top side
		opts["start1"] = 0
		opts["stop1"] = height
		opts["step1"] = 1
		opts["start2"] = 0
		opts["stop2"] = width
		opts["step2"] = 1
		opts["xy"] = 2

		// Calculate black borders using a function
		top, topPercent := calculatePosition(loadedImage2, opts)

		// Options for detecting the left side
		opts["start1"] = 0
		opts["stop1"] = width
		opts["step1"] = 1
		opts["start2"] = 0
		opts["stop2"] = height
		opts["step2"] = 1
		opts["xy"] = 1

		// Calculate black borders using a function
		left, leftPercent := calculatePosition(loadedImage2, opts)

		// Options for detecting the bottom side
		opts["start1"] = height
		opts["stop1"] = 0
		opts["step1"] = -1
		opts["start2"] = 0
		opts["stop2"] = width
		opts["step2"] = 1
		opts["xy"] = 2

		// Calculate black borders using a function
		bottom, bottomPercent := calculatePosition(loadedImage2, opts)

		// Options for detecting the right side
		opts["start1"] = width
		opts["stop1"] = 0
		opts["step1"] = -1
		opts["start2"] = 0
		opts["stop2"] = height
		opts["step2"] = 1
		opts["xy"] = 1

		// Calculate black borders using a function
		right, rightPercent := calculatePosition(loadedImage2, opts)

		// Calculate totalPercent from result
		total = topPercent + rightPercent + bottomPercent + leftPercent

		// Output results
		fmt.Printf("%8v", strconv.FormatFloat(angle, 'f', 2, 64))
		fmt.Printf(" | %8v", total)
		fmt.Printf(" | %8v", top)
		fmt.Printf(" | %8v", right)
		fmt.Printf(" | %8v", bottom)
		fmt.Printf(" | %8v\n", left)

		// Save results
		if results["total"] == 0 || total < results["total"] {
			bestAngle = angle
			results["total"] = total
			results["top"] = top
			results["right"] = right
			results["bottom"] = bottom
			results["left"] = left
		}

		// Decrease angle
		angle -= 0.01

		// Break if angles exceed x degrees in either direction
		if angle < 359.00 {
			if angle > 1.00 {
				break
			}
		}
	}

	// Print empty line
	fmt.Println()

	// Rotate image to best angle
	loadedImage2 = imaging.Rotate(loadedImage, bestAngle, blackColor)

	// Crop image
	croppedImage := imaging.Crop(
		loadedImage2,
		image.Rect(
			results["left"],
			results["top"],
			results["right"],
			results["bottom"]))

	// Save the resulting image
	err = imaging.Save(croppedImage, outputFile)
	if err != nil {
		fmt.Println("Failed to save image: " + outputFile)
	}

	// Notify user
	fmt.Println("File saved as: " + outputFile)
}
