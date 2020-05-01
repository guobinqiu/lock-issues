package main

import (
	"fmt"
	"sync"
	"time"
)

var m sync.Mutex
var x = 1

func main() {
	c1 := make(chan bool)
	c2 := make(chan bool)

	start := time.Now()
	go read(c1)
	go read(c2)
	<-c1
	<-c2
	fmt.Println("time taken:", time.Since(start)) //2s

	start = time.Now()
	go write(2, c1)
	go write(3, c2)
	<-c1
	<-c2
	fmt.Println("time taken:", time.Since(start)) //2s
	fmt.Println(x)
}

//读是互斥锁，读也要排队，sucks
func read(ch chan bool) {
	m.Lock()
	time.Sleep(time.Second)
	fmt.Println(x)
	m.Unlock()
	ch <- true
}

//写是互斥锁
func write(i int, ch chan bool) {
	m.Lock()
	time.Sleep(time.Second)
	x = i
	m.Unlock()
	ch <- true
}
