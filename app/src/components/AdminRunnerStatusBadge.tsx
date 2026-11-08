import { RunnerAdmin } from "../types/runner";

type RunnerStatus = "online" | "offline";

type AdminRunnerStatusBadgeProps =
    | { runner: Pick<RunnerAdmin, "online">; status?: never; label?: string }
    | { runner?: never; status: RunnerStatus; label?: string };


export function AdminRunnerStatusBadge(props: AdminRunnerStatusBadgeProps) {
    const status = props.runner
        ? (props.runner.online ? "online" : "offline")
        : props.status;
    const label = props.label ?? (status === "online" ? "Online" : "Offline");

    return (
        <span className={`runner-status-badge is-${status}`}>
            <span className="runner-status-dot" />
            {label}
        </span>
    );
}