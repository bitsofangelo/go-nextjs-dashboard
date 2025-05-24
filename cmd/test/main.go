package main

import (
	"fmt"
	"time"
)

type T int

func main() {
	_, err := time.Parse(time.RFC3339, "2025-05-20T12:30:05+08:00")
	fmt.Println(err)

	var t *int
	var i any
	i = t
	fmt.Println(i == nil)
}
