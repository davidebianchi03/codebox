import React from "react";
import { RunnerAdmin } from "../types/runner";
import { Badge } from "reactstrap";

export function AdminRunnerStatusBadge({ runner }: { runner: RunnerAdmin }) {
    if (new Date(runner?.last_contact || "") > new Date(Date.now() - 5 * 60 * 1000)) {
        return (
            <React.Fragment>
                <Badge
                    style={{ background: "rgba(var(--tblr-success-rgb), 0.1) !important" }}
                    className="border-success"
                >
                    <span className="text-success pe-1">●</span>
                    <span className="text-success">Online</span>
                </Badge>
            </React.Fragment>
        )
    } else {
        return (
            <React.Fragment>
                <Badge
                    style={{ background: "rgba(var(--tblr-danger-rgb), 0.1) !important" }}
                    className="border-danger"
                >
                    <span className="text-danger pe-1">●</span>
                    <span className="text-danger">Offline</span>
                </Badge>
            </React.Fragment>
        )
    }
}