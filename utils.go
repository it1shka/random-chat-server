package main

import "time"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func chunksOf[T any](array []T, chunkSize int) [][]T {
	chunks := [][]T{}
	for start := 0; start < len(array); start += chunkSize {
		chunk := []T{}
		for i := start; i < min(start+chunkSize, len(array)); i++ {
			chunk = append(chunk, array[i])
		}
		chunks = append(chunks, chunk)
	}
	return chunks
}

func setInfiniteLoop(duration time.Duration, function func()) {
	ticker := time.NewTicker(duration)
	go func() {
		for range ticker.C {
			function()
		}
	}()
}

type Json = map[string]interface{}
