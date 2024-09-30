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
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"gopkg.in/yaml.v2"
)

const (
	AgentSSHServerPort = 2222
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
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("cannot initialize docker client: %s", err)
	}
	defer dockerClient.Close()

	if js.devcontainersInfo.composeProjectName != "" {
		// list containers in stack
		stackContainers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{
			Filters: filters.NewArgs(
				filters.KeyValuePair{
					Key:   "label",
					Value: fmt.Sprintf("com.docker.compose.project=%s", js.devcontainersInfo.composeProjectName),
				},
			),
		})
		if err != nil {
			return fmt.Errorf("cannot list stack containers %s", err)
		}

		// mapping dei containers
		for _, container := range stackContainers {
			// retrieve name of container inside the stack
			containerNameInStack := ""
			for labelKey, labelValue := range container.Labels {
				if labelKey == "com.docker.compose.service" {
					containerNameInStack = labelValue
				}
			}

			if containerNameInStack == "" {
				return fmt.Errorf("failed to list stack containers, unknown name for container %s", container.ID)
			}

			// retrieve more details about the container
			containerInfo, err := dockerClient.ContainerInspect(context.Background(), container.ID)

			if err != nil {
				return fmt.Errorf("failed retrieve more details about container %s", container.ID)
			}

			// retrieve external agent port
			if len(containerInfo.HostConfig.PortBindings) != 1 {
				return fmt.Errorf("there are more than one ports exposed for container %s", container.ID)
			}

			hostConfigKey, ok := reflect.ValueOf(containerInfo.HostConfig.PortBindings).MapKeys()[0].Interface().(nat.Port)
			if !ok {
				return fmt.Errorf("unknown error occured while trying to retrieve exposed ports for container %s", container.ID)
			}

			if hostConfigKey != "55088/tcp" {
				return fmt.Errorf("unknown error occured while trying to retrieve exposed ports for container %s", container.ID)
			}

			if len(containerInfo.HostConfig.PortBindings[hostConfigKey]) != 1 {
				return fmt.Errorf("agent port not exposed for container %s", container.ID)
			}

			agentPort := containerInfo.HostConfig.PortBindings[hostConfigKey][0].HostPort
			agentPortUInt, err := strconv.Atoi(agentPort)
			if err != nil {
				return fmt.Errorf("unknown error occured while trying to retrieve exposed ports for container %s", container.ID)
			}

			canConnectRemoteDeveloping := false
			workspacePathInContainer := ""
			remoteUser := "root"
			if js.devcontainersInfo.developmentContainerName == containerNameInStack {
				canConnectRemoteDeveloping = true
				remoteUser = js.devcontainersInfo.remoteUser
				workspacePathInContainer = js.devcontainersInfo.workspaceLocationInContainer
			}

			workspaceContainer := db.WorkspaceContainer{
				Type:                       "docker_container",
				Name:                       containerNameInStack,
				ContainerUser:              remoteUser,
				ContainerStatus:            container.State,
				AgentStatus:                db.WorkspaceContainerAgentStatusStarting,
				AgentExternalPort:          uint(agentPortUInt),
				CanConnectRemoteDeveloping: canConnectRemoteDeveloping,
				WorkspacePathInContainer:   workspacePathInContainer,
				ExternalIPv4:               "172.17.0.1",
			}

			// mapping of exposed ports
			developmentContainerSSHPortFound := false
			containerForwardedPorts, ok := js.devcontainersInfo.forwardedPorts[containerNameInStack]
			if ok {
				for _, port := range containerForwardedPorts {
					portConnectionType := db.ConnectionTypeHttp
					if port == 22 {
						developmentContainerSSHPortFound = true
						portConnectionType = db.ConnectionTypeWS
					}

					portObj := db.ForwardedPort{
						PortNumber:     port,
						ConnectionType: portConnectionType,
						Public:         true,
					}
					workspaceContainer.ForwardedPorts = append(workspaceContainer.ForwardedPorts, portObj)
				}
			}

			if !developmentContainerSSHPortFound {
				portObj := db.ForwardedPort{
					PortNumber:     22,
					ConnectionType: db.ConnectionTypeWS,
					Public:         true,
				}
				workspaceContainer.ForwardedPorts = append(workspaceContainer.ForwardedPorts, portObj)
			}

			result := db.DB.Create(&workspaceContainer)

			if result.Error != nil {
				return fmt.Errorf("failed to create workspace container %s in DB %s", containerInfo.ID, result.Error)
			}
		}

	} else {
		// mapping dei workspace con singolo container
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
			ExternalIPv4:               "172.17.0.1",
		}

		developmentContainerSSHPortFound := false
		containerForwardedPorts, ok := js.devcontainersInfo.forwardedPorts["development"]
		if ok {
			for _, port := range containerForwardedPorts {
				portConnectionType := db.ConnectionTypeHttp
				if port == AgentSSHServerPort {
					developmentContainerSSHPortFound = true
					portConnectionType = db.ConnectionTypeWS
				}

				portObj := db.ForwardedPort{
					PortNumber:     port,
					ConnectionType: portConnectionType,
					Public:         true,
				}
				workspaceContainer.ForwardedPorts = append(workspaceContainer.ForwardedPorts, portObj)
			}
		}

		if !developmentContainerSSHPortFound {
			portObj := db.ForwardedPort{
				PortNumber:     AgentSSHServerPort,
				ConnectionType: db.ConnectionTypeWS,
				Public:         true,
			}
			workspaceContainer.ForwardedPorts = append(workspaceContainer.ForwardedPorts, portObj)
		}

		result := db.DB.Create(&workspaceContainer)

		if result.Error != nil {
			return fmt.Errorf("failed to create workspace container in DB %s", result.Error)
		}
	}

	return nil
}

