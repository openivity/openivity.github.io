package result

import (
	"time"
)

type Encode struct {
	Err                  error
	DeserializeInputTook time.Duration
	EncodeTook           time.Duration
	SerializationTook    time.Duration
	TotalElapsed         time.Duration
	FileName             string
	FileType             string
	FilesBytes           [][]byte
}

func (e Encode) ToMap() map[string]any {
	if e.Err != nil {
		return map[string]any{"err": e.Err.Error()}
	}

	begin := time.Now()

	filesBytes := make([]any, len(e.FilesBytes))
	for i := range e.FilesBytes {
		fileBytes := make([]any, len(e.FilesBytes[i]))
		for j := range e.FilesBytes[i] {
			fileBytes[j] = e.FilesBytes[i][j]
		}
		filesBytes[i] = fileBytes
	}

	e.SerializationTook = time.Since(begin)
	e.TotalElapsed = e.DeserializeInputTook + e.EncodeTook + e.SerializationTook

	m := map[string]any{
		"err":                  nil,
		"fileName":             e.FileName,
		"fileType":             e.FileType,
		"filesBytes":           filesBytes,
		"deserializeInputTook": e.DeserializeInputTook.Milliseconds(),
		"encodeTook":           e.EncodeTook.Milliseconds(),
		"serializationTook":    e.SerializationTook.Milliseconds(),
		"totalElapsed":         e.TotalElapsed.Milliseconds(),
	}

	return m
}
