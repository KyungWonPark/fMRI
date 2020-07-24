package main

import "sync"

func average(matBuf1 *MatBuffer, n float32, mins *[13362]float32, maxs *[13362]float32, order <-chan int, wg *sync.WaitGroup) {
	for {
		index, ok := <-order
		if ok {
			myMin := matBuf1[index][0] / n
			myMax := matBuf1[index][0] / n

			for i := range matBuf1[index] {
				avgedValue := matBuf1[index][i] / n
				matBuf1[index][i] = avgedValue

				if myMin > avgedValue {
					myMin = avgedValue
				} else if myMax < avgedValue {
					myMax = avgedValue
				}
			}

			mins[index] = myMin
			maxs[index] = myMax

			wg.Done()
		} else {
			break
		}

	}
	return
}

// return min, and max value
func doAverage(matBuf1 *MatBuffer, n float32) (float32, float32) {
	order := make(chan int, numWorkers)

	var wg sync.WaitGroup

	var mins [13362]float32
	var maxs [13362]float32

	for i := 0; i < numWorkers; i++ {
		go average(matBuf1, n, &mins, &maxs, order, &wg)
	}

	wg.Add(13362)
	for i := 0; i < 13362; i++ {
		order <- i
	}
	wg.Wait()
	close(order)
	var min float32
	var max float32

	for _, val := range mins {
		if min > val {
			min = val
		}
	}

	for _, val := range maxs {
		if max < val {
			max = val
		}
	}

	return min, max
}
