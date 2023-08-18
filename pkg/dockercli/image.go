package dockercli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
)

type DockerCli struct {
	Client *client.Client
}

func NewDockerCli() (*DockerCli, error) {
	client, err := NewDockerClient()
	if err != nil {
		return nil, err
	}
	return &DockerCli{
		Client: client,
	}, nil
}

func (d *DockerCli) GetOriginalClient() *client.Client {
	return d.Client
}

func (d *DockerCli) PullImage(ctx context.Context, repo, tag, registryAuth string) (string, error) {
	opts := types.ImagePullOptions{}
	if registryAuth != "" {
		opts = types.ImagePullOptions{
			RegistryAuth: registryAuth,
		}
	}
	imgName := fmt.Sprintf("%s:%s", repo, tag)
	out, err := d.Client.ImagePull(ctx, imgName, opts)
	if err != nil {
		return "", err
	}
	io.Copy(io.Discard, out)
	return imgName, nil
}

func (d *DockerCli) PushImage(ctx context.Context, repo, tag, registryAuth string) (string, error) {
	opts := types.ImagePushOptions{}
	if registryAuth != "" {
		opts = types.ImagePushOptions{
			RegistryAuth: registryAuth,
		}
	}
	imgName := fmt.Sprintf("%s:%s", repo, tag)
	out, err := d.Client.ImagePush(ctx, imgName, opts)
	io.Copy(io.Discard, out)
	if err != nil {
		return "", err
	}
	return imgName, err
}

func (d *DockerCli) Load(ctx context.Context, imageFile string) error {
	file, err := os.Open(imageFile)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = d.Client.ImageLoad(ctx, file, true)
	return err
}

func (d *DockerCli) Save(ctx context.Context, ids []string, path string) error {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	out, err := d.Client.ImageSave(ctx, ids)
	if err != nil {
		return err
	}
	if _, err = io.Copy(file, out); err != nil {
		return err
	}
	return nil
}

func (d *DockerCli) ListImage(ctx context.Context) ([]types.ImageSummary, error) {
	sum, err := d.Client.ImageList(ctx, types.ImageListOptions{
		All:     true,
		Filters: filters.Args{},
	})
	if err != nil {
		return nil, err
	}
	return sum, nil
}

func (d *DockerCli) RemoveImage(ctx context.Context, imgIds []string) error {
	for _, imgId := range imgIds {
		_, err := d.Client.ImageRemove(ctx, imgId, types.ImageRemoveOptions{
			Force:         false,
			PruneChildren: true,
		})
		if errdefs.IsConflict(err) {
			_, err1 := d.Client.ImageRemove(ctx, imgId, types.ImageRemoveOptions{
				Force:         true,
				PruneChildren: true,
			})
			if errdefs.IsUnavailable(err1) {
				continue
			}
		} else {
			return err
		}
	}
	return nil
}

func (d *DockerCli) CleanVolume(ctx context.Context) ([]string, uint64, error) {
	res, err := d.Client.VolumesPrune(ctx, filters.Args{})
	if err != nil {
		return []string{}, 0, err
	}
	return res.VolumesDeleted, res.SpaceReclaimed, nil
}

func (d *DockerCli) CleanImage(ctx context.Context) ([]string, uint64, error) {
	res, err := d.Client.ImagesPrune(ctx, filters.Args{})
	if err != nil {
		return []string{}, 0, err
	}
	imageDeleted := []string{}
	for _, v := range res.ImagesDeleted {
		imageDeleted = append(imageDeleted, v.Deleted)
	}
	return imageDeleted, res.SpaceReclaimed, nil
}

func (d *DockerCli) CleanContainer(ctx context.Context) ([]string, uint64, error) {
	res, err := d.Client.ContainersPrune(ctx, filters.Args{})
	if err != nil {
		return []string{}, 0, err
	}
	return res.ContainersDeleted, res.SpaceReclaimed, nil
}

func (d *DockerCli) ListAllContainerIds(ctx context.Context) ([]types.ContainerJSON, error) {
	result := []types.ContainerJSON{}
	res, err := d.Client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return result, err
	}
	for _, v := range res {
		r, err := d.Client.ContainerInspect(ctx, v.ID)
		if err != nil {
			return result, err
		}
		result = append(result, r)
	}
	return result, nil
}