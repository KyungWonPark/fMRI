package main

import (
	"github.com/KyungWonPark/nifti"
	"sync"
)

func convolution(img *nifti.Nifti1Image, timePoint int, seed Voxel) float32 {
	var value float32
	value = 0

	for k := -1; k < 2; k ++ {
		for j := -1; j < 2; j++ {
			for i := -1; i < 2; i++ {
				value += img.GetAt(uint32(seed.x + i), uint32(seed.y + j), uint32(seed.z + k), uint32(timePoint)) * convKernel[k + 1][j + 1][i + 1]
			}
		}
	}

	return value / 8
}

func sampling(img *nifti.Nifti1Image, order <-chan int, wg *sync.WaitGroup, outBuffer *LinBuffer) {
	for {
		timePoint, ok := <-order
		if ok {
			for i, vox := range greyVoxels {
				seed := Voxel{1 + 2 * vox.x, 1 + 2 * vox.y, 2 + 2 * vox.z}
				outBuffer[i][timePoint - 300] = convolution(img, timePoint, seed)
			}
			wg.Done()
		} else {
			break
		}
	}

	return
}

func doSampling(path string, outBuffer *LinBuffer) {
	var img nifti.Nifti1Image
	img.LoadImage(path, true)

	order := make(chan int, numWorkers)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		go sampling(&img, order, &wg, outBuffer)
	}

	wg.Add(600)
	for timePoint := 300; timePoint < 900; timePoint++ {
		order <- timePoint
	}
	wg.Wait()
	return
}
