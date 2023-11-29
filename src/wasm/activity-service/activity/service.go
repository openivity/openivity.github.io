package activity

import (
	"context"
	"errors"
	"io"
)

var (
	ErrNoActivity = errors.New("no activity")
)

type Service interface {
	Decode(ctx context.Context, r io.Reader) ([]Activity, error)
	Encode(ctx context.Context, activities []Activity) ([][]byte, error)
}
