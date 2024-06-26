// Copyright (C) 2023 Openivity

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

package activity

import (
	"context"
	"errors"
	"io"
)

var ErrNoActivity = errors.New("no activity")

// Service is activity service
type Service interface {
	// Decode decodes the given r into activities and returns any encountered errors.
	Decode(ctx context.Context, r io.Reader) ([]Activity, error)
	// Encode encodes the given activities into a slice of bytes and returns any encountered errors.
	Encode(ctx context.Context, activities []Activity) ([][]byte, error)
}
