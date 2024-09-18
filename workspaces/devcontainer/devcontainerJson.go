package devcontainer

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/tailscale/hujson"
	"gopkg.in/yaml.v2"
)

// https://containers.dev/implementors/json_schema/
type DevcontainerJson struct {
	devcontainerJsonFilePath string
	jsonData                 map[string]interface{}
	dockerComposeExists      bool
	dockerComposeFilePath    string
	dockerComposeYaml        map[string]interface{}
}

func standardizeJSON(b []byte) ([]byte, error) {
	ast, err := hujson.Parse(b)
	if err != nil {
		return b, err
	}
	ast.Standardize()
	return ast.Pack(), nil
}

func InitDevcontainerJson(devcontainerJsonFilePath string) *DevcontainerJson {
	var obj DevcontainerJson
	obj.devcontainerJsonFilePath = devcontainerJsonFilePath
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

	err = json.Unmarshal(data, &dj.jsonData)

	if err != nil {
		return fmt.Errorf("cannot parse devcontainer.json file: %s", err)
	}

	composeFilePathInterface, found := dj.jsonData["dockerComposeFile"]

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
	} else {
		dj.dockerComposeExists = false
	}

	return nil
}

func (js *DevcontainerJson) FixConfigFiles() error {
	return nil
}
