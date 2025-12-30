import axios, { isAxiosError } from "axios";
import { AdminUser } from "../types/user";

export async function AdminListUsers(limit: number = -1): Promise<AdminUser[] | undefined> {
    try {
        const r = await axios.get<AdminUser[]>(`/api/v1/admin/users`, {
            params: {
                limit: limit
            }
        });
        return r.data;
    } catch {
        return undefined;
    }
}

export async function AdminRetrieveUserByEmail(email: string): Promise<AdminUser | undefined | null> {
    try {
        const r = await axios.get<AdminUser>(`/api/v1/admin/users/${email}`);
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
): Promise<AdminUser | undefined> {
    try {
        const r = await axios.post<AdminUser>(
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
    isTemplateManager: boolean,
    isEmailVerified: boolean,
    approved: boolean,
): Promise<AdminUser | undefined> {
    try {
        const r = await axios.put<AdminUser>(
            `/api/v1/admin/users/${userEmail}`,
            {
                first_name: firstName,
                last_name: lastName,
                is_superuser: isAdmin,
                is_template_manager: isTemplateManager,
                email_verified: isEmailVerified,
                approved: approved,
            }
        );
        return r.data;
    } catch {
        return undefined;
    }
}


export async function AdminDeleteUser(
    userEmail: string,
): Promise<boolean> {
    try {
        await axios.delete(
            `/api/v1/admin/users/${userEmail}`
        );
        return true;
    } catch {
        return false;
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

export async function AdminImpersonateUser(userEmail:string): Promise<boolean> {
    try {
        await axios.post(`/api/v1/admin/users/${userEmail}/impersonate`);
        return true;
    } catch {
        return false;
    }
}

export async function StopImpersonation(): Promise<boolean> {
    try {
        await axios.post(`/api/v1/stop-impersonation`);
        return true;
    } catch {
        return false;
    }
}
