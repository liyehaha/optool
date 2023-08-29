package utils

import (
	"bytes"
	"net/http"
	"os"
	"os/exec"
)

func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

func RunCmdLocal(command string, args ...string) ([]byte, error) {
	return exec.Command(command, args...).Output()
}

func DoRequestPost(url string, contentType string, data []byte) (*http.Response, error){
	return http.Post(url, contentType, bytes.NewBuffer(data))
}

func DoRequestGet(url string) (*http.Response, error) {
	return http.Get(url)
}