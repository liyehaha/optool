package utils

import (
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