package helper

import (
	"quizzly/pkg/structs/collections/slices"
	"strings"
	"time"
)

const (
	defaultWordsPerSecond int = 3
)

func ReadEstimation(text ...string) time.Duration {
	words := make([]string, 0, 2)
	for _, t := range text {
		words = append(words, splitText(t)...)
	}

	seconds := float64(len(words)) / float64(defaultWordsPerSecond)
	duration := time.Duration(seconds) * time.Second

	return duration
}

func splitText(text string) []string {
	return slices.Filter(strings.Split(text, " "), func(s string) bool {
		return len(s) > 2
	})
}
