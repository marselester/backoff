package backoff

import (
	"context"
	"fmt"
	"math/rand"
)

func ExampleRun() {
	work := func(attempt int) error {
		fmt.Printf("%d attempt\n", attempt)
		return fmt.Errorf("timeout")
	}
	// Initialize the source of uniformly-distributed pseudo-random numbers for reproducible results.
	random := rand.New(rand.NewSource(1))

	r := NewDecorrJitter(
		WithRand(random),
		WithMaxRetries(1),
	)
	err := Run(context.Background(), r, work)
	fmt.Println(err)
	// Output:
	// 1 attempt
	// 2 attempt
	// timeout
}
