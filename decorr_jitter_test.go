package backoff

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func ExampleDecorrJitter() {
	work := func(attempt int) error {
		fmt.Printf("%d attempt;", attempt)
		return fmt.Errorf("timeout")
	}
	// Initialize the source of uniformly-distributed pseudo-random numbers for reproducible results.
	random := rand.New(rand.NewSource(1))

	r := NewDecorrJitter(
		WithRand(random),
		WithMaxRetries(5),
		WithMultiplier(30*time.Second),
		WithMaxWait(300*time.Second),
	)
	// There will be 6 calls to work func: the first attempt and 5 retries.
	for i := 1; r.Next(); i++ {
		if err := work(i); err == nil {
			break
		}

		if d := r.Delay(); d > 0 {
			fmt.Printf(" retrying in %v\n", d)
		}
	}
	// Output:
	// 1 attempt; retrying in 29.803266325s
	// 2 attempt; retrying in 51.825231051s
	// 3 attempt; retrying in 53.839799981s
	// 4 attempt; retrying in 1m36.445202543s
	// 5 attempt; retrying in 3m48.077368218s
	// 6 attempt;
}

func TestDecorrJitter(t *testing.T) {
	tests := map[string]struct {
		maxRetries int
		multiplier time.Duration
		maxWait    time.Duration

		want string
	}{
		"retry=0 mult=0 wait=0": {
			want: "[0s]",
		},
		"retry=1 mult=0 wait=0": {
			maxRetries: 1,
			want:       "[24.836055ms 0s]",
		},
		"retry=neg mult=0 wait=0": {
			maxRetries: -1,
			want:       "[24.836055ms 43.187692ms 44.866499ms 80.371002ms 190.064473ms 536.390311ms 276.211616ms 1.351602683s 2.341613026s 6.374037452s 0s]",
		},
		"retry=1 mult=1s wait=0": {
			maxRetries: 1,
			multiplier: time.Second,
			want:       "[993.44221ms 0s]",
		},
		"retry=1 mult=neg wait=0": {
			maxRetries: 1,
			multiplier: -1,
			want:       "[24.836055ms 0s]",
		},
		"retry=neg mult=neg wait=0": {
			maxRetries: -1,
			multiplier: -1,
			want:       "[24.836055ms 43.187692ms 44.866499ms 80.371002ms 190.064473ms 536.390311ms 276.211616ms 1.351602683s 2.341613026s 6.374037452s 0s]",
		},
		"retry=1 mult=1s wait=1s": {
			maxRetries: 1,
			multiplier: time.Second,
			maxWait:    time.Second,
			want:       "[993.44221ms 0s]",
		},
		"retry=1 mult=1s wait=neg": {
			maxRetries: 1,
			multiplier: time.Second,
			maxWait:    -1,
			want:       "[993.44221ms 0s]",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewDecorrJitter(
				WithRand(rand.New(rand.NewSource(1))),
				WithMaxRetries(tc.maxRetries),
				WithMultiplier(tc.multiplier),
				WithMaxWait(tc.maxWait),
			)

			var got []time.Duration
			for r.Next() {
				got = append(got, r.Delay())
			}
			if tc.want != fmt.Sprint(got) {
				t.Errorf("expected %v got %v", tc.want, got)
			}
		})
	}
}
