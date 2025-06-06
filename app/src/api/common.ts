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

export async function APIUpdateCurrentUserDetails(firstName: string, lastName: string) {
    try {
        await axios.patch(
            `/api/v1/auth/user-details`,
            {
                first_name: firstName,
                last_name: lastName,
            }
        );
        return true;
    } catch {
        return false;
    }
}

export async function APIRetrieveSshPublicKey(): Promise<string | undefined> {
    try {
        const r = await axios.patch<{ public_key: string }>(`/api/v1/auth/user-ssh-public-key`);
        return r.data.public_key;
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

export async function APIInitialUserExists(): Promise<boolean> {
    try {
        const response = await axios.get<{ exists: boolean }>(`/api/v1/auth/initial-user-exists`);
        return response.data.exists;
    } catch {
        return false;
    }
}

export async function APIChangePassword(currentPassword: string, newPassword: string): Promise<boolean> {
    try {
        await axios.post<{ exists: boolean }>(
            `/api/v1/auth/change-password`,
            {
                current_password: currentPassword,
                new_password: newPassword,
            }
        );
        return true;
    } catch {
        return false;
    }
}
