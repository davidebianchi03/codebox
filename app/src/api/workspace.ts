import axios from "axios";
import { Workspace } from "../types/workspace";

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
        const r = await axios.post<Workspace>(`/api/v1/workspaces`, {
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