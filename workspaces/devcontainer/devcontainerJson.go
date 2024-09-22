package devcontainer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"codebox.com/db"
	"codebox.com/utils"
	"codebox.com/utils/cast"
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
	workingDir               string
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

func InitDevcontainerJson(workspace *db.Workspace, workingDir string) *DevcontainerJson {
	var obj DevcontainerJson
	obj.devcontainerJsonFilePath = path.Join(workingDir, "devcontainer.json")
	obj.workingDir = workingDir
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

		// get names of docker compose services
		services, err := dockerComposeRetrieveListOfServices(js.dockerComposeYaml)
		if err != nil {
			return err
		}

		// assign a port to each container
		usedPorts, err := getUsedPorts()
		if err != nil {
			return err
		}
		// var composeAssignedPorts map[string]int
		composeAssignedPorts := make(map[string]int)
		for _, service := range services {
			for i := forwardedPortsMinInterval; i < forwardedPortsMaxInterval; i++ {
				if !utils.ItemInIntegersList(usedPorts, i) {
					usedPorts = append(usedPorts, i)
					composeAssignedPorts[service] = i
					break
				}

				if i >= forwardedPortsMaxInterval-1 {
					return fmt.Errorf("no free ports available")
				}
			}
		}

		// get workspace volume name
		composeWorkspaceVolume, err := dockerComposeGetWorkspaceVolumeName(
			js.dockerComposeYaml,
			composeServiceName,
			js.devcontainerJson["workspaceFolder"].(string),
		)

		if err != nil {
			return err
		}

		servicesMap, err := dockerComposeRetrieveServicesMap(js.dockerComposeYaml)
		if err != nil {
			return err
		}

		codeBoxWorkspaceVolumeName := "codebox_workspace_volume"
		for serviceName, serviceDefinition := range servicesMap {
			// update volumes
			serviceDefinitionMap, err := cast.Interface2StringMap(serviceDefinition)
			if err != nil {
				return fmt.Errorf("Invalid docker compose syntax")
			}

			serviceVolumes, found := serviceDefinitionMap["volumes"]
			if found {
				serviceVolumesArray, err := cast.Interface2StringArray(serviceVolumes)
				if err != nil {
					return fmt.Errorf("Invalid docker compose syntax")
				}

				var newVolumesArray []string
				for _, volumeStr := range serviceVolumesArray {
					volumeParts := strings.Split(volumeStr, ":")
					if len(volumeParts) < 2 {
						return fmt.Errorf("Invalid docker compose syntax")
					}

					if volumeParts[0] == composeWorkspaceVolume {
						newVolumesArray = append(newVolumesArray, fmt.Sprintf("%s:%s", codeBoxWorkspaceVolumeName, strings.Join(volumeParts[1:], ":")))
					} else {
						newVolumesArray = append(newVolumesArray, volumeStr)
					}
				}

				serviceDefinitionMap["volumes"] = newVolumesArray
			}

			// update ports
			var exposedPorts = []string{fmt.Sprintf("%d:55088", composeAssignedPorts[serviceName])}
			serviceDefinitionMap["ports"] = exposedPorts

			// update labels
			labels, found := serviceDefinitionMap["labels"]
			labelsMap := make(map[string]interface{})
			if found {
				labelsMap, err = cast.Interface2StringMap(labels)
				if err != nil {
					return fmt.Errorf("Invalid docker compose syntax")
				}
			}
			labelsMap["com.codebox.workspace_id"] = js.workspace.ID
			serviceDefinitionMap["labels"] = labelsMap

			// TODO: override entrypoint

			servicesMap[serviceName] = serviceDefinitionMap
		}

		fixedYaml := make(map[string]interface{})
		for key, value := range js.dockerComposeYaml {
			fixedYaml[key] = value
		}
		fixedYaml["services"] = servicesMap
		stackVolumes, found := fixedYaml["volumes"]
		stackVolumesList := make(map[string]interface{})
		if found {
			stackVolumesList, err = cast.Interface2StringMap(stackVolumes)
			if err != nil {
				return fmt.Errorf("Invalid docker compose syntax")
			}
		}
		stackVolumesList[codeBoxWorkspaceVolumeName] = nil
		fixedYaml["volumes"] = stackVolumesList

		// update yaml file
		fixedComposeBytes, err := yaml.Marshal(&fixedYaml)
		if err != nil {
			return fmt.Errorf("cannot serialize docker compose %s", err)
		}

		err = os.WriteFile(js.dockerComposeFilePath, fixedComposeBytes, 0644)
		if err != nil {
			return fmt.Errorf("cannot write docker compose file %s", err)
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

	// update json config file
	configFileBytes, err := json.Marshal(&js.devcontainerJson)
	if err != nil {
		return fmt.Errorf("cannot serialize devcontainer.json %s", err)
	}
	err = os.WriteFile(js.devcontainerJsonFilePath, configFileBytes, 0644)
	if err != nil {
		return fmt.Errorf("cannot write devcontainer.json file %s", err)
	}

	return nil
}

func (js *DevcontainerJson) GoUp() error {
	var stdErrBuffer bytes.Buffer
	var stdOutBuffer bytes.Buffer

	cmd := exec.Command("devcontainer", "up", "--workspace-folder", filepath.Dir(js.workingDir))
	cmd.Stderr = &stdErrBuffer
	cmd.Stdout = &stdOutBuffer
	cmdRunning := true
	stdErrIndex := 0
	stdOutIndex := 0
	go func() {
		// redirect logs to db field
		newStdErrLogs := ""
		newStdOutLogs := ""
		for cmdRunning {
			newStdErrLogs = stdErrBuffer.String()[stdErrIndex:]
			newStdOutLogs = stdOutBuffer.String()[stdOutIndex:]

			stdErrIndex += len(newStdErrLogs)
			stdOutIndex += len(newStdOutLogs)

			js.workspace.Logs += newStdErrLogs
			db.DB.Save(js.workspace)
		}
	}()
	err := cmd.Run()

	_ = err
	return nil
}
