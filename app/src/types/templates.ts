export interface WorkspaceTemplate {
    id: number
    name: string
    type: string
    description: string
    icon: string
}

export interface WorkspaceTemplateVersion {
    id: number
    template_id: number
    name: string
    config_file_relative_path: string
    published: boolean
    published_on: string
    edited_on: string
}

export interface WorkspaceTemplateVersionTreeItem {
    name: string
    full_path: string
    type: "file" | "dir"
    children: WorkspaceTemplateVersionTreeItem[]
}

export interface WorkspaceTemplateVersionEntry {
    name: string
    type: "file" | "dir"
    content: string // base64 encoded
}
