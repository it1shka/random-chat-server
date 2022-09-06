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

func keysOf[K comparable, V any](dict map[K]V) []K {
	keys := make([]K, len(dict))
	index := 0
	for key := range dict {
		keys[index] = key
		index++
	}
	return keys
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
