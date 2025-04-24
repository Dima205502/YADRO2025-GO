package bucket

import (
	"context"
	"time"
)

type Bucket struct {
	tokens chan struct{}
}

func New(ctx context.Context, rps, burst int) *Bucket {
	if rps <= 0 || burst <= 0 {
		panic("both rps and burst must be positive")
	}

	b := Bucket{
		tokens: make(chan struct{}, burst),
	}

	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(rps))
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				select {
				case b.tokens <- struct{}{}:
				default:
				}
			}

		}
	}()

	return &b
}

func (b *Bucket) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-b.tokens:
	}

	return nil
}
