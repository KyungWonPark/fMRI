package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	"github.com/ghetzel/shmtool/shm"
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
	avg    float32
	stddev float32
}

type LinStat [13362]LinStatEle

func init() {
	numWorkers = runtime.NumCPU()

	for z := 0; z < 2; z++ {
		for y := 0; y < 2; y++ {
			for x := 0; x < 2; x++ {
				taxiDist := math.Abs(float64(x-1)) + math.Abs(float64(y-1)) + math.Abs(float64(z-1))
				convKernel[z][y][x] = float32(math.Pow(2, -1*taxiDist))
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
	var matBuf1 *MatBuffer // after averaging
	// var matBuf2 *MatBuffer // after thresholding

	// shared memory preparation
	shMem, err := shm.Create(13362 * 13362 * 32)
	defer shMem.Destroy()
	if err != nil {
		log.Fatal("Failed to create shared memory region")
	}

	pBase, err := shMem.Attach()
	defer shMem.Detach(pBase)
	if err != nil {
		log.Fatal("Failed to attach shared memory region")
	}

	matBuf1 = (*MatBuffer)(pBase)
	// matBuf2 = (*MatBuffer)(unsafe.Pointer(uintptr(pBase) + uintptr(13362*13362*int(unsafe.Sizeof(float32(0))))))

	min, max := pre(matBuf1)
	// Sampling, accumulation and averaging are finished

	// Call histogram CUDA
	lower_level := fmt.Sprintf("%f", min)
	fmt.Printf("LEVEL-Lower: %s\n", lower_level)
	upper_level := fmt.Sprintf("%f", max)
	fmt.Printf("LEVEL-Upper: %s\n", upper_level)
	num_bins := 20
	fmt.Printf("NUM-Bins: %s\n", strconv.Itoa(num_bins))
	num_samples := strconv.Itoa(13362 * 13362)
	fmt.Printf("NUM-Samples: %s\n", string(num_samples))

	shMemHist, err := shm.Create(num_bins * 32)
	defer shMemHist.Destroy()
	if err != nil {
		log.Fatal("Failed to create shared memory region")
	}

	pHist, err := shMemHist.Attach()
	defer shMemHist.Detach(pHist)
	if err != nil {
		log.Fatal("Failed to attach shared memory region")
	}

	id0 := strconv.Itoa(shMem.Id)
	fmt.Printf("SHM-matBuf1 ID: %s\n", id0)
	id1 := strconv.Itoa(shMemHist.Id)
	fmt.Printf("SHM-histogram ID: %s\n", id1)

	fmt.Println("Calling CUDA...")
	cuda := exec.Command("./hist/a.out", lower_level, upper_level, strconv.Itoa(num_bins), num_samples, id0, id1)
	err = cuda.Start()
	if err != nil {
		fmt.Printf("CUDA Launch FAILED: %v\n", err)
	}
	cuda.Wait()
	if err != nil {
		fmt.Printf("CUDA Finish FAILED: %v\n", err)
	}

	fmt.Println("Printing histogram...")
	for i := 0; i < num_bins; i++ {
		index := uintptr(i)
		stride := uintptr(unsafe.Sizeof(int32(0)))

		num := *(*int32)(unsafe.Pointer(uintptr(pHist) + index*stride))

		fmt.Printf("%d ", num)
	}

	return
}

func pre(matBuf1 *MatBuffer) (float32, float32) {
	var linBuf0 LinBuffer
	var linBuf1 LinBuffer

	var matBuf0 MatBuffer

	var linStat0 LinStat

	for i := range matBuf1 {
		for j := range matBuf1[i] {
			matBuf1[i][j] = 0
		}
	}

	for _, path := range fileList {
		fmt.Printf("Processing: %s\n", path)

		doSampling(path, &linBuf0)
		fmt.Printf("Sampling finished: %s\n", path)

		doZScoring(&linBuf0, &linBuf1)
		fmt.Printf("Zscoring finished: %s\n", path)

		doSigmoid(&linBuf1, &linBuf0, &linStat0)
		fmt.Printf("Sigmoid finished: %s\n", path)

		doPearson(&linBuf0, &linStat0, &matBuf0)
		fmt.Printf("Pearson finished: %s\n", path)

		doAccumulation(&matBuf0, matBuf1)
		fmt.Printf("Accumulation finished: %s\n", path)
	}

	min, max := doAverage(matBuf1, float32(len(fileList)))

	return min, max
}
