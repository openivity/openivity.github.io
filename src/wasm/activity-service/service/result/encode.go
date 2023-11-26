package result

import "time"

type Encode struct {
	Err               error
	EncodeTook        time.Duration
	SerializationTook time.Duration
	TotalElapsed      time.Duration
	FileName          string
	FileType          string
	FileBytes         []byte
}

func (e Encode) ToMap() map[string]any {
	if e.Err != nil {
		return map[string]any{"err": e.Err.Error()}
	}

	begin := time.Now()

	fileBytes := make([]any, len(e.FileBytes))
	for i := range e.FileBytes {
		fileBytes[i] = e.FileBytes[i]
	}

	e.SerializationTook = time.Since(begin)
	e.TotalElapsed = e.EncodeTook + e.SerializationTook

	m := map[string]any{
		"err":               nil,
		"fileName":          e.FileName,
		"fileType":          e.FileType,
		"fileBytes":         fileBytes,
		"decodeTook":        e.EncodeTook.Milliseconds(),
		"serializationTook": e.SerializationTook.Milliseconds(),
		"totalElapsed":      e.TotalElapsed.Milliseconds(),
	}

	return m
}
