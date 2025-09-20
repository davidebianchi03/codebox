import axios, { isAxiosError } from "axios";
import { User } from "../types/user";
import { Workspace } from "../types/workspace";
import { Runner, RunnerWithToken } from "../types/runner";
import { AdminStats } from "../types/admin";

export async function AdminRetrieveStats() : Promise<AdminStats | undefined> {
    try {
        const r = await axios.get<AdminStats>(`/api/v1/admin/stats`);
        return r.data;
    } catch {
        return undefined;
    }
}

export async function AdminListUsers(limit: number = -1): Promise<User[] | undefined> {
    try {
        const r = await axios.get<User[]>(`/api/v1/admin/users`, {
            params: {
                limit: limit
            }
        });
        return r.data;
    } catch {
        return undefined;
    }
}

export async function AdminRetrieveUserByEmail(email: string): Promise<User | undefined | null> {
    try {
        const r = await axios.get<User>(`/api/v1/admin/users/${email}`);
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

export async function AdminCreateUser(
    email: string,
    password: string,
    firstName: string,
    lastName: string,
    isAdmin: boolean,
    isTemplateManager: boolean
): Promise<User | undefined> {
    try {
        const r = await axios.post<User>(
            `/api/v1/admin/users`,
            {
                email: email,
                password: password,
                first_name: firstName,
                last_name: lastName,
                is_superuser: isAdmin,
                is_template_manager: isTemplateManager,
            }
        );
        return r.data;
    } catch {
        return undefined;
    }
}

export async function AdminUpdateUser(
    userEmail: string,
    firstName: string,
    lastName: string,
    isAdmin: boolean,
    isTemplateManager: boolean
): Promise<User | undefined> {
    try {
        const r = await axios.put<User>(
            `/api/v1/admin/users/${userEmail}`,
            {
                first_name: firstName,
                last_name: lastName,
                is_superuser: isAdmin,
                is_template_manager: isTemplateManager,
            }
        );
        return r.data;
    } catch {
        return undefined;
    }
}

export async function AdminSetUserPassword(userEmail: string, newPassword: string): Promise<boolean> {
    try {
        await axios.post(
            `/api/v1/admin/users/${userEmail}/set-password`, {
            password: newPassword,
        }
        );
        return true;
    } catch {
        return false;
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

export async function AdminListRunners(): Promise<Runner[] | undefined> {
    try {
        const r = await axios.get<Runner[]>(`/api/v1/admin/runners`);
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