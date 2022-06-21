package main

import (
	"fmt"
	"github.com/go-thic/generic/stream"
)

func main() {
	stream.New(stream.StartCountingFrom(1.5)).
		Limit(stream.Count(100)).
		Finally(stream.Do(func(elem float64) {
			fmt.Println(elem)
		}))
}
