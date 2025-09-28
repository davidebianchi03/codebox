import { AdminUser } from "./user";

export interface ImpersonationLogs {
    id: number;
    impersonator: AdminUser;
    impersonator_ip_address:string;
    impersonation_started_at:string;
    impersonation_finished_at:string | null;
    session_expired:boolean;
}
