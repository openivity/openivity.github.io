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
