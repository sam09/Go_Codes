package main 

import (
	"fmt"
	"time"
)

type Empty interface {}
type semaphore chan Empty
var empty Empty

func producer (ch chan int, sem semaphore) {
	for i :=0 ; i < 10 ; i++{
		ch <-i*10
		time.Sleep(1e9)
		sem <- empty
	}
}

func consumer (ch chan int, sem semaphore) {
	for {
		temp := <-ch
		fmt.Printf("Consumed %d\n", temp)
	}
}

func main() {
	buf := 2
	ch := make(chan int, buf)
	sem := make(semaphore, buf)
	go producer(ch, sem)
	time.Sleep(2e9)
	go consumer(ch, sem)
	go consumer(ch, sem)

	for i:= 0; i<10; i++ {
		<-sem
	}

}