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

package spec

import "fmt"

type FileType byte

const (
	FileTypeUnsupported FileType = iota
	FileTypeFIT
	FileTypeGPX
	FileTypeTCX
)

var _ fmt.Stringer = FileType(0)

func (f FileType) String() string {
	switch f {
	case FileTypeFIT:
		return "fit"
	case FileTypeGPX:
		return "gpx"
	case FileTypeTCX:
		return "tcx"
	}
	return "unsupported"
}

func FileTypeFromString(s string) FileType {
	switch s {
	case "fit":
		return FileTypeFIT
	case "gpx":
		return FileTypeGPX
	case "tcx":
		return FileTypeTCX
	default:
		return FileTypeUnsupported
	}
}
