package spec

import (
	"github.com/muktihari/openactivity-fit/activity"
)

type Encode struct {
	ToolMode       EncodeToolMode      `json:"toolMode"`       // Selected Encode Mode
	TargetFileType FileType            `json:"targetFileType"` // Either fit, gpx, or tcx
	ManufacturerID uint16              `json:"manufacturerId"` // Only for FIT FileType
	ProductID      uint16              `json:"productId"`      // Only for FIT FileType
	DeviceName     string              `json:"deviceName"`     // Only for non-FIT FileType
	Sports         []string            `json:"sports"`         // Change sports
	TrimMarkers    []EncodeMarker      `json:"trimMarkers"`    // Trim markers; If specified, len should match len(sessions).
	ConcealMarkers []EncodeMarker      `json:"concealMarkers"` // Conceal markers; If specified, len should match len(sessions).
	RemoveFields   []string            `json:"removeFields"`   // Remove spefified fields from all records.
	Activities     []activity.Activity `json:"-"`
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

type EncodeMarker struct {
	StartN int `json:"startN"`
	EndN   int `json:"endN"`
}
