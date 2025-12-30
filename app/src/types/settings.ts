export interface AuthenticationSettings {
    is_signup_open: boolean;
    is_signup_restricted: boolean;
    allowed_emails_regex: string;
    blocked_emails_regex: string;
    users_must_be_approved: boolean;
    approved_by_default_emails_regex: string;
}