package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	runChannel := make(chan string, 10)
	var doneSlice []string

	wg.Add(1)
	go func() {
		defer wg.Done()
		var rSlice []string
		for {
			for _, done := range doneSlice {
				fmt.Println(done)
			}
			r := <-runChannel
			rSlice = append(rSlice, r)
			if len(rSlice) == 10 {
				fmt.Printf("\r%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
					rSlice[0], rSlice[1], rSlice[2], rSlice[3], rSlice[4], rSlice[5],
					rSlice[6], rSlice[7], rSlice[8], rSlice[9])
				rSlice = []string{}
				fmt.Printf("==================")
				time.Sleep(time.Second * 5)
			}
		}
	}()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(x int) {
			defer wg.Done()
			for j := 0; j <= 10; j++ {
				if j == 10 {
					s := fmt.Sprintf("%d-Done", i)
					runChannel <- s
				}
				s := fmt.Sprintf("%d-%d", i, j)
				runChannel <- s

			}
		}(i)

	}
	wg.Wait()
}
