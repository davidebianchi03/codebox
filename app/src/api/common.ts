import axios from "axios";
import { User } from "../types/user";
import { InstanceSettings } from "../types/settings";

export async function RetrieveCurrentUserDetails(): Promise<User | undefined> {
    try {
        const response = await axios.get<User>(`/api/v1/auth/user-details`);
        return response.data;
    } catch {
        return undefined;
    }
}

export async function Logout(): Promise<boolean> {
    try {
        await axios.post<User>(`/api/v1/auth/logout`);
        return true;
    } catch {
        return false;
    }
}

export async function RetrieveInstanceSettings(): Promise<InstanceSettings | undefined> {
    try {
        const response = await axios.get<InstanceSettings>(`/api/v1/instance-settings`);
        return response.data;
    } catch {
        return undefined;
    }
}

export async function RequestApiToken(): Promise<string | undefined> {
    try {
        const response = await axios.post<{ token: string }>(`/api/v1/auth/cli-login`);
        return response.data.token;
    } catch {
        return undefined;
    }
}
