package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	num := uint32(18)
	fmt.Printf("The number: %d\n", num)
	delta := int32(-3)
	atomic.AddUint32(&num, uint32(delta))
	fmt.Printf("The number: %d\n", num)
	atomic.AddUint32(&num, ^uint32(-(-3) - 1))
	fmt.Printf("The number: %d\n", num)
}
