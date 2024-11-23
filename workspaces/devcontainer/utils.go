package devcontainer

import (
	"context"
	"fmt"

	"codebox.com/db"
	"codebox.com/env"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/tailscale/hujson"
)

// funzione che rimuove i commenti da una stringa json
func standardizeJSON(b []byte) ([]byte, error) {
	ast, err := hujson.Parse(b)
	if err != nil {
		return b, err
	}
	ast.Standardize()
	return ast.Pack(), nil
}

// funzione che restituisce la lista delle porte utilizzate
func getUsedPorts() ([]int, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize docker client: %s", err)
	}
	defer dockerClient.Close()

	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list existing docker containers: %s", err)
	}

	var usedPorts []int

	for _, container := range containers {
		for _, port := range container.Ports {
			usedPorts = append(usedPorts, int(port.PublicPort))
		}
	}

	return usedPorts, nil
}

func getWorkspaceVolumeId(workspace db.Workspace) string {
	return fmt.Sprintf("%s-workspace-%s-%d-data", env.CodeBoxEnv.WorkspaceObjectsPrefix, workspace.Name, workspace.ID)
}
