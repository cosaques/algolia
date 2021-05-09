package main

import "fmt"

func main() {
	m := map[string]int{"artem": 1}
	if m, exist := m["artem2"]; !exist {
		fmt.Println("not exists", m)
	} else {
		fmt.Println(m)
	}
}
