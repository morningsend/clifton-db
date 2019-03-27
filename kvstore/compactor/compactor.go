package compactor

import "context"

type Compactor interface {
	Compact(ctx context.Context) <-chan struct{}
}
