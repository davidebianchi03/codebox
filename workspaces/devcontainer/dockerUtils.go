package devcontainer

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func putFileInContainer(dockerClient *client.Client, containerID string, destinationPath string, localPath string) error {
	// read file
	content, err := os.ReadFile(localPath)
	if err != nil {
		return fmt.Errorf("failed to read %s, %s", localPath, err)
	}

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err = tw.WriteHeader(&tar.Header{
		Name: filepath.Base(localPath), // filename
		Mode: 0777,                     // permissions
		Size: int64(len(content)),      // filesize
	})
	if err != nil {
		return fmt.Errorf("docker copy: %v", err)
	}
	tw.Write([]byte(content))
	tw.Close()

	err = dockerClient.CopyToContainer(context.Background(), containerID, destinationPath, &buf, types.CopyToContainerOptions{})
	return err
}

func runCommandInContainer(
	dockerClient *client.Client,
	containerID string,
	command []string,
	workingDir string,
	user string,
	env []string,
	detach bool,
) (logs string, err error) {
	ctx := context.Background()
	execConfig := container.ExecOptions{
		Cmd:          command, // The command to run
		AttachStdout: true,    // Attach stdout
		AttachStderr: true,    // Attach stderr
		Tty:          false,   // Disable tty
		User:         user,
		WorkingDir:   workingDir,
		Env:          env,
		Detach:       detach,
	}

	execIDResp, err := dockerClient.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec instance: %v", err)
	}

	// Start the exec instance
	if detach {
		execStartCheck := container.ExecStartOptions{
			Tty:    false,
			Detach: true,
		}
		err := dockerClient.ContainerExecStart(ctx, execIDResp.ID, execStartCheck)
		if err != nil {
			return "", fmt.Errorf("failed to start exec instance: %v", err)
		}
		return "", nil
	} else {
		execStartCheck := container.ExecAttachOptions{
			Tty: false,
		}
		resp, err := dockerClient.ContainerExecAttach(ctx, execIDResp.ID, execStartCheck)
		if err != nil {
			return "", fmt.Errorf("failed to start exec instance: %v", err)
		}
		defer resp.Close()

		// Read output from the command
		var logsBuf bytes.Buffer
		_, err = io.Copy(&logsBuf, resp.Reader)
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("error reading exec output: %v", err)
		}

		return logsBuf.String(), nil
	}
}
