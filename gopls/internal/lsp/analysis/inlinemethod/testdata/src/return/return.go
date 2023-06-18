package main

import "fmt"

func Hello(name string) string {
	return "Hello, " + name
}

func main() {
	fmt.Println(Hello("世界")) // want "inline method"
}
