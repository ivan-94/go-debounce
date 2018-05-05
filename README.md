# go-debounce 

[![GoDoc](https://godoc.org/github.com/carney520/go-debounce?status.svg)](https://godoc.org/github.com/carney520/go-debounce) [![Build Status](https://travis-ci.org/carney520/go-debounce.svg?branch=master)](https://travis-ci.org/carney520/go-debounce)

go-debounce implement Golang version debounce.
The passed function which will postpone its execution until after `time.Duration`
have elapsed since the last time it was invoked.

## Usage

* Go get

```shell
go get github.com/carney520/go-debounce
```

* Example

```go
package main

import (
  "github.com/carney520/go-debounce"
  "time"
)

func main() {
  db := debounce.New(500 * time.Millisecond, func() {
    rebuild()
  })

  onFileChange(db.Trigger)
  <-exit
  db.Stop()
}
```
