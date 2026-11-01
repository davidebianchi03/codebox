import React from "react";
import { RunnerAdmin } from "../types/runner";
import { Badge } from "reactstrap";

export function AdminRunnerStatusBadge({ runner }: { runner: RunnerAdmin }) {
    if (new Date(runner?.last_contact || "") > new Date(Date.now() - 5 * 60 * 1000)) {
        return (
            <React.Fragment>
                <Badge
                    className="bg-success-lt"
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
                    className="bg-danger-lt"
                >
                    <span className="text-danger pe-1">●</span>
                    <span className="text-danger">Offline</span>
                </Badge>
            </React.Fragment>
        )
    }
}