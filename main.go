package main

import (
	"flag"
	"fmt"
	clr "github.com/gookit/color"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"os/exec"
)

var (
	imgFile             *os.File
	img                 image.Image
	imgFileName         string
	pxlType             string
	ttyWidth, ttyHeight int
	imgWidth, imgHeight int
	x0, y0              int
	koffX, koffY        int
	verbose             bool
	err                 error
)

func main() {
	// getting tty size for beautifying outputed image

	// reading program flags
	flag.StringVar(&imgFileName, "image", "", "path to image to replicate")
	flag.StringVar(&pxlType, "type", "ascii", "characters used for outputting (values={ascii, win})")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.IntVar(&ttyWidth, "width", 0, "set tty width")
	flag.IntVar(&ttyHeight, "height", 0, "set tty height")
	flag.Parse()
	// set tty size if not set by user
	if ttyWidth == 0 {
		ttyWidth, ttyHeight = getTtySize()
	}
	if verbose {
		fmt.Printf("tty Width=%d\ntty Height=%d\n", ttyWidth, ttyHeight)
	}
	// read image
	imgFile, err = os.Open(imgFileName)
	if err != nil {
		fmt.Printf("no such image %s found, may be no such file exists?", imgFileName)
		os.Exit(0)
	}
	defer imgFile.Close()

	// decode png image to class image
	img, err = png.Decode(imgFile)
	if err != nil {
		fmt.Println("error reading image, may be it is not PNG?")
	}
	if verbose {
		fmt.Println("image decoded succesfully")
	}
	// set image values
	x0, y0 = img.Bounds().Min.X, img.Bounds().Min.Y
	imgHeight = img.Bounds().Max.Y
	imgWidth = img.Bounds().Max.X

	// find quality reduction
	koffX, koffY = imgWidth/ttyWidth+1, imgHeight/ttyHeight
	if verbose {
		fmt.Printf("quality reduction found. it is %d by X and %d by Y\n", koffX, koffY)
	}
	printAscii()
}

func getTtySize() (width, height int) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Sscanf(string(out), "%d %d", &height, &width)
	return width, height
}

func getRGB(pxl color.Color) (r, g, b uint8) {
	uint32r, uint32g, uint32b, _ := pxl.RGBA()
	r, g, b = uint8(uint32r), uint8(uint32g), uint8(uint32b)
	return
}

func printAscii() {
	// set output format chars
	levels := []string{" ", ".", ",", ";", "!", "v", "l", "L", "F", "E", "$"}
	if pxlType == "win" {
		levels = []string{" ", "░", "▒", "▓", "█"}
	}

	// print image
	length := len(levels)
	for y := y0; y < imgHeight-koffY; y += koffY {
		for x := x0; x < imgWidth-koffX; x += koffX {
			c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			r, g, b := getRGB(img.At(x, y))
			level := c.Y / uint8((255 / length))
			if level == uint8(length) {
				level--
			}
			// formatted := fmt.Sprintf("<fg=%d,%d,%d>%s</>", r, g, b, levels[level])
			// formatted := "<fg=red>A</>"
			clr.RGB(r, g, b).Print(levels[level])
		}
		clr.Print("\n")
	}

}
