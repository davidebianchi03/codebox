import { RunnerAdmin } from "../types/runner";

type RunnerStatus = "online" | "offline";

type AdminRunnerStatusBadgeProps =
    | { runner: Pick<RunnerAdmin, "last_contact">; status?: never; label?: string }
    | { runner?: never; status: RunnerStatus; label?: string };

export function IsRunnerOnline(runner: Pick<RunnerAdmin, "last_contact">, now = Date.now()) {
    return now - new Date(runner.last_contact || "").getTime() < 5 * 60 * 1000;
}

export function AdminRunnerStatusBadge(props: AdminRunnerStatusBadgeProps) {
    const status = props.runner
        ? (IsRunnerOnline(props.runner) ? "online" : "offline")
        : props.status;
    const label = props.label ?? (status === "online" ? "Online" : "Offline");

    return (
        <span className={`runner-status-badge is-${status}`}>
            <span className="runner-status-dot" />
            {label}
        </span>
    );
}