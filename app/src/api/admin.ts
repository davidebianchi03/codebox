import axios, { isAxiosError } from "axios";
import { AdminStats } from "../types/admin";
import { ImpersonationLogs } from "../types/impersonationLogs";
import { AnalyticsConfig } from "../types/analytics";

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

export async function APIAdminGetAnalyticsConfig(): Promise<AnalyticsConfig | undefined> {
    try {
        const r = await axios.get<AnalyticsConfig>(
            `/api/v1/admin/analytics-config`
        );
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APIAdminUpdateAnalyticsConfig(sendData: boolean): Promise<AnalyticsConfig | undefined> {
    try {
        const r = await axios.put<AnalyticsConfig>(
            `/api/v1/admin/analytics-config`,
            {
                send_analytics_data: sendData,
            }
        );
        return r.data;
    } catch {
        return undefined;
    }
}

export async function APIAdminGetAnalyticsPreviewContent(): Promise<string | undefined> {
    try {
        const r = await axios.get<object>(
            `/api/v1/admin/analytics-data-preview`
        );
        return JSON.stringify(r.data, null, "\t");
    } catch {
        return undefined;
    }
}

export async function APIAdminGetAnalyticsDataBannerSent(): Promise<boolean> {
    try {
        const r = await axios.get<{analytics_banner_sent: boolean}>(
            `/api/v1/admin/analytics-banner-sent`
        );
        return r.data.analytics_banner_sent;
    } catch {
        return false;
    }
}

export async function APIAdminSetAnalyticsDataBannerSent(): Promise<boolean> {
    try {
        await axios.post(`/api/v1/admin/analytics-banner-sent`);
        return true;
    } catch {
        return false;
    }
}
