export interface AuthenticationSettings {
    is_signup_open: boolean;
    is_signup_restricted: boolean;
    allowed_emails_regex: string;
    blocked_emails_regex: string;
}