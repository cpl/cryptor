package chunker

import "errors"

var (
	errorDataSize   = errors.New("Chunker data size too small")
	errorChunkCount = errors.New("Invalid chunk count")
)
