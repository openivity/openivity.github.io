package activity

import (
	"bytes"
	"sync"
)

// Safe for concurrency, used for performance-critical situation such as marshaling activity.
var bufPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}
