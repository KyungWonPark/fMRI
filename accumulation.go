package main

import "sync"

func accumulation(matBuf0 *MatBuffer, matBuf1 *MatBuffer, order <-chan int, wg *sync.WaitGroup) {
	for {
		index, ok := <-order
		if ok {
			for i := range matBuf1[index] {
				matBuf1[index][i] += matBuf0[index][i]
			}

			wg.Done()
		} else {
			break
		}

	}
	return
}

func doAccumulation(matBuf0 *MatBuffer, matBuf1 *MatBuffer) {
	order := make(chan int, numWorkers)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		go accumulation(matBuf0, matBuf1, order, &wg)
	}

	wg.Add(13362)
	for i := 0; i < 13362; i++ {
		order <- i
	}
	wg.Wait()
	close(order)
	return
}
