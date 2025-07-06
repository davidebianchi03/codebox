import axios, { isAxiosError } from "axios";
import { WorkspaceTemplate, WorkspaceTemplateVersion, WorkspaceTemplateVersionEntry, WorkspaceTemplateVersionTreeItem } from "../types/templates";
import { Workspace } from "../types/workspace";

export async function APIListTemplates(): Promise<WorkspaceTemplate[] | undefined> {
    try {
        const r = await axios.get<WorkspaceTemplate[]>(`/api/v1/templates`);
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APIRetrieveTemplateById(id: number): Promise<WorkspaceTemplate | undefined | null> {
    try {
        const r = await axios.get<WorkspaceTemplate>(`/api/v1/templates/${id}`);
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

export async function APIRetrieveTemplateByName(name: string): Promise<WorkspaceTemplate | undefined | null> {
    try {
        const r = await axios.get<WorkspaceTemplate>(`/api/v1/templates/${encodeURIComponent(name)}`);
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

export async function APICreateTemplate(
    name: string,
    type: string,
    description: string,
    icon: string,
): Promise<WorkspaceTemplate | undefined> {
    try {
        const r = await axios.post<WorkspaceTemplate>(
            `/api/v1/templates/`,
            {
                name: name,
                type: type,
                description: description,
                icon: icon
            }
        );
        return r.data;
    } catch (error) {
        return undefined;
    }
}

export async function APIUpdateTemplate(templateId: number, name: string): Promise<WorkspaceTemplate | undefined> {
    try {
        const r = await axios.put<WorkspaceTemplate>(
            `/api/v1/templates/${templateId}`,
            { name: name }
        );
        return r.data;
    } catch (error) {
        return undefined;
    }
}

export async function APIDeleteTemplate(templateId: number): Promise<boolean> {
    try {
        await axios.delete(`/api/v1/templates/${templateId}`);
        return true;
    } catch {
        return false;
    }
}

export async function APIRetrieveTemplateLatestVersion(templateId: number): Promise<WorkspaceTemplateVersion | undefined | null> {
    try {
        const r = await axios.get<WorkspaceTemplateVersion>(`/api/v1/templates/${templateId}/latest-version`);
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

export async function APIListTemplateVersionsByTemplate(templateId: number): Promise<WorkspaceTemplateVersion[] | undefined> {
    try {
        const r = await axios.get<WorkspaceTemplateVersion[]>(`/api/v1/templates/${templateId}/versions`);
        return r.data;
    } catch {
        return undefined;
    }
}


export async function APIRetrieveTemplateVersion(templateId: number, versionId: number): Promise<WorkspaceTemplateVersion | undefined | null> {
    try {
        const r = await axios.get<WorkspaceTemplateVersion>(`/api/v1/templates/${templateId}/versions/${versionId}`);
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

export async function APIUpdateTemplateVersion(
    templateId: number,
    versionId: number,
    name: string,
    configPath: string,
    published: boolean,
): Promise<WorkspaceTemplateVersion | undefined> {
    try {
        const r = await axios.patch<WorkspaceTemplateVersion>(
            `/api/v1/templates/${templateId}/versions/${versionId}`,
            {
                name: name,
                config_file_path: configPath,
                published: published,
            }
        );
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APIListWorkspacesByTemplate(templateId: number): Promise<Workspace[] | undefined> {
    try {
        const r = await axios.get<Workspace[]>(`/api/v1/templates/${templateId}/workspaces`);
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APIListTemplateVersionEntry(
    templateId: number,
    versionId: number,
): Promise<WorkspaceTemplateVersionTreeItem[] | undefined> {
    try {
        const r = await axios.get<WorkspaceTemplateVersionTreeItem[]>(`/api/v1/templates/${templateId}/versions/${versionId}/entries`);
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APIRetrieveTemplateVersionEntry(
    templateId: number,
    versionId: number,
    path: string
): Promise<WorkspaceTemplateVersionEntry | undefined> {
    try {
        const r = await axios.get<WorkspaceTemplateVersionEntry>(`/api/v1/templates/${templateId}/versions/${versionId}/entries/${encodeURIComponent(path)}`);
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APICreateTemplateVersionEntry(
    templateId: number,
    versionId: number,
    path: string,
    type: "file" | "dir",
    content: string
): Promise<WorkspaceTemplateVersionEntry | undefined> {
    try {
        const r = await axios.post<WorkspaceTemplateVersionEntry>(
            `/api/v1/templates/${templateId}/versions/${versionId}/entries`,
            {
                path: path,
                type: type,
                content: content,
            }
        );
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APIUpdateTemplateVersionEntry(
    templateId: number,
    versionId: number,
    path: string,
    newPath: string,
    content: string
): Promise<WorkspaceTemplateVersionEntry | undefined> {
    try {
        const r = await axios.put<WorkspaceTemplateVersionEntry>(
            `/api/v1/templates/${templateId}/versions/${versionId}/entries/${encodeURIComponent(path)}`,
            {
                path: newPath,
                content: content,
            }
        );
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APIDeleteTemplateVersionEntry(
    templateId: number,
    versionId: number,
    path: string
): Promise<boolean> {
    try {
        await axios.delete(`/api/v1/templates/${templateId}/versions/${versionId}/entries/${encodeURIComponent(path)}`);
        return true;
    } catch {
        return false;
    }
}
