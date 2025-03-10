import { WorkspaceType } from "./workspace"

export interface Runner {
    id: number;
    name: string;
    type: string;
    use_public_url: boolean;
    public_url: string;
    last_contact: string;
}

export interface RunnerType {
    id: string
    name: string
    description: string
    supported_types: WorkspaceType[]
}