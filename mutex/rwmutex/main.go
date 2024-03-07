package main

import (
	"fmt"
	"sync"
	"time"
)

type resource struct {
	m    sync.RWMutex
	data int //临界资源
}

func (r *resource) read() {
	r.m.RLock()
	defer r.m.RUnlock()
	time.Sleep(time.Second)
	fmt.Println("Reading data:", r.data)
}

func (r *resource) write(i int) {
	r.m.Lock()
	defer r.m.Unlock()
	r.data = i
	time.Sleep(time.Second)
	fmt.Println("Writing data:", i)
}

func main() {
	r := &resource{data: 1}

	var wg sync.WaitGroup

	//读-读不互斥
	start := time.Now()
	wg.Add(2)
	go func() {
		r.read()
		wg.Done()
	}()
	go func() {
		r.read()
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("time taken for read-read concurrency:", time.Since(start)) //1s

	//读-写互斥
	start = time.Now()
	wg.Add(2)
	go func() {
		r.write(2)
		wg.Done()
	}()
	go func() {
		r.read()
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("time taken for read-write exclusivity:", time.Since(start)) //2s

	//写-写互斥
	start = time.Now()
	wg.Add(2)
	go func() {
		r.write(2)
		wg.Done()
	}()
	go func() {
		r.write(3)
		wg.Done()
	}()
	wg.Wait()
	fmt.Println("time taken for write-write exclusivity:", time.Since(start)) //2s
}
