package main

import "fmt"

type person struct {
	name string
}

func main() {
	p1 := person{"artem"}
	p2 := person{"artem2"}
	fmt.Println(p1 == p2)
}
