package main

import (
	"math/rand"
)

func AppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if len(i) <= 0 {
			return slice
		}
		if ele == i {
			return slice
		}

	}
	return append(slice, i)
}

func Find(slice []string, val string) (string, bool) {
	for _, item := range slice {
		if item == val {
			return item, true
		}
	}
	return "", false
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
