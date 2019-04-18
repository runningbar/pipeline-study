package main 

import (
	"fmt"
	"sync"
	"time"
)

// 第一阶段：发送数组中所有int，发完之后关闭channel
func gen1(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// 第2阶段：从inbound channel中接受int，同时算出平方值传给下游
// range可以感知到channel关闭
// 所以inbound channel关闭，传给下游的outbound channel也会关闭
func sq1(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}



// 组装版本2：因为sq的inbound和outbound channel类型一样，所以可以任意组合多次
func main() {
	in := gen1(2, 3)

	// fan-out，不过是不同协程用同一个函数读inbound channel
	c1 := sq1(in)
	c2 := sq1(in)


	/* // 提前停止消费，阻塞生产协程，导致生产协程泄露
	out := merge1(c1, c2)
	fmt.Println(<-out)
	time.Sleep(100000 * time.Second) // 模拟父协程去干其他事情了，这样就会导致merge1中有一个生产协程被阻塞而泄露
	return */

	for n := range merge1(c1, c2) {
		fmt.Println(n)
	}
}

// fan-in，接受多个inbound channel，归一化到一个outbound channel中
// 当所有的inbound channel关闭后，再关闭outbound channel
func merge1(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}