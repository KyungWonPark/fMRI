package main

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
)

type MatDim struct {
	rows int
	cols int
}

func writeMat32(path string, arr []float32, dim MatDim) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("Failed to write matrix!")
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for row := 0; row < dim.rows; row++ {
		line := ""
		for col := 0; col < dim.cols; col++ {
			line += strconv.FormatFloat(float64(arr[row * dim.rows + col]), 'E', -1, 32)
			if col < dim.cols - 1 {
				line += ","
			}
		}
		line += "\n"
		w.WriteString(line)
	}
	w.Flush()

	return
}

func writeVec32(path string, arr [600]float32) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("Failed to write vector")
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	txt := ""
	for _, val := range arr {
		txt += strconv.FormatFloat(float64(val), 'E', -1, 32) + "\n"
	}
	w.WriteString(txt)
	w.Flush()

	return
}
