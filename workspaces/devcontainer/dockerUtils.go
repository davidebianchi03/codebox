package devcontainer

// func putFileInContainer(dockerClient *client.Client, containerID string, destinationPath string, localPath string) error {
// 	// read file
// 	content, err := os.ReadFile(localPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to read %s, %s", localPath, err)
// 	}

// 	var buf bytes.Buffer
// 	tw := tar.NewWriter(&buf)
// 	err = tw.WriteHeader(&tar.Header{
// 		Name: filepath.Base(localPath), // filename
// 		Mode: 0777,                     // permissions
// 		Size: int64(len(content)),      // filesize
// 	})
// 	if err != nil {
// 		return fmt.Errorf("docker copy: %v", err)
// 	}
// 	tw.Write([]byte(content))
// 	tw.Close()

// 	err = dockerClient.CopyToContainer(context.Background(), containerID, destinationPath, &buf, types.CopyToContainerOptions{})
// 	return err
// }

// func parseDockerMultiplexedStream(streamReader io.Reader) (stdout string, stderr string, err error) {
// 	header := make([]byte, 8)
// 	for {
// 		// Read the 8-byte header
// 		_, err := io.ReadFull(streamReader, header)
// 		if err == io.EOF {
// 			break // End of stream
// 		}
// 		if err != nil {
// 			return "", "", fmt.Errorf("failed to read header: %v", err)
// 		}

// 		// Parse the header
// 		streamType := header[0] // 1 for stdout, 2 for stderr
// 		payloadSize := binary.BigEndian.Uint32(header[4:8])

// 		// Read the payload
// 		payload := make([]byte, payloadSize)
// 		_, err = io.ReadFull(streamReader, payload)
// 		if err != nil {
// 			return "", "", fmt.Errorf("failed to read payload: %v", err)
// 		}

// 		// Append to the appropriate output
// 		if streamType == 1 {
// 			stdout += string(payload)
// 		} else if streamType == 2 {
// 			stderr += string(payload)
// 		}
// 	}
// 	return stdout, stderr, nil
// }

// type DockerCommandOutput struct {
// 	detached bool
// 	exitCode int
// 	stdOut   string
// 	stdErr   string
// }

// func runCommandInContainer(
// 	dockerClient *client.Client,
// 	containerID string,
// 	command []string,
// 	workingDir string,
// 	user string,
// 	env []string,
// 	detach bool,
// ) (cmdOut DockerCommandOutput, err error) {
// 	ctx := context.Background()
// 	execConfig := container.ExecOptions{
// 		Cmd:          command, // The command to run
// 		AttachStdout: true,    // Attach stdout
// 		AttachStderr: true,    // Attach stderr
// 		Tty:          false,   // Disable tty
// 		User:         user,
// 		WorkingDir:   workingDir,
// 		Env:          env,
// 		Detach:       detach,
// 	}

// 	execIDResp, err := dockerClient.ContainerExecCreate(ctx, containerID, execConfig)
// 	if err != nil {
// 		return DockerCommandOutput{}, fmt.Errorf("failed to create exec instance: %v", err)
// 	}

// 	// Start the exec instance
// 	if detach {
// 		execStartCheck := container.ExecStartOptions{
// 			Tty:    false,
// 			Detach: true,
// 		}
// 		err := dockerClient.ContainerExecStart(ctx, execIDResp.ID, execStartCheck)
// 		if err != nil {
// 			return DockerCommandOutput{}, fmt.Errorf("failed to start exec instance: %v", err)
// 		}
// 		return DockerCommandOutput{detached: true}, nil
// 	} else {
// 		execStartCheck := container.ExecAttachOptions{
// 			Tty: false,
// 		}
// 		attachResp, err := dockerClient.ContainerExecAttach(ctx, execIDResp.ID, execStartCheck)
// 		if err != nil {
// 			return DockerCommandOutput{}, fmt.Errorf("failed to start exec instance: %v", err)
// 		}
// 		defer attachResp.Close()

// 		inspectResp, err := dockerClient.ContainerExecInspect(ctx, execIDResp.ID)
// 		if err != nil {
// 			return DockerCommandOutput{}, fmt.Errorf("failed to inspect exec instance: %v", err)
// 		}

// 		stdOut, stdErr, err := parseDockerMultiplexedStream(attachResp.Reader)
// 		if err != nil {
// 			return DockerCommandOutput{}, fmt.Errorf("failed to parse logs stream, %s", err)
// 		}

// 		return DockerCommandOutput{
// 			detached: false,
// 			exitCode: inspectResp.ExitCode,
// 			stdOut:   stdOut,
// 			stdErr:   stdErr,
// 		}, nil
// 	}
// }
