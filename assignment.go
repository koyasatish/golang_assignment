package main

import (
	"fmt"
	"sync"
)

func genVal(comp_struct <-chan struct{}, nums int) <-chan int {
	val := make(chan int)
	go func() {
		defer close(val)
		for n := 1; n < nums; n += 9 {
			select {
			case val <- n:
			case <-comp_struct:
				return
			}
		}
	}()
	return val
}

func square_root(comp_struct <-chan struct{}, gen_val <-chan int) <-chan int {
	val := make(chan int)
	go func() {
		defer close(val)
		for n := range gen_val {
			select {
			case val <- n * n:
			case <-comp_struct:
				return
			}
		}
	}()
	return val
}

func merge_chanel(comp_struct <-chan struct{}, cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	val := make(chan int)

	output := func(c <-chan int) {
		defer wg.Done()
		for n := range c {
			select {
			case val <- n:
			case <-comp_struct:
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
		close(val)
	}()
	return val
}
func main() {
	fmt.Println("Starting-------->>>")
	comp_struct := make(chan struct{})
	defer close(comp_struct)

	gen_val := genVal(comp_struct, 1234567)
	channel_a := square_root(comp_struct, gen_val)
	channel_b := square_root(comp_struct, gen_val)
	channel_c := square_root(comp_struct, gen_val)

	for final_number := range merge_chanel(comp_struct, channel_a, channel_b, channel_c) {
		fmt.Println(final_number)
	}

}
