// Command funnel creates images of scaled funnels  based on input percentages
// Copyright 2016, Hariharan Srinath

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"

	"github.com/llgcode/draw2d/draw2dimg"
)

const docString = `
funnel creates a scaled funnel graph based on percentage sizes by segment.

USAGE: funnel -width [width] -height [height] -out [filename] [entries...]

PARAMETERS:
-width [width]      The width of the output image in pixels (default 400)
-height [height]    The height of the output image pixels (default 600)
-out [filename]     The filename of the output PNG image (default funnel.png)        
[entries...]        Percentages representing each segment. Max 10 entries.
                    MUST be between 0 and 100 (inclusive). 
EXAMPLE:
funnel 100 70 40 10 0
`

var colorpal = []color.RGBA{
	color.RGBA{13, 71, 161, 255},
	color.RGBA{21, 101, 192, 255},
	color.RGBA{25, 118, 210, 255},
	color.RGBA{30, 136, 229, 255},
	color.RGBA{33, 150, 243, 255},
	color.RGBA{66, 165, 245, 255},
	color.RGBA{100, 181, 246, 255},
	color.RGBA{144, 202, 249, 255},
	color.RGBA{187, 222, 251, 255},
	color.RGBA{227, 242, 253, 255},
}

func getColorPal(n int) []color.RGBA {
	switch n {
	case 1:
		return []color.RGBA{colorpal[9]}
	case 2:
		return []color.RGBA{colorpal[0], colorpal[9]}
	case 3:
		return []color.RGBA{colorpal[0], colorpal[4], colorpal[9]}
	case 4:
		return []color.RGBA{colorpal[0], colorpal[3], colorpal[6], colorpal[9]}
	case 5:
		return []color.RGBA{colorpal[0], colorpal[3], colorpal[5], colorpal[7], colorpal[9]}
	case 6:
		return []color.RGBA{colorpal[0], colorpal[2], colorpal[4], colorpal[6], colorpal[8], colorpal[9]}
	case 7:
		return []color.RGBA{colorpal[0], colorpal[1], colorpal[3], colorpal[4], colorpal[6], colorpal[8], colorpal[9]}
	case 8:
		return []color.RGBA{colorpal[0], colorpal[1], colorpal[3], colorpal[4], colorpal[5], colorpal[6], colorpal[8], colorpal[9]}
	case 9:
		return []color.RGBA{colorpal[0], colorpal[1], colorpal[2], colorpal[3], colorpal[4], colorpal[5], colorpal[6], colorpal[8], colorpal[9]}
	}
	return colorpal
}

func parseFunnel(entries []string) ([]float64, error) {
	ret := []float64{}

	for _, entry := range entries {
		d, err := strconv.ParseFloat(entry, 64)
		if err != nil {
			return nil, fmt.Errorf("Entry %s could not be interpreted as a number: %s", entry, err)
		}

		if d > 100 || d < 0 {
			return nil, fmt.Errorf("All entries must be between 100 and 0. Got:%d", d)
		}

		ret = append(ret, d)
	}

	if len(ret) == 0 {
		return nil, fmt.Errorf("Funnel must have at least 1 entry")
	}

	if len(ret) > len(colorpal) {
		return nil, fmt.Errorf("We support a max of %d funnel entries.", len(colorpal))
	}

	return ret, nil

}

func main() {
	var width, height int
	var outfile string

	flag.IntVar(&width, "width", 400, "width of the complete funnel")
	flag.IntVar(&height, "height", 600, "height of the complete funnel")
	flag.StringVar(&outfile, "out", "funnel.png", "output file name for the image generated")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, docString)
	}

	flag.Parse()

	funnel, err := parseFunnel(flag.Args())
	if err != nil {
		log.Fatalf("Error reading funnel entries: %s", err)
	}

	skipLast := false
	if funnel[len(funnel)-1] == 0 {
		skipLast = true
	}

	dest := image.NewRGBA(image.Rect(0, 0, width, height))
	gc := draw2dimg.NewGraphicContext(dest)
	verticalDelta := float64(height) / float64(len(funnel))
	if skipLast {
		verticalDelta = float64(height) / float64(len(funnel)-1)
	}

	colorPal := getColorPal(len(funnel))
	if skipLast {
		colorPal = getColorPal(len(funnel) - 1)
	}

	for j, funnelVal := range funnel {

		if skipLast && j == len(funnel)-1 {
			break
		}

		topY := float64(j) * verticalDelta
		botY := float64(j+1) * verticalDelta
		topX := float64(width) * funnelVal / 200
		botX := topX
		if j+1 < len(funnel) {
			botX = float64(width) * funnel[j+1] / 200
		}

		gc.SetFillColor(colorPal[j])
		gc.SetStrokeColor(colorPal[j])
		gc.MoveTo(float64(width)/2-topX, topY)
		gc.LineTo(float64(width)/2+topX, topY)
		gc.LineTo(float64(width)/2+botX, botY)
		gc.LineTo(float64(width)/2-botX, botY)
		gc.LineTo(float64(width)/2-topX, topY)

		gc.Close()
		gc.FillStroke()
	}

	//gc.SetLineWidth(5)

	// Save to file
	//	draw2dimg.SaveToPngFile("hello.png", dest)

	f, err := os.Create(outfile)
	if err != nil {

		log.Fatal(err)
	}

	if err := png.Encode(f, dest); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}
