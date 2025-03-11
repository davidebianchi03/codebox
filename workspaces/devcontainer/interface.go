package devcontainer

// type DevcontainerWorkspace struct {
// 	common.CommonWorkspace
// }

// func (dw *DevcontainerWorkspace) StartWorkspace() {
// 	// ottengo la configurazione dal repo git se non esiste il file
// 	if dw.Workspace.WorkspaceConfigurationFiles == "" {
// 		err := common.RetrieveWorkspaceConfigFilesFromGitRepo(dw.Workspace)
// 		if err != nil {
// 			dw.Workspace.AppendLogs(err.Error() + "\n")
// 			dw.Workspace.Status = db.WorkspaceStatusError
// 			db.DB.Save(&dw.Workspace)
// 			return
// 		}
// 		dw.Workspace.AppendLogs("\n")
// 		db.DB.Save(&dw.Workspace)
// 	}

// 	if dw.Workspace.WorkspaceConfigurationFiles == "" {
// 		dw.Workspace.AppendLogs("missing workspace configuration files, if problem persists contact us\n")
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}

// 	// estraggo i file di configurazione in una cartella temporanea
// 	workingDir, err := os.MkdirTemp("", fmt.Sprintf("tmp_workspace_%d", dw.Workspace.ID))
// 	if err != nil {
// 		dw.Workspace.AppendLogs("cannot create temp working directory, " + err.Error() + "\n")
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}
// 	defer os.RemoveAll(workingDir)

// 	workingDir = path.Join(workingDir, fmt.Sprintf("%s_workspace_%d", env.CodeBoxEnv.WorkspaceObjectsPrefix, dw.Workspace.ID))
// 	err = os.MkdirAll(workingDir, 0777)
// 	if err != nil {
// 		dw.Workspace.AppendLogs("cannot create temp working directory, " + err.Error() + "\n")
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}

// 	configFilesLocation := path.Join(workingDir, dw.Workspace.GitRepoConfigurationFolder)
// 	err = utils.ExtractTarGz(dw.Workspace.WorkspaceConfigurationFiles, configFilesLocation)
// 	if err != nil {
// 		dw.Workspace.AppendLogs("cannot extract configuration files from targz archive, " + err.Error() + "\n")
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}

// 	// caricamento della configurazione dal file .devcontainer.json
// 	devcontainerConfig := InitDevcontainerJson(dw.Workspace, configFilesLocation)
// 	err = devcontainerConfig.LoadConfigFromFiles()
// 	if err != nil {
// 		dw.Workspace.AppendLogs(err.Error() + "\n")
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}

// 	// aggiustamento dei file di configurazione
// 	err = devcontainerConfig.FixConfigFiles()
// 	if err != nil {
// 		dw.Workspace.AppendLogs(err.Error() + "\n")
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}

// 	// creazione e avvio dei containers
// 	err = devcontainerConfig.GoUp()
// 	if err != nil {
// 		dw.Workspace.AppendLogs(err.Error() + "\n")
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}

// 	// mapping dei containers
// 	err = devcontainerConfig.MapContainers()
// 	if err != nil {
// 		dw.Workspace.AppendLogs(err.Error() + "\n")
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}

// 	// install and start agents
// 	dw.Workspace.AppendLogs("Starting agents...")
// 	err = devcontainerConfig.StartAgents()
// 	if err != nil {
// 		dw.Workspace.AppendLogs(err.Error() + "\n")
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}

// 	// ping agents
// 	dw.Workspace.AppendLogs("Checking status of agents...")
// 	err = devcontainerConfig.CheckAgents()
// 	if err != nil {
// 		dw.Workspace.AppendLogs(err.Error() + "\n")
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}

// 	// clone repo in development container
// 	err = devcontainerConfig.CloneRepoInWorkspace()
// 	if err != nil {
// 		dw.Workspace.AppendLogs(err.Error() + "\n")
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}

// 	dw.Workspace.AppendLogs("Workspace is now running...")
// 	dw.Workspace.Status = db.WorkspaceStatusRunning
// 	db.DB.Save(&dw.Workspace)
// }

// func (dw *DevcontainerWorkspace) StopWorkspace() {
// 	// retrieve workspace containers
// 	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
// 	if err != nil {
// 		dw.Workspace.AppendLogs(fmt.Sprintf("cannot initialize docker client: %s", err))
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}
// 	defer dockerClient.Close()

// 	// list workspace containers
// 	workspaceContainers, err := dockerClient.ContainerList(context.Background(), container.ListOptions{
// 		Filters: filters.NewArgs(
// 			filters.KeyValuePair{
// 				Key:   "label",
// 				Value: fmt.Sprintf("com.%s.workspace_id=%d", env.CodeBoxEnv.WorkspaceObjectsPrefix, dw.Workspace.ID),
// 			},
// 		),
// 	})
// 	if err != nil {
// 		dw.Workspace.AppendLogs(fmt.Sprintf("cannot list workspace containers %s", err))
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}

