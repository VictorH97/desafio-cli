package main

import (
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"
)

type Report struct {
	TotalTime          int
	NumberOfRequests   int
	NumberOfRequestsOk int
	Requests           map[int]int
}

func main() {
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(7)

	for i := 0; i < int(math.Ceil(7.0/4.0)); i++ {
		for j := 0; j < int(math.Abs(float64(7*i-4))); j++ {
			go func(i int) {
				resp, err := http.Get("http://google.com")
				if err != nil {
					panic(err)
				}
				fmt.Printf("Chamou %d aqui\n", i)
				defer resp.Body.Close()

				waitGroup.Done()
			}(i)
		}

		time.Sleep(time.Second)
	}

	waitGroup.Wait()
}
