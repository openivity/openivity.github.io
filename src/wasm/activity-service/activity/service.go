package activity

import (
	"context"
	"io"
)

type Service interface {
	Decode(ctx context.Context, r io.Reader) ([]Activity, error)
}
