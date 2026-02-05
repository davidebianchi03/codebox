import axios from "axios";
import { CurrentUser, User } from "../types/user";
import { AuthenticationSettings } from "../types/settings";

export async function RetrieveCurrentUserDetails(): Promise<CurrentUser | undefined> {
    try {
        const response = await axios.get<CurrentUser>(`/api/v1/auth/user-details`);
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
        const r = await axios.get<{ public_key: string }>(`/api/v1/auth/user-ssh-public-key`);
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

export async function APIAdminRetrieveAuthenticationSettings(): Promise<AuthenticationSettings | undefined> {
    try {
        const response = await axios.get<AuthenticationSettings>(`/api/v1/admin/authentication-settings`);
        return response.data;
    } catch {
        return undefined;
    }
}

export async function APIAdminUpdateAuthenticationSettings(
    isSignupOpen: boolean,
    isSignupRestricted: boolean,
    allowedEmailsRegex: string,
    blockedEmailsRegex: string,
    usersMustBeApproved: boolean,
    approvedByDefaultEmailsRegex: string,
): Promise<AuthenticationSettings | undefined> {
    try {
        const response = await axios.put<AuthenticationSettings>(
            `/api/v1/admin/authentication-settings`,
            {
                is_signup_open: isSignupOpen,
                is_signup_restricted: isSignupRestricted,
                allowed_emails_regex: allowedEmailsRegex,
                blocked_emails_regex: blockedEmailsRegex,
                users_must_be_approved: usersMustBeApproved,
                approved_by_default_emails_regex: approvedByDefaultEmailsRegex,
            }
        );
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

export async function APISignUpOpen(): Promise<boolean> {
    try {
        const response = await axios.get<{ is_signup_open: boolean }>(`/api/v1/auth/is-signup-open`);
        return response.data.is_signup_open;
    } catch {
        return false;
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

export async function APIAdminEmailServiceConfigured(): Promise<boolean> {
    try {
        const response = await axios.get<{ is_configured: boolean }>(`/api/v1/admin/email-service-configured`);
        return response.data.is_configured;
    } catch (error) {
        return false;
    }
}

export async function APICanResetPassword(): Promise<boolean> {
    try {
        const response = await axios.get<{ can_reset_password: boolean }>(`/api/v1/auth/can-reset-password`);
        return response.data.can_reset_password;
    } catch {
        return false;
    }
}
