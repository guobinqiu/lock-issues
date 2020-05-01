package main

import (
	"fmt"
	"sync"
	"time"
)

var m sync.RWMutex
var x = 1

func main() {
	c1 := make(chan bool)
	c2 := make(chan bool)

	start := time.Now()
	go read(c1)
	go read(c2)
	<-c1
	<-c2
	fmt.Println("time taken:", time.Since(start)) //1s

	start = time.Now()
	go write(2, c1)
	go write(3, c2)
	<-c1
	<-c2
	fmt.Println("time taken:", time.Since(start)) //2s
	fmt.Println(x)
}

//读是共享锁，可以并发的读，great
func read(ch chan bool) {
	m.RLock()
	time.Sleep(time.Second)
	fmt.Println(x)
	m.RUnlock()
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
