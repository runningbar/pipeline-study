package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		fmt.Println("father alive")

		go func() {
			//defer fmt.Println("child exit")
			time.Sleep(time.Second * 2)
			fmt.Println("child alive")
			//time.Sleep(time.Second * 3)
		}()
		defer fmt.Println("father exit")
		return 
	}()
	time.Sleep(time.Second * 3)
}