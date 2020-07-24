package main

import (
	"math"
	"sync"
)

func zScoring(linBuf0 *LinBuffer, linBuf1 *LinBuffer, order <-chan int, wg *sync.WaitGroup) {
	for {
		index, ok := <-order
		if ok {
			var valAcc float32
			var sqrAcc float32

			for _, value := range linBuf0[index] {
				valAcc += value
				sqrAcc += value * value
			}

			avg := valAcc / 600
			sqrMean := sqrAcc / 600
			stddev := float32(math.Sqrt(float64(sqrMean) - float64(avg*avg)))

			for i, value := range linBuf0[index] {
				linBuf1[index][i] = (value - avg) / stddev
			}

			wg.Done()
		} else {
			break
		}
	}

	return
}

func doZScoring(linBuf0 *LinBuffer, linBuf1 *LinBuffer) {
	order := make(chan int, numWorkers)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		go zScoring(linBuf0, linBuf1, order, &wg)
	}

	wg.Add(13362)
	for i := 0; i < 13362; i++ {
		order <- i
	}
	wg.Wait()
	close(order)
	return
}