func (js *DevcontainerJson) StartAgents() error {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("cannot initialize docker client: %s", err)
	}
	defer dockerClient.Close()

	// list workspace containers
	workspaceContainers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{
		Filters: filters.NewArgs(
			filters.KeyValuePair{
				Key:   "label",
				Value: fmt.Sprintf("com.codebox.workspace_id=%d", js.workspace.ID),
			},
		),
	})
	if err != nil {
		return fmt.Errorf("cannot list workspace containers %s", err)
	}

	for _, container := range workspaceContainers {
		containerInfo, err := dockerClient.ContainerInspect(context.Background(), container.ID)

		if err != nil {
			return fmt.Errorf("failed retrieve more details about container %s", container.ID)
		}

		containerName := "development"
		if js.dockerComposeExists {
			containerName = ""
			for labelKey, labelValue := range containerInfo.Config.Labels {
				if labelKey == "com.docker.compose.service" {
					containerName = labelValue
					break
				}
			}

			if containerName == "" {
				return fmt.Errorf("failed to retrieve container name in stack for container %s", containerInfo.ID)
			}
		}

		// create agent folder
		logs, err := runCommandInContainer(dockerClient, container.ID, []string{"mkdir -p codebox"}, "/opt", "root", []string{}, true)
		if err != nil {
			return fmt.Errorf("failed create agent folder on container %s, %s", container.ID, err)
		}

		// install agent
		err = putFileInContainer(dockerClient, container.ID, "/opt/codebox", "./agent.bin")
		if err != nil {
			return fmt.Errorf("failed to add agent to container %s, %s", container.ID, err)
		}

		// add server private key
		err = putFileInContainer(dockerClient, container.ID, "/opt/codebox", "./id_rsa")
		if err != nil {
			return fmt.Errorf("failed to add server's private key to container %s, %s", container.ID, err)
		}

		// add server public key
		err = putFileInContainer(dockerClient, container.ID, "/opt/codebox", "./id_rsa.pub")
		if err != nil {
			return fmt.Errorf("failed to add server's public key to container %s, %s", container.ID, err)
		}

		// start agent
		logs, err = runCommandInContainer(dockerClient, container.ID, []string{"./agent.bin"}, "/opt", "root", []string{}, true)
		if err != nil {
			return fmt.Errorf("failed to start agent on container %s, %s", container.ID, err)
		}

		js.workspace.Logs += fmt.Sprintf("<Container: %s> %s\n", container.ID, logs)
		db.DB.Save(js.workspace)
	}

	return nil
}
