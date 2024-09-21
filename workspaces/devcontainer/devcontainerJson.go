package devcontainer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"codebox.com/db"
	"codebox.com/utils"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/tailscale/hujson"
	"gopkg.in/yaml.v2"
)

type DevcontainerJson struct {
	workspace                *db.Workspace
	devcontainerJsonFilePath string
	devcontainerJson         map[string]interface{}
	dockerComposeExists      bool
	dockerComposeFilePath    string
	dockerComposeYaml        map[string]interface{}
}

var codeBoxDeniedDevcontainerJsonKeys = [...]string{
	"workspaceMount",
}

func standardizeJSON(b []byte) ([]byte, error) {
	ast, err := hujson.Parse(b)
	if err != nil {
		return b, err
	}
	ast.Standardize()
	return ast.Pack(), nil
}

func InitDevcontainerJson(workspace *db.Workspace, devcontainerJsonFilePath string) *DevcontainerJson {
	var obj DevcontainerJson
	obj.devcontainerJsonFilePath = devcontainerJsonFilePath
	obj.workspace = workspace
	return &obj
}

func (dj *DevcontainerJson) LoadConfigFromFiles() error {
	info, err := os.Stat(dj.devcontainerJsonFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", dj.devcontainerJsonFilePath)
		} else {
			return fmt.Errorf("unknown error: %s", err.Error())
		}
	}

	if info.IsDir() {
		return fmt.Errorf("%s is a directory", dj.devcontainerJsonFilePath)
	}

	data, err := os.ReadFile(dj.devcontainerJsonFilePath)
	if err != nil {
		return fmt.Errorf("cannot read devcontainer.json file: %s", err)
	}

	data, err = standardizeJSON(data)
	if err != nil {
		return fmt.Errorf("cannot parse devcontainer.json file: %s", err)
	}

	err = json.Unmarshal(data, &dj.devcontainerJson)

	if err != nil {
		return fmt.Errorf("cannot parse devcontainer.json file: %s", err)
	}

	for key := range dj.devcontainerJson {
		for _, deniedKey := range codeBoxDeniedDevcontainerJsonKeys {
			if key == deniedKey {
				return fmt.Errorf("%s: this devcontainer.json key is not allowed in codebox", key)
			}
		}
	}

	composeFilePathInterface, found := dj.devcontainerJson["dockerComposeFile"]

	if found {
		composeFilePath, ok := composeFilePathInterface.(string)
		if !ok {
			return fmt.Errorf("'dockerComposeFile' key in devcontainer.json is not a string")
		}
		composeAbsPath := path.Join(filepath.Dir(dj.devcontainerJsonFilePath), composeFilePath)
		_, err := os.Stat(composeAbsPath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("docker compose file not found")
			} else {
				return fmt.Errorf("unknown error: %s", err.Error())
			}
		}

		dj.dockerComposeFilePath = composeAbsPath

		data, err = os.ReadFile(dj.dockerComposeFilePath)
		if err != nil {
			return fmt.Errorf("cannot read docker-compose file: %s", err)
		}

		err = yaml.Unmarshal(data, &dj.dockerComposeYaml)
		if err != nil {
			return fmt.Errorf("cannot parse docker-compose yaml: %s", err)
		}

		dj.dockerComposeExists = true
	} else {
		dj.dockerComposeExists = false
	}

	return nil
}

func getWorkspaceVolumeId(workspace db.Workspace) string {
	return fmt.Sprintf("codebox-workspace-%s-%d-data", workspace.Name, workspace.ID)
}

const STACK_WORKSPACE_VOLUME_NAME = "codebox_workspace"

func getUsedPorts() ([]uint, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize docker client: %s", err)
	}
	defer dockerClient.Close()

	containers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list existing docker containers: %s", err)
	}

	var usedPorts []uint

	for _, container := range containers {
		for _, port := range container.Ports {
			usedPorts = append(usedPorts, uint(port.PublicPort))
		}
	}

	return usedPorts, nil
}

