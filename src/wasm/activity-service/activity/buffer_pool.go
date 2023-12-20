package activity

import (
	"bytes"
	"fmt"
	"sync"
)

// Safe for concurrency, used for performance-critical situation such as marshaling activity.
var bufPool = sync.Pool{
	New: func() any {
		fmt.Println("called")
		return new(bytes.Buffer)
	},
}
