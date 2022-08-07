# backoff 🤚 [![Documentation](https://godoc.org/github.com/marselester/backoff?status.svg)](https://pkg.go.dev/github.com/marselester/backoff) [![Go Report Card](https://goreportcard.com/badge/github.com/marselester/backoff)](https://goreportcard.com/report/github.com/marselester/backoff)

This package generates sleep durations in an exponentially backing-off,
jittered manner, making sure to mitigate any correlations.
It is a Go port of the exponential backoff algorithm
[recommended in Polly project](https://github.com/Polly-Contrib/Polly.Contrib.WaitAndRetry#new-jitter-recommendation).

> a new jitter formula characterised by very smooth and even distribution of retry intervals,
> a well-controlled median initial retry delay, and broadly exponential backoff.

```go
package main

import (
	"context"
	"fmt"

	"github.com/marselester/backoff"
)

func main() {
	ctx := context.Background()
	r := backoff.NewDecorrJitter(
		backoff.WithMaxRetries(5),
	)

	err := backoff.Run(ctx, r, func(attempt int) error {
		fmt.Printf("%d attempt\n", attempt)
		return fmt.Errorf("timeout")
	})
	if err != nil {
		fmt.Printf("all the attempts failed: %v\n", err)
	}
}
```

There are 6 calls to a function: the first attempt and 5 retries.

```
1 attempt
2 attempt
3 attempt
4 attempt
5 attempt
6 attempt
all the attempts failed: timeout
```