// https://containers.dev/implementors/json_schema/
func (js *DevcontainerJson) FixConfigFiles() error {
	// modifiche al file devcontainer.json
	_, found := js.devcontainerJson["workspaceFolder"].(string)
	if !found {
		js.devcontainerJson["workspaceFolder"] = "/workspace"
	}

	js.devcontainerJson["workspaceMount"] = fmt.Sprintf(
		"source=%s,target=%s,type=volume",
		getWorkspaceVolumeId(*js.workspace),
		js.devcontainerJson["workspaceFolder"],
	)

	usedPorts, err := getUsedPorts()
	if err != nil {
		return err
	}

	forwardedPortsMinInterval := 50000
	forwardedPortsMaxInterval := 60000

	if js.dockerComposeExists {
		// aggiustamenti specifici per i workspaces con lo stack
		composeServiceName, ok := js.devcontainerJson["service"].(string)
		if !ok {
			return fmt.Errorf("missing or invalid key 'service' in devcontainer.json")
		}

		dockerComposeServices, found := js.dockerComposeYaml["services"].(map[interface{}]interface{})
		if !found {
			return fmt.Errorf("missing 'services' tag in docker-compose.yml")
		}

		var dockerStackContainerNames []string
		for serviceName, serviceDefinition := range dockerComposeServices {
			dockerStackContainerNames = append(dockerStackContainerNames, serviceName.(string))

			service, err := utils.Interface2StringMap(serviceDefinition)
			if err != nil {
				return fmt.Errorf("invalid docker compose syntax for service '%s'", serviceName)
			}

			// gestione dei volumi
			volumes, found := service["volumes"]
			if found {
				volumesList, ok := volumes.([]interface{})
				if !ok {
					return fmt.Errorf("invalid docker compose syntax for service '%s', 'volumes' tag is not a list", serviceName)
				}

				// ottengo il nome del volume dove è verrà clonato il repository
				workspaceVolume := ""
				for _, volume := range volumesList {
					volumeStr, ok := volume.(string)
					if !ok {
						return fmt.Errorf("invalid docker compose syntax for service '%s', volume %s is not a string", serviceName, volume)
					}
					volumeParts := strings.Split(volumeStr, ":")
					if len(volumeParts) != 2 {
						return fmt.Errorf("invalid docker compose syntax for service '%s'", serviceName)
					}

					workspaceMountPoint := volumeParts[1]

					if workspaceMountPoint == js.devcontainerJson["workspaceFolder"] && composeServiceName == serviceName {
						// questo è il volume del workspace, sostituire il bind alla macchina locale con un volume
						workspaceVolume = volumeParts[0]
						break
					}
				}

				for _, volume := range volumesList {
					volumeStr, ok := volume.(string)
					if !ok {
						return fmt.Errorf("invalid docker compose syntax for service '%s', volume %s is not a string", serviceName, volume)
					}
					volumeParts := strings.Split(volumeStr, ":")
					if len(volumeParts) != 2 {
						return fmt.Errorf("invalid docker compose syntax for service '%s'", serviceName)
					}

					if volumeParts[0] == workspaceVolume {
						volume = STACK_WORKSPACE_VOLUME_NAME
					}
				}

				service["volumes"] = volumesList
			}

			// gestione delle porte
		}
	} else {
		agentExternalPort := -1
		// aggiustamenti specifici per i workspaces senza lo stack
		for i := forwardedPortsMinInterval; i < forwardedPortsMaxInterval; i++ {
			for _, port := range usedPorts {
				if i != int(port) {
					agentExternalPort = i
					break
				}
			}
		}

		if agentExternalPort < 0 {
			return fmt.Errorf("no free ports found")
		}

		js.devcontainerJson["runArgs"] = [...]string{
			"--publish",
			fmt.Sprintf("%d:55088", agentExternalPort),
			"--label",
			fmt.Sprintf("com.codebox.workspace_id=%d", js.workspace.ID),
		}
	}

	return nil
}
