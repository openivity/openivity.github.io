package mem_test

import (
	"testing"

	"github.com/openivity/activity-service/mem"
)

func TestAllocBuffer(t *testing.T) {
	alloc := testing.AllocsPerRun(1000, func() {
		buf := mem.GetBuffer()
		mem.PutBuffer(buf)
	})
	if alloc != 0 {
		t.Fatalf("expected alloc: 0, got: %f", alloc)
	}
}
