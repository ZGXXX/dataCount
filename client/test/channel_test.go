package test

import (
	"fmt"
	"sync"
	"testing"
)

func TestChannel(t *testing.T) {
	ch := make(chan int, 10)
	wg := sync.WaitGroup{}
	for i:=0;i<5;i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 10; i++ {
				ch <- i
			}
			wg.Done()
		}()
	}
	wg.Wait()

	for item := range ch{
		fmt.Println(item)
	}
}