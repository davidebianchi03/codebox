import axios, { isAxiosError } from "axios";
import { Workspace } from "../types/workspace";
import { Runner, RunnerWithToken } from "../types/runner";
import { AdminStats } from "../types/admin";
import { ImpersonationLogs } from "../types/impersonationLogs";

export async function AdminRetrieveStats() : Promise<AdminStats | undefined> {
    try {
        const r = await axios.get<AdminStats>(`/api/v1/admin/stats`);
        return r.data;
    } catch {
        return undefined;
    }
}

export async function AdminListWorkspaces(): Promise<Workspace[] | undefined> {
    try {
        const r = await axios.get<Workspace[]>(`/api/v1/admin/workspaces`);
        return r.data;
    } catch {
        return undefined;
    }
}

export async function AdminListRunners(limit: number = -1): Promise<Runner[] | undefined> {
    try {
        const r = await axios.get<Runner[]>(`/api/v1/admin/runners`, {
            params: {
                limit: limit
            }
        });
        return r.data;
    } catch {
        return undefined;
    }
}

export async function AdminRetrieveRunnerById(runnerId: number): Promise<Runner | undefined | null> {
    try {
        const r = await axios.get<Runner>(`/api/v1/admin/runners/${runnerId}`);
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

export async function AdminCreateRunner(
    name: string,
    type: string,
    use_public_url: boolean,
    public_url: string,
): Promise<RunnerWithToken | undefined> {
    try {
        const r = await axios.post<RunnerWithToken>(
            `/api/v1/admin/runners`,
            {
                name: name,
                type: type,
                use_public_url: use_public_url,
                public_url: public_url,
            }
        );
        return r.data;
    } catch {
        return undefined;
    }
}

export async function AdminUpdateRunner(
    runnerId: number,
    name: string,
    type: string,
    use_public_url: boolean,
    public_url: string,
): Promise<Runner | undefined> {
    try {
        const r = await axios.put<Runner>(
            `/api/v1/admin/runners/${runnerId}`,
            {
                name: name,
                type: type,
                use_public_url: use_public_url,
                public_url: public_url,
            }
        );
        return r.data;
    } catch {
        return undefined;
    }
}

export async function AdminListImpersonationLogs(
    userEmail:string
) : Promise<ImpersonationLogs[]|undefined> {
     try {
        const r = await axios.get<ImpersonationLogs[]>(
            `/api/v1/admin/users/${userEmail}/impersonation-logs`
        );
        return r.data;
    } catch {
        return undefined;
    }
}
