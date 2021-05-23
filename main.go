package main

import (
	"drawbot/lib"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

var drawing = false
var start_drawing chan bool
var event_start_chan = make(chan chan hook.Event)

var startX = -1
var startY = -1

var scale *float64

func main() {
	filename := flag.String("img", "", "Input path of image you wish to use")
	scale = flag.Float64("scale", 1, "Image scale")

	flag.Parse()

	if *filename == "" {
		fmt.Println("You haven't specified an image with CLI flag --img")
		return
	}

	img := loadimage(*filename)

	fmt.Println("[Ctrl] + [Click] to set top left drawing position!\n[Ctrl] + [Shift] + [Q] to exit")

	register_ctrl_click(img, scale)
	register_alt_click()
	register_exit_keydown()

	s := robotgo.EventStart()
	for {
		<-robotgo.EventProcess(s)
		s = <-event_start_chan
	}
}

func register_alt_click() {
	robotgo.EventHook(hook.MouseDown, []string{"alt"}, func(e hook.Event) {
		if drawing {
			start_drawing <- true
		}
	})
}

func register_ctrl_click(img image.Image, scale *float64) {
	robotgo.EventHook(hook.MouseDown, []string{"ctrl"}, func(e hook.Event) {
		if e.Button == hook.MouseMap["left"] {
			startX = int(e.X)
			startY = int(e.Y)

			go draw_manager(
				imaging.Resize(
					img,
					int(
						float64(
							img.Bounds().Max.X,
						)*(*scale),
					),
					int(
						float64(
							img.Bounds().Max.Y,
						)*(*scale),
					),
					imaging.Lanczos,
				),
				startX,
				startY,
			)
		}
	})
}

func register_exit_keydown() {
	robotgo.EventHook(hook.KeyDown, []string{"ctrl", "shift", "q"}, func(e hook.Event) {
		drawing = false
		robotgo.EventEnd()
		os.Exit(0)
	})
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

func draw_manager(img image.Image, startX, startY int) {
	for {
		draw(img, startX, startY)
		if strings.ToLower(lib.Input("Do you wish to draw again at the same position? (yes / no)")) != "yes" {
			fmt.Println("[Ctrl] + [Click] to set new drawing position")
			register_ctrl_click(img, scale) //? Re-add ctrl click handler
			return
		}
	}
}

func draw(img image.Image, startX, startY int) {
	fmt.Println("=================================")
	draw_segments := lib.Analyze(img, lib.SetupConditionOptions())

	start_drawing = make(chan bool, 1)
	drawing = true

	fmt.Println("[Alt] + [Click] to start drawing!")
	<-start_drawing

	fmt.Println("Drawing...")

	robotgo.EventEnd() //? Remove all clickdown handlers

	event_start_chan <- robotgo.EventStart() //? Add back exit handler in case of cancel
	register_exit_keydown()

	for y, row := range draw_segments {
		columnLevel := startY + y
		for _, segment := range row {
			robotgo.MoveMouse(segment.Start+startX, columnLevel)
			robotgo.DragSmooth(segment.End+startX, columnLevel, 0.1, 0.1)
		}
	}

	register_alt_click()

	drawing = false

	fmt.Println("=================================")
}
