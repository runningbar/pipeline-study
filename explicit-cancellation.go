package main 

import (
	"fmt"
	"sync"
	//"time"
)

// 第一阶段：发送数组中所有int，发完之后关闭channel
func gen2(done <-chan struct{}, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select { 
			case out <- n:
			case <-done:
				return
			}
		}
	}()
	return out
}

// 第2阶段：从inbound channel中接受int，同时算出平方值传给下游
// range可以感知到channel关闭
// 所以inbound channel关闭，传给下游的outbound channel也会关闭
func sq2(done <-chan struct{}, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-done:
				return
			}
		}
	}()
	return out
}



// 组装流水线
func main() {
	done := make(chan struct{})
	defer close(done)

	in := gen2(done, 2, 3)

	// fan-out，不过是不同协程用同一个函数读inbound channel
	c1 := sq2(done, in)
	c2 := sq2(done, in)

	out := merge2(done, c1, c2)
	fmt.Println(<-out)

	//close(done)
	//time.Sleep(5 * time.Second)
}

// fan-in，接受多个inbound channel，归一化到一个outbound channel中
// 当所有的inbound channel关闭后，再关闭outbound channel
func merge2(done <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case out <- n:
			case <-done:
				return
			}
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		fmt.Println("ok")
		close(out)
	}()
	return out
}