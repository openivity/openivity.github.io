// Copyright (C) 2024 Openivity

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package mem

import (
	"bytes"
	"sync"
)

const (
	megabytes = 1000 << 10

	// bufferSize is preallocated buffer size to handle most operations.
	//
	// FIT files are mostly small-sized, for example 190km ride is usually below 1MB in a FIT format,
	// but if it's in a GPX or TCX format, it can jump to 10-20MB.
	//
	// Having slightly bigger buffer can reduce the need to re-allocate slice.
	bufferSize = 50 * megabytes
)

var pool = sync.Pool{New: func() any { return bytes.NewBuffer(make([]byte, 0, bufferSize)) }}

// GetBuffer gets ready to use buffer from the pool. Buffer has been reset before returned.
func GetBuffer() *bytes.Buffer {
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// PutBuffer puts buffer back to the pool.
func PutBuffer(buf *bytes.Buffer) {
	pool.Put(buf)
}
