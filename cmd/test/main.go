package main

import (
	"errors"
	"fmt"
)

type T int

func main() {
	var errs []error
	errs = append(errs, errors.New("test"))
	errs = append(errs, errors.New("test2"))
	errs = append(errs, errors.New("test3"))

	fmt.Println(errors.Join(errs...))
}
