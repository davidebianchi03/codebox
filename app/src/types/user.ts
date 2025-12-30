export interface CurrentUser {
    email: string;
    first_name: string;
    last_name: string;
    is_superuser: boolean;
    is_template_manager: boolean;
    last_login: string | null;
    created_at: string;
    impersonated: boolean;
}

export interface AdminUser {
    email: string;
    first_name: string;
    last_name: string;
    is_superuser: boolean;
    is_template_manager: boolean;
    last_login: string | null;
    created_at: string;
    deletion_in_progress: boolean;
    email_verified: boolean;
    approved: boolean;
}

export interface User {
    first_name: string;
    last_name: string;
    last_login: string | null;
}
