package main

import (
	"math/rand"
	"os"
	"strings"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomString(strlen int) string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	result := make([]byte, strlen)
	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}
	return string(result)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ProcessMessage(message string) (string, string) {
	data := ""

	if strings.Contains(message, "[img]") {
		data = strings.Split(message, "]")[1]
		return data, "IMAGE"
	}

	if strings.Contains(message, "[video]") {
		data = strings.Split(message, "]")[1]
		return data, "VIDEO"
	}

	if strings.Contains(message, "[audio]") {
		data = strings.Split(message, "]")[1]
		return data, "AUDIO"
	}

	if strings.Contains(message, "[file]") {
		data = strings.Split(message, "]")[1]
		return data, "FILE"
	}

	return message, "TEXT"
}
