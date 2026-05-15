package serializers

import "gitlab.com/codebox4073715/codebox/runnerinterface"

type ContainerFileInfoSerializer struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	IsDir    bool   `json:"is_dir"`
	Size     int64  `json:"size"`
	Mode     string `json:"mode"`
	ModTime  int64  `json:"mod_time"`
	Owner    string `json:"owner"`
	Group    string `json:"group"`
	MimeType string `json:"mime_type"`
}

func LoadContainerFileInfoSerializer(
	fileInfo runnerinterface.ContainerFileInfo,
) *ContainerFileInfoSerializer {
	return &ContainerFileInfoSerializer{
		Name:     fileInfo.Name,
		Path:     fileInfo.Path,
		IsDir:    fileInfo.IsDir,
		Size:     fileInfo.Size,
		Mode:     fileInfo.Mode,
		ModTime:  fileInfo.ModTime,
		Owner:    fileInfo.Owner,
		Group:    fileInfo.Group,
		MimeType: fileInfo.MimeType,
	}
}

func LoadMultipleContainerFileInfoSerializers(
	fileInfo []runnerinterface.ContainerFileInfo,
) []ContainerFileInfoSerializer {
	serializers := make([]ContainerFileInfoSerializer, len(fileInfo))
	for i, file := range fileInfo {
		serializers[i] = *LoadContainerFileInfoSerializer(file)
	}
	return serializers
}

type FileContentSerializer struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Size    int64  `json:"size"`
}

func LoadFileContentSerializer(
	file runnerinterface.ContainerReadFileResponse,
) *FileContentSerializer {
	return &FileContentSerializer{
		Path:    file.Path,
		Content: file.Content,
		Size:    file.Size,
	}
}

type CommandResultSerializer struct {
	Command    string `json:"command"`
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	ExitCode   int    `json:"exit_code"`
	WasSuccess bool   `json:"was_success"`
}

func LoadCommandResultSerializer(
	file runnerinterface.ExecuteCommandResponse,
) *CommandResultSerializer {
	return &CommandResultSerializer{
		Command:    file.Command,
		Stdout:     file.Stdout,
		Stderr:     file.Stderr,
		ExitCode:   file.ExitCode,
		WasSuccess: file.WasSuccess,
	}
}
