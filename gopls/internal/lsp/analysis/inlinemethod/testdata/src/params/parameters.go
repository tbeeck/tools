package main

import "fmt"

func PrintHello(name string) {
	fmt.Println("Hello,", name)
}

func main() {
	PrintHello("世界") // want "inline method"
}
