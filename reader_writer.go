package main

import (
	"fmt"
	"time"
	"math/rand"
)

type empty interface {}
type semaphore chan empty

func (s semaphore) P(n int) {
	var e empty
		for i:= 0 ; i<n; i++ {
		s <- e
	}

}

func (s semaphore) V( n int )  {
	for i := 0; i < n; i++ {
		<-s
	}
}


func pow( a int, b int ) int {
	p := 1

	for i := 0; i < b; i++ {
		p *= a
	}
	return p
}
var shared int
var readCount int
var resource semaphore
var mutex semaphore
var done semaphore
var e empty

func writer () {
	resource.P(1)
	shared++
	a := rand.Intn(2) + 2
	time.Sleep(time.Duration(a  * 1e8))
	resource.V(1)
	done <- e
}
func reader () {
	mutex.P(1)
	readCount++
	if readCount == 1 {
		resource.P(1)
	}
	mutex.V(1)
	//Critical section
	fmt.Printf("The value of shared: %d\n", shared)

	mutex.P(1)
	readCount--
	if readCount == 0 {
		resource.V(1)
	}
	mutex.V(1)
	done <- e
}
func main() {	
	rand.Seed(10)
	readCount = 0
	resource = make(semaphore,1)
	mutex = make(semaphore,1)
	done = make(semaphore, 20)
	shared = 4
	
	fmt.Printf("The value of shared is %d\n", shared)
	go reader()
	go writer() //5
	go reader()  
	go writer() //6
	go reader() 
	go reader()
	go writer() //7
	go reader()
	go reader()
	go reader() 
	go writer() //8
	go reader() 

	for i := 0; i < 12; i++ {
		<-done
	}

	fmt.Printf("The value of shared is %d\n", shared)
}