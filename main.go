package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"

	"github.com/disintegration/imaging"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

var drawing = false

func main() {
	filename := flag.String("img", "", "Input path of image you wish to use")
	scale := flag.Float64("scale", 1, "Image scale")

	flag.Parse()

	if *filename == "" {
		fmt.Println("You haven't specified an image with CLI flag --img")
		return
	}

	img := loadimage(*filename)

	fmt.Println("Ctrl + Click to start drawing! Ctrl + Shift + Q to cancel")
	robotgo.EventHook(hook.MouseDown, []string{"ctrl"}, func(e hook.Event) {
		if drawing {
			return
		}

		if e.Button == hook.MouseMap["left"] {
			go draw(
				imaging.Resize(
					img,
					int(
						float64(
							img.Bounds().Max.X,
						)*(*scale),
					),
					int(
						float64(
							img.Bounds().Max.X,
						)*(*scale),
					),
					imaging.Lanczos,
				),
				int(e.X),
				int(e.Y),
			)
		}
	})

	robotgo.EventHook(hook.KeyDown, []string{"ctrl", "shift", "q"}, func(e hook.Event) {
		drawing = false
		robotgo.EventEnd()
		os.Exit(0)
	})

	s := robotgo.EventStart()
	<-robotgo.EventProcess(s)
}

func loadimage(filename string) image.Image {
	reader, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	return img
}

func draw(img image.Image, startX, startY int) {
	drawing = true

	imageX := img.Bounds().Max.X
	imageY := img.Bounds().Max.Y

	currentX := 0
	currentY := 0

	for i := 0; i < imageX*imageY; i++ {
		if !drawing {
			return
		}

		color := img.At(currentX, currentY)
		r, g, b, _ := color.RGBA()

		val := max(int(r), int(g), int(b))
		if val > 200 { //? If brightness of pixel is larger than threshhold
			robotgo.Move(
				int(startX)+currentX,
				int(startY)+currentY,
			)
			robotgo.Click()
		}

		currentX += 1

		if currentX == imageX {
			currentX = 0
			currentY += 1
		}
	}

	robotgo.EventEnd()
}

func max(numbers ...int) int {
	largest := numbers[0]
	for n := range numbers {
		if n > largest {
			largest = n
		}
	}

	return largest
}
