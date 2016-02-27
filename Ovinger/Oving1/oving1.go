
package main

import (
	"fmt"
	"runtime"
	"time"
)

var i int = 0

func increment(){
	for j := 0; j < 1000000; j++ {
    	i++
    }	
}

func decrement(){
	for j := 0; j < 1000000; j++ {
    	i--
    }
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	go increment()
	go decrement()

	time.Sleep(100*time.Millisecond)
	fmt.Printf("i: %d\n", i)

}
