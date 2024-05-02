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

import (
	"github.com/muktihari/fit/profile/typedef"
	"github.com/openivity/activity-service/activity"
)

type Encode struct {
	ToolMode       EncodeToolMode       `json:"toolMode"`       // Selected Encode Mode
	TargetFileType FileType             `json:"targetFileType"` // Either fit, gpx, or tcx
	ManufacturerID typedef.Manufacturer `json:"manufacturerId"` // Only for FIT FileType
	ProductID      uint16               `json:"productId"`      // Only for FIT FileType
	DeviceName     string               `json:"deviceName"`     // Only for non-FIT FileType
	Sports         []string             `json:"sports"`         // Change sports
	TrimMarkers    []EncodeMarker       `json:"trimMarkers"`    // Trim markers; If specified, len should match len(sessions).
	ConcealMarkers []EncodeMarker       `json:"concealMarkers"` // Conceal markers; If specified, len should match len(sessions).
	RemoveFields   []string             `json:"removeFields"`   // Remove spefified fields from all records.
	Activities     []activity.Activity  `json:"-"`
}

type EncodeToolMode byte

const (
	ToolModeUnknown EncodeToolMode = iota
	ToolModeEdit
	ToolModeCombine
	ToolModeSplitPerSession
)

func (e EncodeToolMode) String() string {
	switch e {
	case ToolModeEdit:
		return "edit"
	case ToolModeCombine:
		return "combine"
	case ToolModeSplitPerSession:
		return "split"
	default:
		return "unknown"
	}
}

type FileType byte

const (
	FileTypeUnsupported FileType = iota
	FileTypeFIT
	FileTypeGPX
	FileTypeTCX
)

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

type EncodeMarker struct {
	StartN int `json:"startN"`
	EndN   int `json:"endN"`
}
