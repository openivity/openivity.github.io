package activity

import (
	"fmt"
	"io"
	"time"
)

type Logging interface {
	Service
}

type logging struct {
	svc Service
}

func NewLogging(svc Service) Logging {
	return &logging{svc: svc}
}

func (l *logging) DecodeFit(rs []io.Reader) []Result {
	defer func(begin time.Time) {
		fmt.Printf("wasm: decode fit took: %s\n", time.Since(begin))
	}(time.Now())

	return l.svc.DecodeFit(rs)
}
