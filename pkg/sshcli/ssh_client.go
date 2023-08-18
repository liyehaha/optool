package sshcli

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/liyehaha/optool/pkg/utils"
)

func NewServerInfo(host, port, username, password string) ServerInfo {
	return ServerInfo{
		Host: host,
		Port: port,
		User: username,
		Password: password,
	}
}

func NewSshClientWithServerInfo(serverInfo ServerInfo, timeout time.Duration) (*SshClient, error) {
	return NewSshClient(serverInfo.Host, serverInfo.Port, serverInfo.User, serverInfo.Password, timeout)
}

func NewSshClient(host, port, username, password string, timeout time.Duration) (*SshClient, error) {
	if timeout == 0 {
		timeout = DEFAULT_SSH_TIMEOUT
	}
	serverInfo := NewServerInfo(host, port, username, password)
	clientConfig, err := NewSSHClientConfigWithPassword(serverInfo.User, serverInfo.Password, timeout)
	if err != nil {
		err = utils.NewError(err, "failed to create ssh client config")
		return nil, err
	}
	sshClient, err := NewOriginalSSHClient(serverInfo.Host, serverInfo.Port, clientConfig)
	if err != nil {
		err = utils.NewError(err, fmt.Sprintf("failed to create ssh client for host: %s", serverInfo.Host))
		return nil, err
	}
	scpClient, err := scp.NewClientBySSHWithTimeout(sshClient, timeout)
	if err != nil {
		err = utils.NewError(err, fmt.Sprintf("failed to create scp client for host: %s", serverInfo.Host))
		return nil, err
	}
	return &SshClient{
		ServerInfo: serverInfo,
		Sshclient: sshClient,
		Scpclient: scpClient,
	}, nil
}

func (c *SshClient) Copy(ctx context.Context, src, dst string) error {
	file, err := os.Open(src)
	if err != nil {
		err = utils.NewError(err, fmt.Sprintf("failed to open file %s", src))
		return err
	}
	err = c.Scpclient.CopyFromFile(ctx, *file, dst, "0744")
	defer file.Close()
	return err
}

func (c *SshClient) Cmd(cmd string) *remoteScript {
	return &remoteScript{
		_type:  cmdLine,
		client: c.Sshclient,
		script: bytes.NewBufferString(cmd + "\n"),
	}
}

func (c *SshClient) Terminal(config *TerminalConfig) *remoteShell {
	return &remoteShell{
		client:         c.Sshclient,
		terminalConfig: config,
		requestPty:     true,
	}
}

func (c *SshClient) Shell() *remoteShell {
	return &remoteShell{
		client:     c.Sshclient,
		requestPty: false,
	}
}