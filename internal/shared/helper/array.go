package helper

func ChunkSlice[T any](slice []T, chunkSize int) [][]T {
	if chunkSize <= 0 {
		return [][]T{slice}
	}

	var chunks [][]T
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
