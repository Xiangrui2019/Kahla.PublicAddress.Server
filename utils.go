package main

import (
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"Kahla.PublicAddress.Server/kahla"
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

func ProcessMessage(message string, client *kahla.Client) (string, string) {
	data := ""

	if strings.Contains(message, "[img]") {
		data = strings.Split(message, "]")[1]
		filekey := strings.Split(data, "-")[0]
		downloadurl := "https://oss.aiursoft.com/download/fromkey/" + filekey
		return downloadurl + "-" + data, "IMAGE"
	}

	if strings.Contains(message, "[video]") {
		data = strings.Split(message, "]")[1]
		downloadurl := "https://oss.aiursoft.com/download/fromkey/" + data
		return downloadurl + "-" + data, "VIDEO"
	}

	if strings.Contains(message, "[audio]") {
		data = strings.Split(message, "]")[1]
		filekey, _ := strconv.Atoi(data)
		downloadurl, _ := client.Oss.FileDownloadAddress(filekey)
		downloadurl = strings.Replace(downloadurl, "audio", "audio.ogg", -1)
		return downloadurl + "-" + data, "AUDIO"
	}

	if strings.Contains(message, "[file]") {
		data = strings.Split(message, "]")[1]
		filekey, _ := strconv.Atoi(strings.Split(data, "-")[0])
		downloadurl, _ := client.Oss.FileDownloadAddress(filekey)
		return downloadurl + "-" + data, "FILE"
	}

	return message, "TEXT"
}
