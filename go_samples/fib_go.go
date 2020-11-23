package main

import "fmt"

func main() {
	var f1 int = 0
	var f2 int = 1

	const n int = 50

	for i := 0; i < n-2; i++ {
		temp := f2
		f2 = f1 + f2
		f1 = temp
	}

	fmt.Println(f2)
}
