package sshcli

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

var DEFAULT_SSH_TIMEOUT = 20 * time.Second

func NewSSHClientConfigWithPassword(username, passwordk string, timeout time.Duration) (*ssh.ClientConfig, error) {
	if timeout == 0 {
		timeout = DEFAULT_SSH_TIMEOUT
	}
	ssh_config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwordk),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
		Timeout: timeout,
	}
	return ssh_config, nil
}

func NewSSHClientConfigWithKeyPath(username, keyFilePath string, timeout time.Duration) (*ssh.ClientConfig, error) {
	if timeout == 0 {
		timeout = DEFAULT_SSH_TIMEOUT
	}
	key, err := os.ReadFile(keyFilePath)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: "user",
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: timeout,
	}
	return config, nil
}

func NewOriginalSSHClient(host, port string, clientConfig *ssh.ClientConfig) (*ssh.Client, error) {
	addr := fmt.Sprintf("%s:%s", host, port)
	return ssh.Dial("tcp", addr, clientConfig)
}