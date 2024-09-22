package devcontainer

import (
	"fmt"
	"strings"

	"codebox.com/utils/cast"
)

func dockerComposeRetrieveServicesMap(composeDefinition map[string]interface{}) (map[string]interface{}, error) {
	dockerComposeServicesInt, found := composeDefinition["services"]
	if !found {
		return nil, fmt.Errorf("missing 'services' tag in docker-compose.yml")
	}

	dockerComposeServices, err := cast.Interface2StringMap(dockerComposeServicesInt)
	if err != nil {
		return nil, fmt.Errorf("invalid docker-compose syntax")
	}

	return dockerComposeServices, nil
}

// Ottiene la lista dei nomi dei container di uno stack docker
func dockerComposeRetrieveListOfServices(composeDefinition map[string]interface{}) ([]string, error) {
	dockerComposeServices, err := dockerComposeRetrieveServicesMap(composeDefinition)
	if err != nil {
		return nil, err
	}

	var services []string
	for serviceName := range dockerComposeServices {
		services = append(services, serviceName)
	}
	return services, nil
}

// Ottiene il nome del volume con i file del workspace nello stack,
// se non trova il volume ritorna una stringa vuota
func dockerComposeGetWorkspaceVolumeName(composeDefinition map[string]interface{}, devcontainerServiceName string, workspaceLocation string) (string, error) {
	dockerComposeServices, err := dockerComposeRetrieveServicesMap(composeDefinition)
	if err != nil {
		return "", err
	}

	for serviceName, serviceDefinition := range dockerComposeServices {
		if serviceName == devcontainerServiceName {
			serviceDefinitionMap, err := cast.Interface2StringMap(serviceDefinition)
			if err != nil {
				return "", fmt.Errorf("invalid docker-compose syntax")
			}

			volumesObj, found := serviceDefinitionMap["volumes"]
			if !found {
				return "", nil
			}

			volumesList, err := cast.Interface2StringArray(volumesObj)
			if err != nil {
				return "", fmt.Errorf("invalid docker-compose syntax")
			}

			for _, volume := range volumesList {
				volumeParts := strings.Split(volume, ":")
				if len(volumeParts) < 2 {
					return "", fmt.Errorf("invalid docker-compose syntax")
				}

				if volumeParts[1] == workspaceLocation {
					return volumeParts[0], nil
				}
			}

			return "", nil
		}
	}

	return "", fmt.Errorf("service with name %s doesn't exist in docker compose file", devcontainerServiceName)
}
