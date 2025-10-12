import axios, { isAxiosError } from "axios";
import { ContainerPort, Workspace, WorkspaceContainer, WorkspaceType } from "../types/workspace";

export async function APIListWorkspaces(): Promise<Workspace[] | undefined> {
    try {
        const r = await axios.get<Workspace[]>(`/api/v1/workspace`);
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APIRetrieveWorkspaceById(id: number): Promise<Workspace | null | undefined> {
    try {
        const r = await axios.get<Workspace>(`/api/v1/workspace/${id}`);
        return r.data;
    } catch (error) {
        if (isAxiosError(error)) {
            if (error.response?.status === 404) {
                return null;
            }
        }
        return undefined;
    }
}

export async function APICreateWorkspace(
    name: string,
    type: string,
    runner_id: number,
    config_source: string,
    git_repo_url: string,
    git_ref_name: string,
    config_source_path: string,
    environment_variables: string[],
    template_version_id: number,
): Promise<Workspace | undefined> {
    try {
        const r = await axios.post<Workspace>(`/api/v1/workspace`, {
            name: name,
            type: type,
            runner_id: runner_id,
            config_source: config_source,
            git_repo_url: git_repo_url,
            git_ref_name: git_ref_name,
            config_source_path: config_source_path,
            environment_variables: environment_variables,
            template_version_id: template_version_id,
        });
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APIUpdateWorkspace(
    workspaceId: number,
    git_repo_url: string,
    git_ref_name: string,
    config_source_path: string,
    environment_variables: string[],
    runner_id: number | null,
): Promise<Workspace | undefined> {
    try {
        const r = await axios.put<Workspace>(`/api/v1/workspace/${workspaceId}`, {
            git_repo_url: git_repo_url,
            git_ref_name: git_ref_name,
            config_source_path: config_source_path,
            environment_variables: environment_variables,
            runner_id: runner_id,
        });
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APIDeleteWorkspace(workspaceId: number, skipErrors: boolean): Promise<boolean> {
    try {
        await axios.delete<Workspace>(
            `/api/v1/workspace/${workspaceId}?skip_errors=${skipErrors ? "true" : "false"}`
        );
        return true
    } catch {
        return false;
    }
}

export async function APIStartWorkspace(workspaceId: number): Promise<boolean> {
    try {
        await axios.post<Workspace>(`/api/v1/workspace/${workspaceId}/start`);
        return true
    } catch {
        return false;
    }
}

export async function APIStopWorkspace(workspaceId: number): Promise<boolean> {
    try {
        await axios.post<Workspace>(`/api/v1/workspace/${workspaceId}/stop`);
        return true
    } catch {
        return false;
    }
}

export async function APIUpdateWorkspaceConfig(workspaceId: number): Promise<boolean> {
    try {
        await axios.post<Workspace>(`/api/v1/workspace/${workspaceId}/update-config`);
        return true
    } catch {
        return false;
    }
}

export async function APIRetrieveWorkspaceLogs(workspaceId: number): Promise<string | undefined> {
    try {
        const r = await axios.get<{ logs: string }>(`/api/v1/workspace/${workspaceId}/logs`);
        return r.data.logs;
    } catch {
        return undefined;
    }

}

export async function APIListWorkspacesTypes(): Promise<WorkspaceType[] | undefined> {
    try {
        const r = await axios.get<WorkspaceType[]>(`/api/v1/workspace-types`);
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APIListWorkspaceContainers(
    workspaceId: number,
): Promise<WorkspaceContainer[] | undefined> {
    try {
        const r = await axios.get<WorkspaceContainer[]>(`/api/v1/workspace/${workspaceId}/container`);
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APIRetrieveWorkspaceContainer(
    workspaceId: number,
    containerName: string,
): Promise<WorkspaceContainer | null | undefined> {
    try {
        const r = await axios.get<WorkspaceContainer>(`/api/v1/workspace/${workspaceId}/container/${containerName}`);
        return r.data;
    } catch (error) {
        if (isAxiosError(error)) {
            if (error.response?.status === 404) {
                return null;
            }
        }
        return undefined;
    }
}

export async function APIListWorkspaceContainerPorts(
    workspaceId: number,
    containerName: string,
): Promise<ContainerPort[] | undefined> {
    try {
        const r = await axios.get<ContainerPort[]>(`/api/v1/workspace/${workspaceId}/container/${containerName}/port`);
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APICreateWorkspaceContainerPort(
    workspaceId: number,
    containerName: string,
    portNumber: number,
    serviceName: string,
    publicExposed: boolean
): Promise<ContainerPort | undefined> {
    try {
        const r = await axios.post<ContainerPort>(
            `/api/v1/workspace/${workspaceId}/container/${containerName}/port`,
            {
                port_number: portNumber,
                service_name: serviceName,
                public: publicExposed,
            }
        );
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APIDeleteWorkspaceContainerPort(
    workspaceId: number,
    containerName: string,
    portNumber: number,
): Promise<boolean> {
    try {
        await axios.delete<ContainerPort>(
            `/api/v1/workspace/${workspaceId}/container/${containerName}/port/${portNumber}`
        );
        return true;
    } catch {
        return false;
    }
}
