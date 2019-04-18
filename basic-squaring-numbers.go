package main 

import (
	"fmt"
)

// 第一阶段：发送数组中所有int，发完之后关闭channel
func gen(nums ...int) <-chan int {
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
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

// 组装流水线版本1：组装流水线并运行最终阶段
func main() {
	// 组装流水线
	c := gen(2, 3)
	out := sq(c)

	// 最终阶段，消费数据
	fmt.Println(<-out)
	fmt.Println(<-out)
}

/* // 组装版本2：因为sq的inbound和outbound channel类型一样，所以可以任意组合多次
func main() {
	for n := range sq(sq(gen(2, 3))) {
		fmt.Println(n)
	}
} */