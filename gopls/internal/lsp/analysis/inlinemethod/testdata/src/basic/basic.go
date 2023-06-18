package main

import "fmt"

func PrintHello() {
	fmt.Println("Hello, 世界")
}

func main() {
	PrintHello() // want "inline method"
}
