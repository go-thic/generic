package main

import (
	"fmt"
	"github.com/go-thic/generic/stream"
	"strconv"
)

func FizzBuzz(number int) string {
	theString := func(s string) string { return s }

	makeMapper := func(forNumber int, toString string) func(func(string) string) func(string) string {
		return func(fun func(string) string) func(string) string {
			if number%forNumber == 0 {
				return func(_ string) string {
					return toString + fun("")
				}
			}
			return fun
		}
	}

	fizz := makeMapper(3, "Fizz")
	buzz := makeMapper(5, "Buzz")

	return fizz(buzz(theString))(strconv.Itoa(number))
}

func main() {
	stream.NewProvider(stream.StartCountingFrom(0)).
		Limit(stream.Count(100)).
		Do(stream.Map(FizzBuzz)).
		Finally(stream.Do(func(s string) {
			fmt.Println(s)
		}))
}
