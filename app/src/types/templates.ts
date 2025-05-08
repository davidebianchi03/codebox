export interface WorkspaceTemplate {
    id: number
    name: string
    type: string
    description: string
    icon: string
}

export interface WorkspaceTemplateVersion {
    id: number
    name: string
    config_file_relative_path: string
    published: boolean
    edited_on: string
}
