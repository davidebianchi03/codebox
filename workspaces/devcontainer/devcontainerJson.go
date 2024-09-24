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
	"reflect"
	"strconv"
	"strings"

	"codebox.com/db"
	"codebox.com/utils"
	"codebox.com/utils/cast"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"gopkg.in/yaml.v2"
)

type DevcontainerConfigInfo struct {
	containerNames               []string
	developmentContainerName     string
	workspaceLocationInContainer string
	containerId                  string
	composeProjectName           string
	remoteUser                   string
	forwardedPorts               map[string][]uint
}

type DevcontainerJson struct {
	workspace                *db.Workspace
	devcontainerJsonFilePath string
	devcontainerJson         map[string]interface{}
	dockerComposeExists      bool
	dockerComposeFilePath    string
	dockerComposeYaml        map[string]interface{}
	workingDir               string
	devcontainersInfo        DevcontainerConfigInfo
}

var codeBoxDeniedDevcontainerJsonKeys = [...]string{
	"workspaceMount",
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

const STACK_WORKSPACE_VOLUME_NAME = "codebox_workspace"

// https://containers.dev/implementors/json_schema/
func (js *DevcontainerJson) FixConfigFiles() error {
	// modifiche al file devcontainer.json
	_, found := js.devcontainerJson["workspaceFolder"].(string)
	if !found {
		js.devcontainerJson["workspaceFolder"] = "/workspace"
	}
	js.devcontainersInfo.workspaceLocationInContainer = js.devcontainerJson["workspaceFolder"].(string)

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

		js.devcontainersInfo.developmentContainerName = composeServiceName
		js.devcontainersInfo.containerNames = services

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
			if !utils.ItemInIntegersList(usedPorts, i) {
				agentExternalPort = i
				break
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

		forwardedPortsStrArray, ok := js.devcontainerJson["forwardPorts"].([]string)
		if ok {
			var forwardedPortsUIntArray []uint
			ok = true
			for _, port := range forwardedPortsStrArray {
				intPort, err := strconv.Atoi(port)
				if err != nil || intPort < 0 {
					ok = false
					break
				}
				forwardedPortsUIntArray = append(forwardedPortsUIntArray, uint(intPort))
			}
			if ok {
				js.devcontainersInfo.forwardedPorts["development"] = forwardedPortsUIntArray
			}
		} else {
			forwardedPortsUIntArray, ok := js.devcontainerJson["forwardPorts"].([]uint)
			if ok {
				js.devcontainersInfo.forwardedPorts["development"] = forwardedPortsUIntArray
			}
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
	if err != nil {
		return err
	}

	var parsedStdOut map[string]interface{}
	err = json.Unmarshal(stdOutBuffer.Bytes(), &parsedStdOut)
	if err != nil {
		return fmt.Errorf("cannot parse stdout logs")
	}

	containerId, ok := parsedStdOut["containerId"].(string)
	if !ok {
		return fmt.Errorf("invalid stdout logs")
	}
	js.devcontainersInfo.containerId = containerId

	remoteUser, ok := parsedStdOut["remoteUser"].(string)
	if !ok {
		return fmt.Errorf("invalid stdout logs")
	}
	js.devcontainersInfo.remoteUser = remoteUser

	composeProjectNameInt, found := parsedStdOut["composeProjectName"]
	if found {
		composeProjectName := composeProjectNameInt.(string)
		if !ok {
			return fmt.Errorf("invalid stdout logs")
		}
		js.devcontainersInfo.composeProjectName = composeProjectName
	}

	return nil
}

func (js *DevcontainerJson) MapContainers() error {
	if js.devcontainersInfo.composeProjectName != "" {
		// list containers in stack
	} else {
		// mapping dei workspace con singolo container
		dockerClient, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return fmt.Errorf("cannot initialize docker client: %s", err)
		}
		defer dockerClient.Close()

		containerInfo, err := dockerClient.ContainerInspect(context.Background(), js.devcontainersInfo.containerId)
		if err != nil {
			return fmt.Errorf("cannot retrieve container details: %s", err)
		}

		if len(containerInfo.HostConfig.PortBindings) != 1 {
			return fmt.Errorf("there are more than one ports exposed for container %s", js.devcontainersInfo.containerId)
		}

		hostConfigKey, ok := reflect.ValueOf(containerInfo.HostConfig.PortBindings).MapKeys()[0].Interface().(nat.Port)
		if !ok {
			return fmt.Errorf("unknown error occured while trying to retrieve exposed ports for container %s", js.devcontainersInfo.containerId)
		}

		if hostConfigKey != "55088/tcp" {
			return fmt.Errorf("unknown error occured while trying to retrieve exposed ports for container %s", js.devcontainersInfo.containerId)
		}

		if len(containerInfo.HostConfig.PortBindings[hostConfigKey]) != 1 {
			return fmt.Errorf("agent port not exposed for container %s", js.devcontainersInfo.containerId)
		}

		agentPort := containerInfo.HostConfig.PortBindings[hostConfigKey][0].HostPort
		agentPortUInt, err := strconv.Atoi(agentPort)
		if err != nil {
			return fmt.Errorf("unknown error occured while trying to retrieve exposed ports for container %s", js.devcontainersInfo.containerId)
		}

		workspaceContainer := db.WorkspaceContainer{
			Type:                       "docker_container",
			Name:                       "development",
			ContainerUser:              js.devcontainersInfo.remoteUser,
			ContainerStatus:            containerInfo.State.Status,
			AgentStatus:                db.WorkspaceContainerAgentStatusStarting,
			AgentExternalPort:          uint(agentPortUInt),
			CanConnectRemoteDeveloping: true,
			WorkspacePathInContainer:   js.devcontainersInfo.workspaceLocationInContainer,
			// ExternalIPv4
			// ForwardedPorts
		}

		// containerForwardedPorts, ok := js.devcontainersInfo.forwardedPorts["development"]
		// if ok {
		// 	for _, port := range containerForwardedPorts {
		// 		portObj := db.ForwardedPort{
		// 			PortNumber: port,
		// 			ConnectionType: ConnectionTypeHttp
		// 		}

		// 		workspaceContainer.ForwardedPorts = append(workspaceContainer.ForwardedPorts, &portObj)
		// 	}
		// }
	}

	return nil
}
