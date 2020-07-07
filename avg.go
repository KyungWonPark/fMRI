package main

import "sync"

func average(matBuf1 *MatBuffer, n float32, order <-chan int, wg *sync.WaitGroup) {
	for {
		index, ok := <-order
		if ok {
			for i := range matBuf1[index] {
				matBuf1[index][i] = matBuf1[index][i] / n
			}

			wg.Done()
		} else {
			break
		}

	}
	return
}

func doAverage(matBuf1 *MatBuffer, n float32) {
	order := make(chan int, numWorkers)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		go average(matBuf1, n, order, &wg)
	}

	wg.Add(13362)
	for i := 0; i < 13362; i++ {
		order <- i
	}
	wg.Wait()

	return
}
