package chunker

import "errors"

var (
	ErrorDataSize            = errors.New("Chunker data size too small")
	ErrorDataSizeCompressoin = errors.New("Chunker data size too small, after compression")
	ErrorChunkCount          = errors.New("Invalid chunk count")
)
