package main

import (
	"fmt"
	"strconv"

	"github.com/go-thic/generic/stream"
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
		Do(stream.Filter(func(s string) bool {
			if _, err := strconv.Atoi(s); err == nil {
				return true
			}
			return false
		})).
		Finally(stream.Do(func(s string) {
			fmt.Println(s)
		}))
}
