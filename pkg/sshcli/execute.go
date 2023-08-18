package sshcli

import (
	"context"
	"time"
)

func ExecuteShell(command string, sshi ServerInfo, timeout time.Duration) (string, error) {
	client, err := NewSshClientWithServerInfo(sshi, timeout)
	if err != nil {
		return "", err
	}
	output, err := client.Cmd(command).Output()
	return string(output), err
}


func CopyFile(ctx context.Context, src, dst string, sshi ServerInfo, timeout time.Duration) error {
	client, err := NewSshClientWithServerInfo(sshi, timeout)
	if err != nil {
		return err
	}
	return client.Copy(ctx, src, dst)
}
