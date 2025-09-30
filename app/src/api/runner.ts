import axios, { isAxiosError } from "axios";
import { Runner, RunnerAdmin, RunnerType, RunnerWithToken } from "../types/runner";

export async function ListRunners(): Promise<Runner[] | undefined> {
    try {
        const r = await axios.get<Runner[]>(`/api/v1/runners`);
        return r.data;
    } catch {
        return undefined;
    }
}

export async function ListRunnerTypes(): Promise<RunnerType[] | undefined> {
    try {
        const r = await axios.get<RunnerType[]>(`/api/v1/runner-types`);
        return r.data;
    } catch {
        return undefined;
    }
}

export async function AdminListRunners(limit: number = -1): Promise<RunnerAdmin[] | undefined> {
    try {
        const r = await axios.get<RunnerAdmin[]>(`/api/v1/admin/runners`, {
            params: {
                limit: limit
            }
        });
        return r.data;
    } catch {
        return undefined;
    }
}

export async function AdminRetrieveRunnerById(runnerId: number): Promise<RunnerAdmin | undefined | null> {
    try {
        const r = await axios.get<RunnerAdmin>(`/api/v1/admin/runners/${runnerId}`);
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
): Promise<RunnerAdmin | undefined> {
    try {
        const r = await axios.put<RunnerAdmin>(
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