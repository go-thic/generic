package main

import (
	"fmt"
	"github.com/go-thic/generic/stream"
)

func main() {
	stream.New(stream.WithGenerator(fibo())).
		Limit(stream.Count(30)).
		Finally(stream.Do(func(fibo int64) {
			fmt.Println(fibo)
		}))
}

func fibo() func() int64 {
	previous := int64(0)
	next := int64(1)
	return func() int64 {
		next, previous = previous+next, next
		return next
	}
}