// 	// stop and remove containers
// 	allContainersStopped := true
// 	for _, workspaceContainer := range workspaceContainers {
// 		dw.Workspace.AppendLogs(fmt.Sprintf("stopping container %s...", workspaceContainer.ID))
// 		err = dockerClient.ContainerStop(context.Background(), workspaceContainer.ID, container.StopOptions{})
// 		if err != nil {
// 			dw.Workspace.AppendLogs(fmt.Sprintf("cannot stop container %s, %s", workspaceContainer.ID, err))
// 			allContainersStopped = false
// 			continue
// 		}

// 		dw.Workspace.AppendLogs(fmt.Sprintf("removing container %s...", workspaceContainer.ID))
// 		err = dockerClient.ContainerRemove(context.Background(), workspaceContainer.ID, container.RemoveOptions{})
// 		if err != nil {
// 			dw.Workspace.AppendLogs(fmt.Sprintf("cannot remove container %s, %s", workspaceContainer.ID, err))
// 			allContainersStopped = false
// 			continue
// 		}
// 	}

// 	if !allContainersStopped {
// 		dw.Workspace.AppendLogs("error, cannot stop some containers")
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}

// 	// remove unused networks
// 	pruneReport, err := dockerClient.NetworksPrune(context.Background(), filters.Args{})
// 	if err != nil {
// 		dw.Workspace.AppendLogs(fmt.Sprintf("failed to remove unused networks %s", err))
// 	} else {
// 		dw.Workspace.AppendLogs(
// 			fmt.Sprintf("successfully removed %d networks: %s",
// 				len(pruneReport.NetworksDeleted),
// 				strings.Join(pruneReport.NetworksDeleted, ","),
// 			),
// 		)
// 	}

// 	// remove containers from DB
// 	dbContainers := []db.WorkspaceContainer{}
// 	result := db.DB.Where(map[string]interface{}{"workspace_id": dw.Workspace.ID}).Preload("ForwardedPorts").Find(&dbContainers)
// 	if result.Error != nil {
// 		dw.Workspace.AppendLogs(fmt.Sprintf("cannot retrieve workspace containers from db, %s", result.Error))
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}

// 	for _, container := range dbContainers {
// 		// remove exposed ports from db
// 		for _, port := range container.ForwardedPorts {
// 			db.DB.Delete(&port)
// 		}
// 		db.DB.Delete(&container)
// 	}

// 	dw.Workspace.AppendLogs("Workspace has been stopped...")
// 	dw.Workspace.Status = db.WorkspaceStatusStopped
// 	db.DB.Save(&dw.Workspace)
// }

// func (dw *DevcontainerWorkspace) DeleteWorkspace() {
// 	dw.Workspace.AppendLogs("Deleting workspace...")
// 	dw.Workspace.Status = db.WorkspaceStatusStopping
// 	db.DB.Save(&dw.Workspace)

// 	// retrieve workspace containers
// 	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
// 	if err != nil {
// 		dw.Workspace.AppendLogs(fmt.Sprintf("cannot initialize docker client: %s", err))
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}
// 	defer dockerClient.Close()

// 	// remove workspace volume
// 	var workspaceVolumesIds []string
// 	volumes, err := dockerClient.VolumeList(context.Background(), volume.ListOptions{})
// 	if err != nil {
// 		dw.Workspace.AppendLogs(fmt.Sprintf("failed to list workspace volumes, %s", err))
// 	}

// 	for _, volume := range volumes.Volumes {
// 		// check if volume is member of a docker compose stack
// 		composeProjectName, found := volume.Labels["com.docker.compose.project"]
// 		if found {
// 			if strings.HasPrefix(composeProjectName, fmt.Sprintf("%s_workspace_%d", env.CodeBoxEnv.WorkspaceObjectsPrefix, dw.Workspace.ID)) {
// 				workspaceVolumesIds = append(workspaceVolumesIds, volume.Name)
// 				continue
// 			}
// 		}

// 		// otherwise workspace is a single container workspace, so find volume and delete it
// 		if strings.HasPrefix(volume.Name, getWorkspaceVolumeId(*dw.Workspace)) {
// 			workspaceVolumesIds = append(workspaceVolumesIds, volume.Name)
// 			continue
// 		}
// 	}

// 	for _, volumeId := range workspaceVolumesIds {
// 		dw.Workspace.AppendLogs(fmt.Sprintf("removing volume %s", volumeId))
// 		err = dockerClient.VolumeRemove(context.Background(), volumeId, true)
// 		if err != nil {
// 			dw.Workspace.AppendLogs(fmt.Sprintf("failed to remove volume %s", volumeId))
// 		}
// 	}

// 	// remove all related items from db
// 	dbContainers := []db.WorkspaceContainer{}
// 	result := db.DB.Where(map[string]interface{}{"workspace_id": dw.Workspace.ID}).Preload("ForwardedPorts").Find(&dbContainers)
// 	if result.Error != nil {
// 		dw.Workspace.AppendLogs(fmt.Sprintf("cannot retrieve workspace containers from db, %s", result.Error))
// 		dw.Workspace.Status = db.WorkspaceStatusError
// 		db.DB.Save(&dw.Workspace)
// 		return
// 	}

// 	for _, container := range dbContainers {
// 		// remove exposed ports from db
// 		for _, port := range container.ForwardedPorts {
// 			db.DB.Delete(&port)
// 		}
// 		db.DB.Delete(&container)
// 	}

// 	dw.Workspace.ClearLogs()
// 	dw.Workspace.Delete()
// }
