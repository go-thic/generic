# generic
Some more or less useful stuff with go generics.

I initially created this repo only for playing around with and eventually understanding go generics. Later on i thought, it would be useful to have some packages for structuring.

## optional

The package contains a generic optional datatype. Sometimes it's perhaps useful to have something like this:
```go
  s := someFuncReturningOptionalString().Val(OrZero[string]())  
  s := someFuncReturningOptionalString().Val(OrElse("SomeDefault"))
  s := someFuncReturningOptionalNumber().Val(OrElse(0))
```

## stream

A generic implementation of some rudimentary stream api

## playground

the place to play around, copy, paste, try and error :-)
