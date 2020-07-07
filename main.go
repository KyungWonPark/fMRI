package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"math"
	"runtime"
)

var numWorkers int

// [z][y][x]
var convKernel [3][3][3]float32

type Voxel struct {
	x int
	y int
	z int
}

var fileList []string
var greyVoxels [13362]Voxel

type LinBuffer [13362][600]float32
type MatBuffer [13362][13362]float32
type LinStatEle struct {
	avg float32
	stddev float32
}

type LinStat [13362]LinStatEle

func init() {
	numWorkers = runtime.NumCPU()

	for z := 0; z < 2; z++ {
		for y := 0; y < 2; y ++ {
			for x := 0; x < 2; x++ {
				taxiDist := math.Abs(float64(x - 1)) + math.Abs(float64(y - 1))	+ math.Abs(float64(z - 1))
				convKernel[z][y][x] = float32(math.Pow(2, -1 * taxiDist))
			}
		}
	}

	f, err := os.Open("./greyList.txt")
	if err != nil {
		fmt.Printf("Failed to open file\n")
		fmt.Println(err)
		os.Exit(3)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		xyz := strings.Split(line, ",")
		x, err0 := strconv.Atoi(xyz[0])
		y, err1 := strconv.Atoi(xyz[1])
		z, err2 := strconv.Atoi(xyz[2])
		if err0 != nil || err1 != nil || err2 != nil {
			fmt.Printf("Failed to do Atoi")
			os.Exit(3)
		}
		greyVoxels[i] = Voxel{x, y, z}
		i++
	}

	g, err := os.Open("./fileList.txt")
	if err != nil {
		fmt.Printf("Failed to open file")
		os.Exit(3)
	}
	defer g.Close()

	scanner = bufio.NewScanner(g)
	for scanner.Scan() {
		line := scanner.Text()
		fileList = append(fileList, line)
	}

	return
}

func main() {
	var linBuf0 LinBuffer
	var linBuf1 LinBuffer
	var matBuf0 MatBuffer

	var linStat0 LinStat

	for _, path := range fileList {
		fmt.Printf("Processing: %s\n", path)

		doSampling(path, &linBuf0)
		fmt.Println("Sampling is done")
		print(linBuf0[9233])

		doZScoring(&linBuf0, &linBuf1)
		fmt.Println("Z Scoring is done")
		print(linBuf1[9233])

		doSigmoid(&linBuf1, &linBuf0, &linStat0)
		fmt.Println("Sigmoid is done")
		print(linBuf0[9233])

		doPearson(&linBuf0, &linStat0, &matBuf0)
		fmt.Println("Pearson is done")

		doSavePearsonHeatmap(path, &matBuf0)
		fmt.Printf("Finished processing: %s\n", path)
	}

	return
}

func print(arr [600]float32) {
	for _, val := range arr {
		fmt.Println(val)
	}

	return
}

