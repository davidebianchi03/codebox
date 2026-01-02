import axios, { isAxiosError } from "axios";
import { Workspace } from "../types/workspace";
import { Runner, RunnerWithToken } from "../types/runner";
import { AdminStats } from "../types/admin";
import { ImpersonationLogs } from "../types/impersonationLogs";

export async function AdminRetrieveStats(): Promise<AdminStats | undefined> {
    try {
        const r = await axios.get<AdminStats>(`/api/v1/admin/stats`);
        return r.data;
    } catch {
        return undefined;
    }
}


export async function AdminListImpersonationLogs(
    userEmail: string
): Promise<ImpersonationLogs[] | undefined> {
    try {
        const r = await axios.get<ImpersonationLogs[]>(
            `/api/v1/admin/users/${userEmail}/impersonation-logs`
        );
        return r.data;
    } catch {
        return undefined;
    }
}

export enum AdminSendTestEmailResponse {
    SUCCESS,
    SEND_ERROR,
    UNKNOWN_ERROR,
}

export async function APIAdminSendTestEmail(): Promise<{ response: AdminSendTestEmailResponse, description: string }> {
    try {
        const r = await axios.post<{
            success: boolean,
            description: string,
        }>(
            `/api/v1/admin/send-test-email`
        );

        if (r.data.success) {
            return {
                response: AdminSendTestEmailResponse.SUCCESS,
                description: `Test email has been sent`,
            }
        } else {
            return {
                response: AdminSendTestEmailResponse.SEND_ERROR,
                description: r.data.description,
            }
        }
    } catch (error) {
        if (isAxiosError(error)) {
            return {
                response: AdminSendTestEmailResponse.UNKNOWN_ERROR,
                description: `Unknown error, server responded with status code ${error.response?.status}`
            }
        } else {
            return {
                response: AdminSendTestEmailResponse.UNKNOWN_ERROR,
                description: `Unknown error`
            }
        }
    }
}