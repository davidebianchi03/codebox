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
