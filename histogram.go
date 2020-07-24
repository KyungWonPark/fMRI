package main

import (
	"fmt"
	"math"
	"sync"
)

func getStat(matBuf1 *MatBuffer, arrAccAvg *[13362]float64, arrAccSqrAvg *[13362]float64, order <-chan int, wg *sync.WaitGroup) {
	for {
		index, ok := <-order
		if ok {
			var accAvg float64
			var accSqrAvg float64

			for i := range matBuf1[index] {
				accAvg += float64(matBuf1[index][i])
				accSqrAvg += math.Pow(float64(matBuf1[index][i]), 2)
			}

			arrAccAvg[index] += (accAvg / (13362 * 13362))
			arrAccSqrAvg[index] += (accSqrAvg / (13362 * 13362))
		} else {
			break
		}
	}

}

func doHistAnalysis(matBuf1 *MatBuffer) {
	var arrAccAvg [13362]float64
	var arrAccSqrAvg [13362]float64

	order := make(chan int, numWorkers)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		go getStat(matBuf1, &arrAccAvg, &arrAccSqrAvg, order, &wg)
	}

	wg.Add(13362)
	for i := 0; i < 13362; i++ {
		order <- i
	}
	wg.Wait()

	var avg float64
	var sqrAvg float64

	for i := 0; i < 13362; i++ {
		avg += arrAccAvg[i]
		sqrAvg += arrAccSqrAvg[i]
	}

	stdDev := math.Sqrt(sqrAvg - (avg * avg))

	fmt.Printf("AVG: %f\n", avg)
	fmt.Printf("STDDEV: %f\n", stdDev)

	return
}
