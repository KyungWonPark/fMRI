package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sync"
)

func savePearsonHeatmap(matBuf0 *MatBuffer, img *image.RGBA, order <-chan int, wg *sync.WaitGroup) {
	for {
		row, ok := <-order
		if ok {
			for col, val := range matBuf0[row] {
				var red float64
				var green float64
				var blue float64

				if val >= 0 {
					red = 255
					green = math.Round(float64(1-val) * 255)
					blue = math.Round(float64(1-val) * 255)
				} else {
					red = math.Round(float64(1+val) * 255)
					green = math.Round(float64(1+val) * 255)
					blue = 255
				}

				color := color.RGBA{uint8(red), uint8(green), uint8(blue), 0xff}
				img.Set(row, col, color)
			}

			wg.Done()
		} else {
			break
		}
	}

	return
}

func doSavePearsonHeatmap(path string, matBuf0 *MatBuffer) {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{13362, 13362}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	order := make(chan int, numWorkers)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		go savePearsonHeatmap(matBuf0, img, order, &wg)
	}

	wg.Add(13362)
	for i := 0; i < 13362; i++ {
		order <- i
	}
	wg.Wait()

	f, err := os.Create("result_" + path + ".png")
	if err != nil {
		fmt.Println("Failed to write heatmap")
		os.Exit(3)
	}
	defer f.Close()
	png.Encode(f, img)

	return
}

func doWrite(matBuf1 *MatBuffer) {
	f, err := os.Create("./result.csv")
	defer f.Close()
	if err != nil {
		fmt.Println("Failed to write result!")
		os.Exit(3)
	}

	for i := range matBuf1 {
		txt := ""
		for j := range matBuf1[i] {
			ttxt := fmt.Sprintf("%f", matBuf1[i][j])
			txt += ttxt
			txt += ","
		}
		txt += "\n"

		f.Write([]byte(txt))
	}

	err = f.Sync()
	if err != nil {
		fmt.Println("Failed to sync result to disk!")
	}

	return
}
