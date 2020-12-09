package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Get user inputs
	channel := readLine("Enter color channel\n")
	chanDiff, err := strconv.ParseFloat(readLine("Enter color channel difference\n"), 64)
	gsFile := readLine("Enter greenscreen image file name\n")
	fiFile := readLine("Enter fill image file name\n")
	outFile := readLine("Enter output file name\n")

	// Validate user input
	if !strings.Contains("rgb", channel) {
		fmt.Println("Invalid color channel. Valid color channels: r, g, b")
		return
	}
	if !(chanDiff >= 1 && chanDiff <= 10) || false {
		fmt.Println("Invalid channel difference. Valid values are between 1.0 and 10.0")
		return
	}
	// Check that both images have the same dimensions
	gsImage, gsDimensions := readPpm(gsFile)
	fiImage, fiDimensions := readPpm(fiFile)
	if gsDimensions != fiDimensions {
		fmt.Println("Invalid images. Both images must be of the same size.")
		return
	}

	// Replace greenscreen with fill image
	screenChan, otherChans := getChannels(channel)
	for i, row := range gsImage {
		for j, pixel := range row {
			if float64(pixel[screenChan])/float64(pixel[otherChans[0]]) > chanDiff &&
				float64(pixel[screenChan])/float64(pixel[otherChans[1]]) > chanDiff {
				gsImage[i][j] = fiImage[i][j]
			}
		}
	}

	// Write result to output file
	file, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("P3\n" + gsDimensions + "\n255\n")
	for _, row := range gsImage {
		for _, pixel := range row {
			for _, subpixel := range pixel {
				writer.WriteString(strconv.Itoa(int(subpixel)))
				writer.Write([]byte(" "))
			}
		}
		writer.Write([]byte("\n"))
	}
}

func getChannels(char string) (int, [2]int) {
	if char == "r" {
		return 0, [2]int{1, 2}
	} else if char == "g" {
		return 1, [2]int{0, 2}
	} else if char == "b" {
		return 2, [2]int{0, 1}
	}
	log.Fatal("Invalid channel input.")
	// Compiler requires me to have a return here
	return 0, [2]int{0, 0}
}

func readPpm(filename string) ([][][3]int, string) {
	// Open file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)

	// Get image data
	dimensions := getImageDimensions(reader)
	image := loadImagePixels(reader)
	file.Close()
	return image, dimensions
}

func getImageDimensions(reader *bufio.Reader) string {
	var dimensions string
	for i := 0; i < 3; i++ {
		line, err := reader.ReadString('\n')
		if err != io.EOF && err != nil {
			break
		}
		// Set dimensions to value of second line
		if i == 1 {
			dimensions = line
		}
	}
	return dimensions
}

func loadImagePixels(reader *bufio.Reader) (image [][][3]int) {
	for {
		line, err := reader.ReadString('\n')
		if err != io.EOF && err != nil {
			break
		}
		subRow := getSubpixels(line)
		rgbRow := make([][3]int, 0, len(subRow)/3)
		for i := 0; i < len(subRow); i += 3 {
			rgbRow = append(rgbRow, [3]int{subRow[i], subRow[i+1], subRow[i+2]})
		}
		if len(rgbRow) != 0 {
			image = append(image, rgbRow)
		}
		if err != nil {
			break
		}
	}
	return
}

func getSubpixels(text string) []int {
	var subpixel int
	row := make([]int, 0, len(text)/3)
	for char := range text {
		if !(char >= 48 && char <= 58) && char != 32 {
			continue
		}
		if char != 32 {
			subpixel = subpixel*10 + (int(char) - 48)
			} else {
				row = append(row, subpixel)
				subpixel = 0
			}
		}
	return row
}
		
// Name inspired by Nim's readLine() and Python's input()
func readLine(text string) string {
	fmt.Print(text)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}
