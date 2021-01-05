package main

import (
	"fmt"
)

func main() {
	var a string
	for i := 0; i < 64; i++ {
		a += "a"
	}
	fmt.Println(a)
	fmt.Printf("字符串长度: %d\n", len(a))
}
