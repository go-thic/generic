package main

import (
	"fmt"
	"github.com/go-thic/generic/stream"
)

func main() {
	stream.New(stream.StartCountingFrom(int64(0))).
		Limit(stream.Count(30)).
		Do(stream.Map(fibo())).
		Finally(stream.Do(func(fibo int64) {
			fmt.Println(fibo)
		}))
}

func fibo() func(f int64) int64 {
	previous := int64(0)
	next := int64(1)
	return func(f int64) int64 {
		next, previous = previous+next, next
		return next
	}
}
