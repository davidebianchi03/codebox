import { Runner } from "./runner"
import { WorkspaceTemplateVersion } from "./templates"
import { User } from "./user"

export interface GitSource {
    id: number
    repository_url: string
    ref_name: string
    config_file_relative_path: string
}

export interface Workspace {
    id: number
    name: string
    user: User
    status: string
    type: string
    runner: Runner | null
    config_source: string
    template_version: WorkspaceTemplateVersion
    git_source: GitSource|null
    config_source_file_path: string
    environment_variables: string[]
    created_at: string
    updated_at: string
}

export interface WorkspaceType {
    id: string
    name: string
    supported_config_sources: string[],
    config_files_default_path: string,
}

export interface WorkspaceContainer {
    container_id: string
    container_name: string
    container_image: string
    container_user_id: number
    container_user_name: string
    agent_last_contact: string
    created_at: string
    updated_at: string
}

export interface ContainerPort {
    service_name: string
    port_number: number
    public: boolean
    port_url: string
    created_at: string
    updated_at: string
}
