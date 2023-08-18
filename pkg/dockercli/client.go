package dockercli

import (
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
)

func NewDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation())
}

func GetRegistryAuthString(username, password, url string) (string, error) {
	authType := registry.AuthConfig{
		Username: username,
		Password: password,
		ServerAddress: url,
	}
	return registry.EncodeAuthConfig(authType)
}