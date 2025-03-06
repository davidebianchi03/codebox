import { User } from "./user"

export interface GitSource {
    id: number
    repository_url: string
}

export interface Workspace {
    id: number
    name: string
    user: User
    status: string
    type: string
    runner: any
    config_source: string
    template_version: any
    git_source: GitSource
    config_source_file_path: string
    environment_variables: string[]
    created_at: string
    updated_at: string
}

export interface WorkspaceType {
    id: string
    name: string
    supported_config_sources: string[]
}
