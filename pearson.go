package main

import (
	"sync"
)

type pWork struct {
	from int
	to   int
}

func pearson(linBuf0 *LinBuffer, linStat0 *LinStat, matBuf0 *MatBuffer, order <-chan int, wg *sync.WaitGroup) {
	for {
		work, ok := <-order
		if ok {
			for i := work; i < 13362; i++ {
				var accProd float32
				for t := 0; t < 600; t++ {
					accProd += linBuf0[work][t] * linBuf0[i][t]
				}

				cov := (accProd / 600) - (linStat0[work].avg * linStat0[i].avg)

				pearson := cov / (linStat0[work].stddev * linStat0[i].stddev)

				matBuf0[work][i] = pearson
				matBuf0[i][work] = pearson
			}

			wg.Done()
		} else {
			break
		}
	}

	return
}

func doPearson(linBuf0 *LinBuffer, linStat0 *LinStat, matBuf0 *MatBuffer) {
	order := make(chan int, numWorkers)

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		go pearson(linBuf0, linStat0, matBuf0, order, &wg)
	}

	wg.Add(13362)
	for i := 0; i < 13362; i++ {
		order <- i
	}
	wg.Wait()

	return
}
