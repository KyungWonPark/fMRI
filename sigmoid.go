package main

import (
	"math"
	"sync"
)

func sigmoid(linBuf1 *LinBuffer, linBuf0 *LinBuffer, linStat0 *LinStat, order <-chan int, wg *sync.WaitGroup) {
	for {
		index, ok := <-order
		if ok {
			var valAcc float32
			var sqrAcc float32

			for i, value := range linBuf1[index] {
				newVal := float32(2/(1+math.Exp(-float64(value))) - 1)
				valAcc += newVal
				sqrAcc += newVal * newVal
				linBuf0[index][i] = newVal
			}

			avgVal := valAcc / 600
			avgSqr := sqrAcc / 600
			stdDev := float32(math.Sqrt(float64(avgSqr) - float64(avgVal*avgVal)))

			linStat0[index] = LinStatEle{
				avg:    avgVal,
				stddev: stdDev,
			}

			wg.Done()
		} else {
			break
		}
	}

	return
}

func doSigmoid(linBuf1 *LinBuffer, linBuf0 *LinBuffer, linStat0 *LinStat) {
	order := make(chan int, numWorkers)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		go sigmoid(linBuf1, linBuf0, linStat0, order, &wg)
	}

	wg.Add(13362)
	for i := 0; i < 13362; i++ {
		order <- i
	}
	wg.Wait()
	close(order)
	return
}
